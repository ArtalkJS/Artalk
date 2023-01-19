package handler

import (
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsAdminSiteGet struct {
}

// POST /api/admin/site-get
func AdminSiteGet(router fiber.Router) {
	router.Post("/site-get", func(c *fiber.Ctx) error {
		var p ParamsAdminSiteGet
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		allSites := query.FindAllSitesCooked()
		sites := allSites

		// 非超级管理员仅显示分配的站点
		if !common.GetIsSuperAdmin(c) {
			sites = []entity.CookedSite{}
			user := common.GetUserByReq(c)
			userCooked := query.CookUser(&user)
			for _, s := range allSites {
				if utils.ContainsStr(userCooked.SiteNames, s.Name) {
					sites = append(sites, s)
				}
			}
		}

		return common.RespData(c, common.Map{
			"sites": sites,
		})
	})
}
