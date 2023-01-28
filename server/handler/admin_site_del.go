package handler

import (
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsAdminSiteDel struct {
	ID uint `form:"id" validate:"required"`
}

// @Summary      Site Delete
// @Description  Delete a specific site
// @Tags         Site
// @Param        id             formData  string  true   "the site ID you want to delete"
// @Security     ApiKeyAuth
// @Success      200  {object}  common.JSONResult
// @Router       /admin/site-del  [post]
func AdminSiteDel(router fiber.Router) {
	router.Post("/site-del", func(c *fiber.Ctx) error {
		var p ParamsAdminSiteDel
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		if !common.GetIsSuperAdmin(c) {
			return common.RespError(c, i18n.T("Access denied"))
		}

		site := query.FindSiteByID(p.ID)
		if site.IsEmpty() {
			return common.RespError(c, i18n.T("{{name}} not found", Map{"name": i18n.T("Site")}))
		}

		err := query.DelSite(&site)
		if err != nil {
			return common.RespError(c, i18n.T("{{name}} deletion failed", Map{"name": i18n.T("Site")}))
		}

		// 刷新 CORS 可信域名
		common.ReloadCorsAllowOrigins()

		return common.RespSuccess(c)
	})
}
