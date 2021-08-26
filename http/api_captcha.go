package http

import (
	"bytes"
	"encoding/base64"
	"strings"
	"time"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/eko/gocache/v2/store"
	"github.com/labstack/echo/v4"
	"github.com/steambap/captcha"
)

func ActionCaptchaGet(c echo.Context) error {
	ip := c.RealIP()

	return RespData(c, Map{
		"img_data": GetNewCaptchaImageBase64(ip),
	})
}

func ActionCaptchaCheck(c echo.Context) error {
	ip := c.RealIP()

	inputVal := c.QueryParam("value")
	if inputVal == "" {
		return RespError(c, "param `value` is required.")
	}

	if strings.ToLower(inputVal) == GetCaptchaRealCode(ip) {
		// 验证码正确
		ResetActionRecord(c) // 重置操作记录
		DisposeCaptcha(ip)   // 销毁验证码

		return RespSuccess(c)
	} else {
		// 验证码错误
		RecordAction(c) // 记录操作

		return RespError(c, "验证码错误", Map{
			"img_data": GetNewCaptchaImageBase64(ip),
		})
	}
}

func GetCaptchaRealCode(ip string) string {
	realVal := ""
	if val, err := lib.CACHE.Get(permCtx, "captcha:"+ip); err == nil {
		realVal = string(val.([]byte))
	}
	return strings.ToLower(realVal)
}

func GetNewCaptchaImageBase64(ip string) string {
	// generate a image
	pngBuffer := bytes.NewBuffer([]byte{})
	data, _ := captcha.New(150, 40)
	data.WriteImage(pngBuffer)
	base64 := "data:image/png;base64," + base64.StdEncoding.EncodeToString(pngBuffer.Bytes())

	// save real code
	lib.CACHE.Set(permCtx, "captcha:"+ip, []byte(data.Text), &store.Options{Expiration: 5 * time.Minute}) // 5分钟失效

	return base64
}

func DisposeCaptcha(ip string) {
	lib.CACHE.Delete(permCtx, "captcha:"+ip)
}
