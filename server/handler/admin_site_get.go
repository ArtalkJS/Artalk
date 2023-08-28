package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsAdminSiteGet struct {
}

type ResponseAdminSiteGet struct {
	Sites []entity.CookedSite `json:"sites"`
}

// @Summary      Site List
// @Description  Get a list of sites by some conditions
// @Tags         Site
// @Security     ApiKeyAuth
// @Success      200  {object}  common.JSONResult{data=ResponseAdminSiteGet}
// @Router       /admin/site-get  [post]
func AdminSiteGet(app *core.App, router fiber.Router) {
	router.Post("/site-get", func(c *fiber.Ctx) error {
		var p ParamsAdminSiteGet
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

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
			Sites: sites,
		})
	})
}
