package handler

import (
	"bytes"
	"io"
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

// 获取当前状态，是否需要验证
func captchaStatus(c *fiber.Ctx) error {
	if common.IsReqNeedCaptchaCheck(c) {
		return common.RespData(c, common.Map{"is_pass": false})
	} else {
		return common.RespData(c, common.Map{"is_pass": true})
	}
}

// 获取验证码
func captchaGet(c *fiber.Ctx) error {
	ip := c.IP()

	// ===========
	//  Geetest
	// ===========
	if config.Instance.Captcha.Geetest.Enabled {
		pageFile, _ := captcha.GetPage("geetest.html")
		buf, _ := io.ReadAll(pageFile)

		var page bytes.Buffer

		t := template.New("")
		t.Parse(string(buf))
		t.Execute(&page, map[string]interface{}{"gt_id": config.Instance.Captcha.Geetest.CaptchaID})

		c.Set("Content-Type", "text/html")
		return c.SendString(page.String())
	}

	// ===========
	//  图片验证码
	// ===========
	return common.RespData(c, common.Map{
		"img_data": common.GetNewImageCaptchaBase64(ip),
	})
}

type ParamsCaptchaCheck struct {
	Value string `form:"value" validate:"required"`
}

// 验证
func captchaCheck(c *fiber.Ctx) error {
	ip := c.IP()

	var p ParamsCaptchaCheck
	if isOK, resp := common.ParamsDecode(c, &p); !isOK {
		return resp
	}
	inputVal := p.Value

	// ===========
	//  Geetest
	// ===========
	if config.Instance.Captcha.Geetest.Enabled {
		isPass, reason, err := captcha.GeetestCheck(inputVal)
		if err != nil {
			logrus.Error("[Geetest] Failed to verify: ", err)
			return common.RespError(c, "Geetest API error")
		}

		if isPass {
			// 验证成功
			common.OnCaptchaPass(c)
			return common.RespSuccess(c)
		} else {
			// 验证失败
			common.OnCaptchaFail(c)
			return common.RespError(c, i18n.T("Verification failed"), common.Map{
				"reason": reason,
			})
		}
	}

	// ===========
	//  图片验证码
	// ===========
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
