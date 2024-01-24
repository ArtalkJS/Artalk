package handler

import (
	"github.com/ArtalkJS/Artalk/internal/captcha"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsCaptchaVerify struct {
	Value string `form:"value" json:"value" validate:"required"` // The captcha value to check
}

// @Id           VerifyCaptcha
// @Summary      Verify Captcha
// @Description  Verify user enters correct captcha code
// @Tags         Captcha
// @Param        data  body  ParamsCaptchaVerify  true  "The data to check"
// @Produce      json
// @Success      200  {object}  Map{}
// @Failure      403  {object}  Map{img_data=string}
// @Router       /captcha/verify [post]
func CaptchaVerify(app *core.App, router fiber.Router) {
	router.Post("/captcha/verify", func(c *fiber.Ctx) error {
		// handle user input
		var p ParamsCaptchaVerify
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		limiter, err := common.GetLimiter(c)
		if err != nil {
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
	})
}
