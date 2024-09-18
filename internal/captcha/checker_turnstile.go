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

const TURNSTILE_API = "https://challenges.cloudflare.com/turnstile/v0/siteverify"

var _ Checker = (*TurnstileChecker)(nil)

type TurnstileChecker struct {
	User       *User
	SiteKey    string
	SecreteKey string
}

func NewTurnstileChecker(conf *config.TurnstileConf, user *User) *TurnstileChecker {
	return &TurnstileChecker{
		User:       user,
		SiteKey:    conf.SiteKey,
		SecreteKey: conf.SecretKey,
	}
}

func (c *TurnstileChecker) Check(value string) (bool, error) {
	// 构建 POST 请求的参数
	values := make(url.Values)
	values.Add("secret", c.SecreteKey)
	values.Add("response", value)
	if c.User.IP != "" {
		values.Add("remoteip", c.User.IP)
	}

	// 发送 POST 请求
	url := TURNSTILE_API
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

func (c *TurnstileChecker) Type() CaptchaType {
	return IFrame
}

func (c *TurnstileChecker) Get() ([]byte, error) {
	return RenderIFrame("turnstile.html", Map{
		"site_key": c.SiteKey,
	})
}
