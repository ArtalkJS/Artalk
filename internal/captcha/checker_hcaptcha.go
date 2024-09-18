package captcha

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/tidwall/gjson"
)

const HCAPTCHA_API = "https://api.hcaptcha.com/siteverify"

var _ Checker = (*HCaptchaChecker)(nil)

type HCaptchaChecker struct {
	User       *User
	SiteKey    string
	SecreteKey string
}

func NewHCaptchaChecker(conf *config.HCaptchaConf, user *User) *HCaptchaChecker {
	return &HCaptchaChecker{
		User:       user,
		SiteKey:    conf.SiteKey,
		SecreteKey: conf.SecretKey,
	}
}

func (c *HCaptchaChecker) Check(value string) (bool, error) {
	// 构建 POST 请求的参数
	values := make(url.Values)
	values.Add("secret", c.SecreteKey)
	values.Add("response", value)
	if c.User.IP != "" {
		values.Add("remoteip", c.User.IP)
	}

	// 发送 POST 请求
	url := HCAPTCHA_API
	cli := http.Client{Timeout: time.Second * 10} // 10s 超时
	resp, err := cli.PostForm(url, values)
	if err != nil || resp.StatusCode != 200 {
		return false, err
	}
	defer resp.Body.Close()

	// 解析响应内容
	respBuf, _ := io.ReadAll(resp.Body)
	success := gjson.GetBytes(respBuf, "success")
	if success.Exists() && success.Bool() {
		// 验证成功
		return true, nil
	} else {
		// 验证失败
		return false, fmt.Errorf("err reason: %s", gjson.GetBytes(respBuf, "error-codes").String())
	}
}

func (c *HCaptchaChecker) Type() CaptchaType {
	return IFrame
}

func (c *HCaptchaChecker) Get() ([]byte, error) {
	return RenderIFrame("hcaptcha.html", Map{
		"site_key": c.SiteKey,
	})
}
