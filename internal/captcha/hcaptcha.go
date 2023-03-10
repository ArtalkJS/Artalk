package captcha

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/tidwall/gjson"
)

const HCAPTCHA_API = "https://api.hcaptcha.com/siteverify"

type HCaptcha struct {
	SiteKey    string
	SecreteKey string
}

var _ Captcha = (*HCaptcha)(nil)

func NewHCaptcha(conf *config.HCaptchaConf) *HCaptcha {
	return &HCaptcha{
		SiteKey:    conf.SiteKey,
		SecreteKey: conf.SecretKey,
	}
}

func (c *HCaptcha) Check(p CaptchaPayload) (bool, error) {
	// 构建 POST 请求的参数
	values := make(url.Values)
	values.Add("secret", c.SecreteKey)
	values.Add("response", p.CheckValue)
	if p.UserIP != "" {
		values.Add("remoteip", p.UserIP)
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
		return false, errors.New("err reason: " + gjson.GetBytes(respBuf, "error-codes").String())
	}
}

func (c *HCaptcha) PageParams() Map {
	return Map{
		"site_key": c.SiteKey,
	}
}
