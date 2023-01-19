package handler

import (
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsAdminUserDel struct {
	ID uint `form:"id" validate:"required"`
}

// POST /api/admin/user-del
func AdminUserDel(router fiber.Router) {
	router.Post("/user-del", func(c *fiber.Ctx) error {
		if !common.GetIsSuperAdmin(c) {
			return common.RespError(c, i18n.T("Access denied"))
		}

		var p ParamsAdminUserDel
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		user := query.FindUserByID(p.ID)
		if user.IsEmpty() {
			return common.RespError(c, i18n.T("{{name}} not found", Map{"name": i18n.T("User")}))
		}

		err := query.DelUser(&user)
		if err != nil {
			return common.RespError(c, i18n.T("{{name}} deletion failed", Map{"name": i18n.T("User")}))
		}
		return common.RespSuccess(c)
	})
}
