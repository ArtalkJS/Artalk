package common

import (
	"fmt"

	"github.com/ArtalkJS/Artalk/internal/captcha"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/gofiber/fiber/v2"
)

func GetLimiter[T any](c *fiber.Ctx) (lmt *T, error error) {
	l := c.Locals("limiter").(*T)
	if l == nil {
		return nil, RespError(c, "limiter is not initialize")
	}
	return l, nil
}

func NewCaptchaChecker(app *core.App, c *fiber.Ctx) captcha.Checker {
	return captcha.NewCaptchaChecker(&captcha.CheckerConf{
		CaptchaConf: app.Conf().Captcha,
		User: captcha.User{
			ID: fmt.Sprint(GetUserByReq(app, c).ID),
			IP: c.IP(),
		},
	})
}
