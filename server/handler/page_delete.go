package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

// @Id           DeletePage
// @Summary      Delete Page
// @Description  Delete a specific page
// @Tags         Page
// @Param        id  path  int  true  "The page ID you want to delete"
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  Map{}
// @Failure      403  {object}  Map{msg=string}
// @Failure      404  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /pages/{id}  [delete]
func PageDelete(app *core.App, router fiber.Router) {
	router.Delete("/pages/:id", common.AdminGuard(app, func(c *fiber.Ctx) error {
		id, _ := c.ParamsInt("id")

		page := app.Dao().FindPageByID(uint(id))
		if page.IsEmpty() {
			return common.RespError(c, 404, i18n.T("{{name}} not found", Map{"name": i18n.T("Page")}))
		}

		err := app.Dao().DelPage(&page)
		if err != nil {
			return common.RespError(c, 500, i18n.T("{{name}} deletion failed", Map{"name": i18n.T("Page")}))
		}

		return common.RespSuccess(c)
	}))
}
