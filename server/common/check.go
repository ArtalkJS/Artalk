package common

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/gofiber/fiber/v2"
)

func CheckIsAllowed(app *core.App, c *fiber.Ctx, name string, email string, page entity.Page, siteName string) (bool, error) {
	isAdminUser := app.Dao().IsAdminUserByNameEmail(name, email)

	// 如果用户是管理员，或者当前页只能管理员评论
	if isAdminUser || page.AdminOnly {
		if !CheckIsAdminReq(app, c) {
			return false, RespError(c, 403, i18n.T("Admin access required"), Map{"need_login": true})
		}
	}

	return true, nil
}

func AdminRequired(app *core.App, c *fiber.Ctx) (bool, error) {
	if !CheckIsAdminReq(app, c) {
		return false, RespError(c, 403, i18n.T("Admin access required"), Map{"need_login": true})
	}

	return true, nil
}

func AdminGuard(app *core.App, handler fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if ok, resp := AdminRequired(app, c); !ok {
			return resp
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
