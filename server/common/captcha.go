package common

import (
	"fmt"

	"github.com/ArtalkJS/Artalk/internal/captcha"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/server/middleware/limiter"
	"github.com/gofiber/fiber/v2"
)

func GetLimiter(c *fiber.Ctx) (lmt *limiter.Limiter, err error) {
	l, ok := c.Locals("limiter").(*limiter.Limiter)
	if l == nil || !ok {
		return nil, RespError(c, 500, "limiter is not initialize, but middleware is used")
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

func LimiterGuard(app *core.App, handler fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		limiter, err := GetLimiter(c)
		if err != nil {
			return err
		}

		// 关闭验证码功能，直接 Skip
		if !app.Conf().Captcha.Enabled {
			return handler(c)
		}

		// 管理员直接忽略
		if CheckIsAdminReq(app, c) {
			return handler(c)
		}

		// 检测是否需要验证码
		ip := c.IP()
		if limiter.IsPass(ip) {
			// 无需验证码
			err := handler(c)

			if c.Method() != fiber.MethodOptions { // 忽略 Options 请求
				limiter.Log(ip) // 记录操作
			}

			return err
		} else {
			// create new captcha checker instance
			cap := NewCaptchaChecker(app, c)

			// response need captcha check
			respData := Map{
				"need_captcha": true,
			}

			switch cap.Type() {
			case captcha.Image:
				// 图片验证码
				img, _ := cap.Get()
				respData["img_data"] = string(img)
			case captcha.IFrame:
				// iFrame 验证模式
				respData["iframe"] = true
			}

			return RespError(c, 403, i18n.T("Captcha required"), respData)
		}
	}
}
