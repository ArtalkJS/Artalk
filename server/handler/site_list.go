package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ResponseAdminSiteGet struct {
	Sites []entity.CookedSite `json:"sites"`
	Count int                 `json:"count"`
}

// @Summary      Get Site List
// @Description  Get a list of sites by some conditions
// @Tags         Site
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  ResponseAdminSiteGet
// @Router       /sites  [get]
func AdminSiteGet(app *core.App, router fiber.Router) {
	router.Get("/sites", common.AdminGuard(app, func(c *fiber.Ctx) error {
		sites := app.Dao().FindAllSitesCooked()

		return common.RespData(c, ResponseAdminSiteGet{
			Sites: sites,
			Count: len(sites),
		})
	}))
}
