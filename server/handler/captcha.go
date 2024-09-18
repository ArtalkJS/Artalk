package handler

import (
	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/gofiber/fiber/v2"
)

func Captcha(app *core.App, router fiber.Router) {
	router.Group("/captcha", func(c *fiber.Ctx) error {
		if !app.Conf().Captcha.Enabled {
			return common.RespError(c, 404, "Captcha disabled")
		}
		return c.Next()
	})
	{
		CaptchaGet(app, router)
		CaptchaStatus(app, router)
		CaptchaVerify(app, router)
	}
}
