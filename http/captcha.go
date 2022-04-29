package http

import (
	"bytes"
	"encoding/base64"
	"image/color"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/lib/captcha"
	"github.com/ArtalkJS/ArtalkGo/pkged"
	"github.com/eko/gocache/v2/store"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	imgCaptcha "github.com/steambap/captcha"
)

var (
	CaptchaExpiration = 5 * time.Minute // 验证码 5 分钟内有效
)

// 获取当前状态，是否需要验证
func (a *action) CaptchaStatus(c echo.Context) error {
	if IsReqNeedCaptchaCheck(c) {
		return RespData(c, Map{"is_pass": false})
	} else {
		return RespData(c, Map{"is_pass": true})
	}
}

// 获取验证码
func (a *action) CaptchaGet(c echo.Context) error {
	ip := c.RealIP()

	// ===========
	//  Geetest
	// ===========
	if config.Instance.Captcha.Geetest.Enabled {
		pageFile, _ := pkged.Open("/lib/captcha/pages/geetest.html")
		buf, _ := ioutil.ReadAll(pageFile)

		var page bytes.Buffer

		t := template.New("")
		t.Parse(string(buf))
		t.Execute(&page, map[string]interface{}{"gt_id": config.Instance.Captcha.Geetest.CaptchaID})

		return c.HTML(http.StatusOK, page.String())
	}

	// ===========
	//  图片验证码
	// ===========
	return RespData(c, Map{
		"img_data": GetNewImageCaptchaBase64(ip),
	})
}

type ParamsCaptchaCheck struct {
	Value string `mapstructure:"value" param:"required"`
}

// 验证
func (a *action) CaptchaCheck(c echo.Context) error {
	ip := c.RealIP()

	var p ParamsCaptchaCheck
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}
	inputVal := p.Value

	// ===========
	//  Geetest
	// ===========
	if config.Instance.Captcha.Geetest.Enabled {
		isPass, reason, err := captcha.GeetestCheck(inputVal)
		if err != nil {
			logrus.Error("[Geetest] 验证发生错误 ", err)
			return RespError(c, "Geetest API 错误")
		}

		if isPass {
			// 验证成功
			onCaptchaPass(c)
			return RespSuccess(c)
		} else {
			// 验证失败
			onCaptchaFail(c)
			return RespError(c, "验证失败", Map{
				"reason": reason,
			})
		}
	}

	// ===========
	//  图片验证码
	// ===========
	isPass := strings.ToLower(inputVal) == GetImageCaptchaRealCode(ip)
	if isPass {
		// 验证码正确
		DisposeImageCaptcha(ip) // 销毁图片验证码
		onCaptchaPass(c)
		return RespSuccess(c)
	} else {
		// 验证码错误
		DisposeImageCaptcha(ip)
		onCaptchaFail(c)
		return RespError(c, "验证码错误", Map{
			"img_data": GetNewImageCaptchaBase64(ip),
		})
	}
}

// 验证成功操作
func onCaptchaPass(c echo.Context) {
	ip := c.RealIP()

	setActionCount(c, config.Instance.Captcha.ActionLimit-1)
	SetAlwaysCaptchaMode_Pass(ip, true) // 允许 always mode pass
}

// 验证失败操作
func onCaptchaFail(c echo.Context) {
	ip := c.RealIP()

	RecordAction(c)                      // 记录操作
	SetAlwaysCaptchaMode_Pass(ip, false) // 取消 always mode pass
}

//#region 图片验证码
// 获取对应 IP 图片验证码正确的值
func GetImageCaptchaRealCode(ip string) string {
	var realVal string
	lib.CACHE.Get(lib.Ctx, "captcha:"+ip, &realVal)
	return strings.ToLower(realVal)
}

// 获取新验证码 base64 格式图片
func GetNewImageCaptchaBase64(ip string) string {
	// generate a image
	pngBuffer := bytes.NewBuffer([]byte{})
	data, _ := imgCaptcha.New(160, 40, func(o *imgCaptcha.Options) {
		o.FontScale = 1
		o.CurveNumber = 2
		o.FontDPI = 85.0
		o.Noise = 0.7
		o.BackgroundColor = color.White
	})
	data.WriteImage(pngBuffer)
	base64 := "data:image/png;base64," + base64.StdEncoding.EncodeToString(pngBuffer.Bytes())

	// save real code
	lib.CACHE.Set(lib.Ctx, "captcha:"+ip, data.Text, &store.Options{Expiration: CaptchaExpiration})

	return base64
}

// 销毁图片验证码
func DisposeImageCaptcha(ip string) {
	lib.CACHE.Delete(lib.Ctx, "captcha:"+ip)
}

//#endregion

// AlwaysMode 是否能 Pass (for 总是需要验证码的选项)
func GetAlwaysCaptchaMode_Pass(ip string) bool {
	var val string
	_, err := lib.CACHE.Get(lib.Ctx, "captcha-am-pass:"+ip, &val)
	return err == nil && val == "1"
}

// 设置 AlwaysMode 允许 Pass (for 总是需要验证码的选项)
func SetAlwaysCaptchaMode_Pass(ip string, pass bool) {
	val := "0"
	if pass {
		val = "1"
	}

	lib.CACHE.Set(lib.Ctx, "captcha-am-pass:"+ip, val, &store.Options{Expiration: CaptchaExpiration})
}
