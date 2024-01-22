package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ResponsePageFetch struct {
	entity.CookedPage
}

// @Id           FetchPage
// @Summary      Fetch Page Data
// @Description  Fetch the data of a specific page
// @Tags         Page
// @Security     ApiKeyAuth
// @Param        id       path  int                   true  "The page ID you want to fetch"
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponsePageFetch
// @Failure      404  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /pages/{id}/fetch  [post]
func PageFetch(app *core.App, router fiber.Router) {
	router.Post("/pages/:id/fetch", common.AdminGuard(app, func(c *fiber.Ctx) error {
		id, _ := c.ParamsInt("id")

		page := app.Dao().FindPageByID(uint(id))
		if page.IsEmpty() {
			return common.RespError(c, 404, i18n.T("{{name}} not found", Map{"name": i18n.T("Page")}))
		}

		if err := app.Dao().FetchPageFromURL(&page); err != nil {
			return common.RespError(c, 500, i18n.T("Page fetch failed")+": "+err.Error())
		}

		return common.RespData(c, ResponsePageFetch{
			CookedPage: app.Dao().CookPage(&page),
		})
	}))
}
