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

var _ Checker = (*CapJSCaptchaChecker)(nil)

type CapJSCaptchaChecker struct {
	User        *User
	KeyID       string
	SecretKey   string
	APIEndpoint string
}

func NewCapJSChecker(conf *config.CapJSConf, user *User) *CapJSCaptchaChecker {
	return &CapJSCaptchaChecker{
		User:        user,
		KeyID:       conf.KeyID,
		SecretKey:   conf.SecretKey,
		APIEndpoint: conf.APIEndpoint,
	}
}

func (c *CapJSCaptchaChecker) Check(value string) (bool, error) {
	// 构建 POST 请求的参数
	values := make(url.Values)
	values.Add("secret", c.SecretKey)
	values.Add("response", value)

	// 发送 POST 请求
	url := c.APIEndpoint + "/" + c.KeyID + "/siteverify"
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

func (c *CapJSCaptchaChecker) Type() CaptchaType {
	return IFrame
}

func (c *CapJSCaptchaChecker) Get() ([]byte, error) {
	return RenderIFrame("capjs.html", Map{
		"api_endpoint": c.APIEndpoint,
		"key_id":       c.KeyID,
		"secret_key":   c.SecretKey,
	})
}
