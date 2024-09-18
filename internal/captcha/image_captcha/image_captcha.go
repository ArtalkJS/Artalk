package image_captcha

import (
	"bytes"
	"encoding/base64"
	"image/color"
	"strings"
	"time"

	"github.com/artalkjs/artalk/v2/internal/cache/simple_cache"
	"github.com/steambap/captcha"
)

const (
	CaptchaExpiration  = 5 * time.Minute // 验证码 5 分钟内有效
	CaptchaCachePrefix = "atk_captcha:"
)

var captchaStore = simple_cache.New()

// 获取对应 IP 图片验证码正确的值
func CheckImageCaptchaCode(ip string, code string) bool {
	realVal, isFound := captchaStore.Get(CaptchaCachePrefix + ip)
	return isFound && strings.EqualFold(realVal.(string), code)
}

// 获取新验证码 base64 格式图片
// (调用该函数将销毁原有验证码)
func GetNewImageCaptchaBase64(ip string) ([]byte, error) {
	// generate a image
	pngBuffer := bytes.NewBuffer([]byte{})
	data, err := captcha.New(160, 40, func(o *captcha.Options) {
		o.FontScale = 1
		o.CurveNumber = 2
		o.FontDPI = 85.0
		o.Noise = 0.7
		o.BackgroundColor = color.White
	})
	if err != nil {
		return nil, err
	}

	if err := data.WriteImage(pngBuffer); err != nil {
		return nil, err
	}

	base64 := "data:image/png;base64," + base64.StdEncoding.EncodeToString(pngBuffer.Bytes())

	// save real code
	captchaStore.Set(CaptchaCachePrefix+ip, data.Text, CaptchaExpiration)

	return []byte(base64), nil
}

// 销毁图片验证码
func InvalidateImageCaptcha(ip string) {
	captchaStore.Delete(CaptchaCachePrefix + ip)
}
