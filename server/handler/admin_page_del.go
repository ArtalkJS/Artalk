package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

// @Summary      Delete Page
// @Description  Delete a specific page
// @Tags         Page
// @Param        id  path  int  true  "The page ID you want to delete"
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  common.JSONResult
// @Router       /pages/{id}  [delete]
func AdminPageDel(app *core.App, router fiber.Router) {
	router.Delete("/pages/:id", func(c *fiber.Ctx) error {
		id, _ := c.ParamsInt("id")

		page := app.Dao().FindPageByID(uint(id))
		if page.IsEmpty() {
			return common.RespError(c, i18n.T("{{name}} not found", Map{"name": i18n.T("Page")}))
		}

		if !common.IsAdminHasSiteAccess(app, c, page.SiteName) {
			return common.RespError(c, i18n.T("Access denied"))
		}

		err := app.Dao().DelPage(&page)
		if err != nil {
			return common.RespError(c, i18n.T("{{name}} deletion failed", Map{"name": i18n.T("Page")}))
		}

		return common.RespSuccess(c)
	})
}
