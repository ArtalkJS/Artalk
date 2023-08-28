package middleware

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

func AdminOnlyMiddleware(app *core.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !common.CheckIsAdminReq(app, c) {
			return common.RespError(c, i18n.T("Admin access required"), common.Map{"need_login": true})
		}

		return c.Next()
	}
}
