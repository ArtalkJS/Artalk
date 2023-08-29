package limiter

import (
	"path"

	"github.com/ArtalkJS/Artalk/internal/captcha"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/limiter"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ActionLimitConf struct {
	ProtectPaths []string
}

const LimiterLocalKey = "limiter"

// 操作限制 中间件
func ActionLimitMiddleware(app *core.App, conf ActionLimitConf) fiber.Handler {
	limiter := limiter.NewLimiter(app, &limiter.LimiterConf{
		AlwaysMode:          app.Conf().Captcha.Always,
		MaxActionDuringTime: app.Conf().Captcha.ActionLimit,
		ResetTimeout:        app.Conf().Captcha.ActionReset,
	})

	return func(c *fiber.Ctx) error {
		// 任何页面都保存 limiter 实例
		c.Locals(LimiterLocalKey, limiter)

		// 关闭验证码功能，直接 Skip
		if !app.Conf().Captcha.Enabled {
			return c.Next()
		}

		// 路径跳过
		if !isProtectPath(c.Path(), conf.ProtectPaths) {
			return c.Next()
		}

		// 管理员直接忽略
		if common.CheckIsAdminReq(app, c) {
			return c.Next()
		}

		// 检测是否需要验证码
		ip := c.IP()
		if limiter.IsPass(ip) {
			// 无需验证码
			err := c.Next()

			// 若为保护路径
			if isProtectPath(c.Path(), conf.ProtectPaths) {
				limiter.Log(ip)
			}

			return err
		} else {
			// create new captcha checker instance
			cap := common.NewCaptchaChecker(app, c)

			// response need captcha check
			respData := common.Map{
				"need_captcha": true,
			}

			switch cap.Type() {
			case captcha.Image:
				// 图片验证码
				img, _ := cap.Get()
				respData["img_data"] = img
			case captcha.IFrame:
				// iFrame 验证模式
				respData["iframe"] = true
				// 前端新版不会再用到 img_data，向旧版响应 ArtalkFrontend out-of-date 图片
				respData["img_data"] = "data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 160 40'%3E%3Cdefs%3E%3Cstyle%3E.a%7Bfill:%23328ce6%3B%7D.b%7Bfont-size:12px%3Bfill:%23fff%3Bfont-family:sans-serif%3B%7D%3C/style%3E%3C/defs%3E%3Crect class='a' width='160' height='40'/%3E%3Ctext class='b' transform='translate(18.37 16.67)'%3EArtalk Frontend%3Ctspan x='0' y='14.4'%3EOut-Of-Date.%3C/tspan%3E%3C/text%3E%3C/svg%3E"
			}

			return common.RespError(c, i18n.T("Captcha required"), respData)
		}
	}
}

func isProtectPath(pathTarget string, protectPaths []string) bool {
	for _, p := range protectPaths {
		if path.Clean(pathTarget) == path.Clean(p) {
			return true
		}
	}
	return false
}
