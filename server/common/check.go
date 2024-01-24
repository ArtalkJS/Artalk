package common

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/gofiber/fiber/v2"
)

func AdminGuard(app *core.App, handler fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !CheckIsAdminReq(app, c) {
			return RespError(c, 403, i18n.T("Admin access required"), Map{"need_login": true})
		}

		return handler(c)
	}
}

func CheckIsAdminReq(app *core.App, c *fiber.Ctx) bool {
	jwt := GetJwtInstanceByReq(app, c)
	if jwt == nil {
		return false
	}

	user := GetUserByJwt(app, jwt)
	return user.IsAdmin
}
