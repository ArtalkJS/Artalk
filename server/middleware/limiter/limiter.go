package limiter

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/limiter"
	"github.com/gofiber/fiber/v2"
)

type ActionLimitConf struct {
}

const LimiterLocalKey = "limiter"

type Limiter = limiter.Limiter

// 操作限制 中间件
func ActionLimitMiddleware(app *core.App, conf ActionLimitConf) fiber.Handler {
	limiter := limiter.NewLimiter(&limiter.LimiterConf{
		AlwaysMode:          app.Conf().Captcha.Always,
		MaxActionDuringTime: app.Conf().Captcha.ActionLimit,
		ResetTimeout:        app.Conf().Captcha.ActionReset,
	})

	return func(c *fiber.Ctx) error {
		// 任何页面都保存 limiter 实例
		c.Locals(LimiterLocalKey, limiter)

		return c.Next()
	}
}
