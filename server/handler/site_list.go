package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ResponseSiteList struct {
	Sites []entity.CookedSite `json:"sites"`
	Count int                 `json:"count"`
}

// @Id           GetSites
// @Summary      Get Site List
// @Description  Get a list of sites by some conditions
// @Tags         Site
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  ResponseSiteList
// @Router       /sites  [get]
func SiteList(app *core.App, router fiber.Router) {
	router.Get("/sites", common.AdminGuard(app, func(c *fiber.Ctx) error {
		sites := app.Dao().FindAllSitesCooked()

		return common.RespData(c, ResponseSiteList{
			Sites: sites,
			Count: len(sites),
		})
	}))
}
