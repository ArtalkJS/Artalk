package captcha

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
)

type TurnstileParams struct {
	SecreteKey string
	UserToken  string
	UserIP     string
}

func TurnstileCheck(p TurnstileParams) (isPass bool, reason string, err error) {
	apiURL := "https://challenges.cloudflare.com/turnstile/v0/siteverify"

	// 构建 POST 请求的参数
	values := url.Values{}
	values.Add("secret", p.SecreteKey)
	values.Add("response", p.UserToken)
	if p.UserIP != "" {
		values.Add("remoteip", p.UserIP)
	}

	// 发送 POST 请求
	resp, err := http.PostForm(apiURL, values)
	if err != nil {
		return false, "", errors.New("service interface exception: " + err.Error())
	}
	defer resp.Body.Close()

	// 解析响应内容
	resJsonBuf, _ := io.ReadAll(resp.Body)
	resJson := string(resJsonBuf)

	gSuccess := gjson.Get(resJson, "success")
	if !gSuccess.Exists() {
		return false, "", errors.New("response results are not as expected: " + resJson)
	}
	success := gSuccess.Bool()
	if success {
		// 验证成功
		return true, "", nil
	} else {
		// 验证失败
		gErrs := gjson.Get(resJson, "error-codes")
		if gErrs.Exists() {
			reason = gErrs.Raw
		}
		return false, reason, nil
	}
}
