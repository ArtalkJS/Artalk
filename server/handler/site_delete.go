package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

// @Id           deleteSite
// @Summary      Site Delete
// @Description  Delete a specific site
// @Tags         Site
// @Security     ApiKeyAuth
// @Param        id  path  int  true   "The site ID you want to delete"
// @Produce      json
// @Success      200  {object}  Map{}
// @Failure      403  {object}  Map{msg=string}
// @Failure      404  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /sites/{id}  [delete]
func SiteDelete(app *core.App, router fiber.Router) {
	router.Delete("/sites/:id", common.AdminGuard(app, func(c *fiber.Ctx) error {
		id, _ := c.ParamsInt("id")

		site := app.Dao().FindSiteByID(uint(id))
		if site.IsEmpty() {
			return common.RespError(c, 404, i18n.T("{{name}} not found", Map{"name": i18n.T("Site")}))
		}

		err := app.Dao().DelSite(&site)
		if err != nil {
			return common.RespError(c, 500, i18n.T("{{name}} deletion failed", Map{"name": i18n.T("Site")}))
		}

		return common.RespSuccess(c)
	}))
}
