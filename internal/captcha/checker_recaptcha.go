package captcha

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/ArtalkJS/Artalk/v2/internal/config"
	"github.com/tidwall/gjson"
)

const RECAPTCHA_API = "https://www.google.com/recaptcha/api/siteverify"

var _ Checker = (*ReCaptchaChecker)(nil)

type ReCaptchaChecker struct {
	User       *User
	SiteKey    string
	SecreteKey string
}

func NewReCaptchaChecker(conf *config.ReCaptchaConf, user *User) *ReCaptchaChecker {
	return &ReCaptchaChecker{
		User:       user,
		SiteKey:    conf.SiteKey,
		SecreteKey: conf.SecretKey,
	}
}

func (c *ReCaptchaChecker) Check(value string) (bool, error) {
	// 构建 POST 请求的参数
	values := make(url.Values)
	values.Add("secret", c.SecreteKey)
	values.Add("response", value)
	if c.User.IP != "" {
		values.Add("remoteip", c.User.IP)
	}

	// 发送 POST 请求
	url := RECAPTCHA_API
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

func (c *ReCaptchaChecker) Type() CaptchaType {
	return IFrame
}

func (c *ReCaptchaChecker) Get() ([]byte, error) {
	return RenderIFrame("recaptcha.html", Map{
		"site_key": c.SiteKey,
	})
}
