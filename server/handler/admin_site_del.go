package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

// @Summary      Delete Site
// @Description  Delete a specific site
// @Tags         Site
// @Security     ApiKeyAuth
// @Param        id  path  int  true   "The site ID you want to delete"
// @Produce      json
// @Success      200  {object}  common.JSONResult
// @Router       /sites/{id}  [delete]
func AdminSiteDel(app *core.App, router fiber.Router) {
	router.Delete("/sites/:id", func(c *fiber.Ctx) error {
		id, _ := c.ParamsInt("id")

		if !common.GetIsSuperAdmin(app, c) {
			return common.RespError(c, i18n.T("Access denied"))
		}

		site := app.Dao().FindSiteByID(uint(id))
		if site.IsEmpty() {
			return common.RespError(c, i18n.T("{{name}} not found", Map{"name": i18n.T("Site")}))
		}

		err := app.Dao().DelSite(&site)
		if err != nil {
			return common.RespError(c, i18n.T("{{name}} deletion failed", Map{"name": i18n.T("Site")}))
		}

		return common.RespSuccess(c)
	})
}
