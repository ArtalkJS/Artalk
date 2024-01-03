package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ResponseAdminSiteGet struct {
	Data []entity.CookedSite `json:"data"`
}

// @Summary      Get Site List
// @Description  Get a list of sites by some conditions
// @Tags         Site
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  ResponseAdminSiteGet
// @Router       /sites  [get]
func AdminSiteGet(app *core.App, router fiber.Router) {
	router.Get("/sites", func(c *fiber.Ctx) error {
		allSites := app.Dao().FindAllSitesCooked()
		sites := allSites

		// 非超级管理员仅显示分配的站点
		if !common.GetIsSuperAdmin(app, c) {
			sites = []entity.CookedSite{}
			user := common.GetUserByReq(app, c)
			userCooked := app.Dao().CookUser(&user)
			for _, s := range allSites {
				if utils.ContainsStr(userCooked.SiteNames, s.Name) {
					sites = append(sites, s)
				}
			}
		}

		return common.RespData(c, ResponseAdminSiteGet{
			Data: sites,
		})
	})
}
