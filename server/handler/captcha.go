package handler

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/ArtalkJS/Artalk/internal/captcha"
	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func Captcha(router fiber.Router) {
	ca := router.Group("/captcha/", func(c *fiber.Ctx) error {
		if !config.Instance.Captcha.Enabled {
			return common.RespError(c, "Captcha disabled")
		}
		return c.Next()
	})
	{
		ca.Post("/refresh", captchaGet)
		ca.Get("/get", captchaGet)
		ca.Post("/get", captchaGet)
		ca.Post("/check", captchaCheck)
		ca.Post("/status", captchaStatus)
	}
}

// @Summary      Captcha Status
// @Description  Get the status of the user's captcha verification
// @Tags         Captcha
// @Success      200  {object}  common.JSONResult{data=object{is_pass=bool}}
// @Router       /captcha/status  [post]
func captchaStatus(c *fiber.Ctx) error {
	if common.IsReqNeedCaptchaCheck(c) {
		return common.RespData(c, common.Map{"is_pass": false})
	} else {
		return common.RespData(c, common.Map{"is_pass": true})
	}
}

// @Summary      Captcha Get
// @Description  Get a base64 encoded captcha image or a HTML page to verify for user
// @Tags         Captcha
// @Success      200  {object}  common.JSONResult{data=object{img_data=string}}
// @Router       /captcha/refresh  [post]
// @Router       /captcha/get      [get]
// @Router       /captcha/get      [post]
func captchaGet(c *fiber.Ctx) error {
	ip := c.IP()

	// ==============
	//  图片验证码
	// ==============
	if config.Instance.Captcha.CaptchaType == config.TypeImage {
		return common.RespData(c, common.Map{
			"img_data": common.GetNewImageCaptchaBase64(ip),
		})
	}

	// ==============
	//  iframe 验证码
	// ==============
	captchaType := config.Instance.Captcha.CaptchaType
	ca := captcha.NewCaptcha(captchaType)

	iframeHTML, err := captcha.GetIFrameHTML(captchaType)
	if err != nil {
		return common.RespError(c, "iframe load err: "+err.Error())
	}

	var buf bytes.Buffer

	t := template.New("")
	t.Parse(string(iframeHTML))
	t.Execute(&buf, ca.PageParams())

	// response html body
	c.Set(fiber.HeaderContentType, "text/html")

	// disable cache
	c.Set(fiber.HeaderCacheControl, "no-cache, no-store, must-revalidate")
	c.Set(fiber.HeaderPragma, "no-cache")
	c.Set(fiber.HeaderExpires, "0")

	return c.SendString(buf.String())
}

type ParamsCaptchaCheck struct {
	Value string `form:"value" validate:"required"`
}

// @Summary      Captcha Check
// @Description  Verify user enters correct captcha code
// @Tags         Captcha
// @Param        value  formData  string  true  "the captcha value to check"
// @Success      200  {object}  common.JSONResult
// @Failure      400  {object}  common.JSONResult{data=object{img_data=string}}
// @Router       /captcha/check [post]
func captchaCheck(c *fiber.Ctx) error {
	ip := c.IP()

	var p ParamsCaptchaCheck
	if isOK, resp := common.ParamsDecode(c, &p); !isOK {
		return resp
	}
	inputVal := p.Value

	// ===========
	//  图片验证码
	// ===========
	if config.Instance.Captcha.CaptchaType == config.TypeImage {
		isPass := strings.ToLower(inputVal) == common.GetImageCaptchaRealCode(ip)
		if isPass {
			// 验证码正确
			common.DisposeImageCaptcha(ip) // 销毁图片验证码
			common.OnCaptchaPass(c)
			return common.RespSuccess(c)
		} else {
			// 验证码错误
			common.DisposeImageCaptcha(ip)
			common.OnCaptchaFail(c)
			return common.RespError(c, i18n.T("Wrong captcha"), common.Map{
				"img_data": common.GetNewImageCaptchaBase64(ip),
			})
		}
	}

	// ==============
	//  iframe 验证码
	// ==============
	captchaType := config.Instance.Captcha.CaptchaType
	ca := captcha.NewCaptcha(captchaType)
	isPass, err := ca.Check(captcha.CaptchaPayload{
		CheckValue: p.Value,
		UserIP:     ip,
	})

	if isPass {
		// 验证成功
		common.OnCaptchaPass(c)
		return common.RespSuccess(c)
	} else {
		// 验证失败
		common.OnCaptchaFail(c)
		logrus.Error("[Captcha] Failed to verify: ", err)
		return common.RespError(c, i18n.T("Verification failed"), common.Map{
			"detail": err.Error(),
		})
	}
}
