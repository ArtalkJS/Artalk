package handler

import (
	"github.com/ArtalkJS/Artalk/internal/captcha"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/limiter"
	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

func Captcha(app *core.App, router fiber.Router) {
	ca := router.Group("/captcha/", func(c *fiber.Ctx) error {
		if !app.Conf().Captcha.Enabled {
			return common.RespError(c, 404, "Captcha disabled")
		}
		return c.Next()
	})
	{
		ca.Get("/status", captchaStatus(app))
		ca.Get("/get", captchaGet(app))
		ca.Post("/verify", captchaVerify(app))
	}
}

type ResponseCaptchaStatus struct {
	IsPass bool `json:"is_pass"`
}

// @Summary      Get Captcha Status
// @Description  Get the status of the user's captcha verification
// @Tags         Captcha
// @Produce      json
// @Success      200  {object}  ResponseCaptchaStatus
// @Router       /captcha/status  [get]
func captchaStatus(app *core.App) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		limiter, err := common.GetLimiter[limiter.Limiter](c)
		if limiter == nil {
			return err
		}

		return common.RespData(c, ResponseCaptchaStatus{
			IsPass: limiter.IsPass(c.IP()),
		})
	}
}

type ResponseCaptchaGet struct {
	ImgData string `json:"img_data"`
}

// @Summary      Get Captcha
// @Description  Get a base64 encoded captcha image or a HTML page to verify for user
// @Tags         Captcha
// @Produce      json,html
// @Success      200  {object}  ResponseCaptchaGet
// @Failure      500  {object}  Map{msg=string}
// @Router       /captcha/get  [get]
func captchaGet(app *core.App) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// create new captcha checker instance
		cap := common.NewCaptchaChecker(app, c)

		// get new captcha
		got, err := cap.Get()
		if err != nil {
			return common.RespError(c, 500, "captcha generate err: "+err.Error())
		}

		// response captcha
		switch cap.Type() {
		case captcha.Image:
			return common.RespData(c, ResponseCaptchaGet{
				ImgData: string(got),
			})

		case captcha.IFrame:
			c.Set(fiber.HeaderCacheControl, "no-cache, no-store, must-revalidate") // disable cache
			c.Set(fiber.HeaderPragma, "no-cache")
			c.Set(fiber.HeaderExpires, "0")
			c.Set(fiber.HeaderContentType, "text/html") // response html body
			return c.Send(got)
		}

		return common.RespError(c, 500, "invalid captcha type")
	}
}

type ParamsCaptchaVerify struct {
	Value string `form:"value" json:"value" validate:"required"` // The captcha value to check
}

// @Summary      Verify Captcha
// @Description  Verify user enters correct captcha code
// @Tags         Captcha
// @Param        data  body  ParamsCaptchaVerify  true  "The data to check"
// @Produce      json
// @Success      200  {object}  Map{}
// @Failure      403  {object}  Map{img_data=string}
// @Router       /captcha/verify [post]
func captchaVerify(app *core.App) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// handle user input
		var p ParamsCaptchaVerify
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		limiter, err := common.GetLimiter[limiter.Limiter](c)
		if limiter == nil {
			return err
		}

		// create new captcha checker instance
		cap := common.NewCaptchaChecker(app, c)

		// check user input
		isPass, err := cap.Check(p.Value)

		// trigger global event
		if isPass {
			// 验证码正确
			limiter.MarkVerifyPassed(c.IP())
		} else {
			// 验证码错误
			limiter.MarkVerifyFailed(c.IP())
		}

		// response result
		switch cap.Type() {
		case captcha.Image:
			if !isPass {
				img, err := cap.Get()
				if err != nil {
					return common.RespError(c, 500, "captcha generate err: "+err.Error())
				}

				return common.RespError(c, 403, i18n.T("Wrong captcha"), common.Map{
					"img_data": string(img),
				})
			}
		case captcha.IFrame:
			if !isPass {
				log.Error("[Captcha] Failed to verify: ", err)
				return common.RespError(c, 403, i18n.T("Verification failed"), common.Map{
					"detail": err.Error(),
				})
			}
		}

		return common.RespSuccess(c)
	}
}
