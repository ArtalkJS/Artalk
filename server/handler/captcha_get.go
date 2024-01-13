package handler

import (
	"github.com/ArtalkJS/Artalk/internal/captcha"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

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
