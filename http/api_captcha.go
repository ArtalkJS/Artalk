package http

import (
	"bytes"
	"encoding/base64"
	"time"

	"github.com/dchest/captcha"
	"github.com/labstack/echo/v4"
)

var caGlobalStore = captcha.NewMemoryStore(100, 10*time.Minute) // TODO: 用 redis 或其他kv数据库来存更多值

func ActionCaptchaRefresh(c echo.Context) error {
	ip := c.RealIP()

	return RespData(c, Map{
		"img_data": GetBase64Image(ip, 150, 40),
	})
}

func ActionCaptchaCheck(c echo.Context) error {
	ip := c.RealIP()
	value := c.QueryParam("value")

	captcha.SetCustomStore(caGlobalStore)
	if captcha.VerifyString(ip, value) {
		return RespSuccess(c)
	} else {
		return RespError(c, "验证码错误")
	}
}

func GetBase64Image(id string, width, height int) string {
	png := bytes.NewBuffer([]byte{})
	d := captcha.RandomDigits(6)
	captcha.SetCustomStore(caGlobalStore)
	caGlobalStore.Set(id, d)
	captcha.NewImage(id, d, width, height).WriteTo(png)
	base64Code := "data:image/png;base64," + base64.StdEncoding.EncodeToString(png.Bytes())
	return base64Code
}
