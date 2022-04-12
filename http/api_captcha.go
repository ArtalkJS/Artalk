package http

import (
	"bytes"
	"encoding/base64"
	"image/color"
	"strings"
	"time"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/eko/gocache/v2/store"
	"github.com/labstack/echo/v4"
	"github.com/steambap/captcha"
)

var (
	CaptchaExpiration = 5 * time.Minute // 验证码 5 分钟内有效
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
		return RespError(c, "param `value` is required")
	}

	if strings.ToLower(inputVal) == GetCaptchaRealCode(ip) {
		// 验证码正确
		ResetActionRecord(c)         // 重置操作记录
		DisposeCaptcha(ip)           // 销毁验证码
		SetCaptchaIsCheked(ip, true) // 记录该 IP 已经成功验证

		return RespSuccess(c)
	} else {
		// 验证码错误
		RecordAction(c)               // 记录操作
		DisposeCaptcha(ip)            // 销毁验证码
		SetCaptchaIsCheked(ip, false) // 记录该 IP 验证码验证失败

		return RespError(c, "验证码错误", Map{
			"img_data": GetNewCaptchaImageBase64(ip),
		})
	}
}

// 获取 IP 验证码正确的值
func GetCaptchaRealCode(ip string) string {
	realVal := ""
	if val, err := lib.CACHE.Get(lib.Ctx, "captcha:"+ip); err == nil {
		realVal = string(val.([]byte))
	}
	return strings.ToLower(realVal)
}

// 获取新验证码 base64 格式图片
func GetNewCaptchaImageBase64(ip string) string {
	// generate a image
	pngBuffer := bytes.NewBuffer([]byte{})
	data, _ := captcha.New(160, 40, func(o *captcha.Options) {
		o.FontScale = 1
		o.CurveNumber = 2
		o.FontDPI = 85.0
		o.Noise = 0.7
		o.BackgroundColor = color.White
	})
	data.WriteImage(pngBuffer)
	base64 := "data:image/png;base64," + base64.StdEncoding.EncodeToString(pngBuffer.Bytes())

	// save real code
	lib.CACHE.Set(lib.Ctx, "captcha:"+ip, []byte(data.Text), &store.Options{Expiration: CaptchaExpiration})

	return base64
}

// 销毁验证码
func DisposeCaptcha(ip string) {
	lib.CACHE.Delete(lib.Ctx, "captcha:"+ip)
}

// 获取 IP 是否验证码检测通过，已经有一次
func GetCaptchaIsCheked(ip string) bool {
	val, err := lib.CACHE.Get(lib.Ctx, "captcha-checked:"+ip)
	return err == nil && string(val.([]byte)) == "1"
}

// 设置该 IP 验证码成功验证，已经有一次
func SetCaptchaIsCheked(ip string, isCheked bool) {
	val := "0"
	if isCheked {
		val = "1"
	}

	lib.CACHE.Set(lib.Ctx, "captcha-checked:"+ip, []byte(val), &store.Options{Expiration: CaptchaExpiration})
}
