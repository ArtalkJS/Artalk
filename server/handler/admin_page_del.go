package handler

import (
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsAdminPageDel struct {
	Key      string `form:"key" validate:"required"`
	SiteName string
	SiteID   uint
}

// @Summary      Page Delete
// @Description  Delete a specific page
// @Tags         Page
// @Param        key            formData  string  true   "the page KEY you want to delete"
// @Param        site_name      formData  string  false  "the site name of your content scope"
// @Security     ApiKeyAuth
// @Success      200  {object}  common.JSONResult
// @Router       /admin/page-del  [post]
func AdminPageDel(router fiber.Router) {
	router.Post("/page-del", func(c *fiber.Ctx) error {
		var p ParamsAdminPageDel
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// use site
		common.UseSite(c, &p.SiteName, &p.SiteID, nil)

		page := query.FindPage(p.Key, p.SiteName)
		if page.IsEmpty() {
			return common.RespError(c, i18n.T("{{name}} not found", Map{"name": i18n.T("Page")}))
		}

		if !common.IsAdminHasSiteAccess(c, page.SiteName) {
			return common.RespError(c, i18n.T("Access denied"))
		}

		err := query.DelPage(&page)
		if err != nil {
			return common.RespError(c, i18n.T("{{name}} deletion failed", Map{"name": i18n.T("Page")}))
		}

		return common.RespSuccess(c)
	})
}
