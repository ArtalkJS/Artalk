package http

import (
	"github.com/ArtalkJS/ArtalkGo/internal/entity"
	"github.com/ArtalkJS/ArtalkGo/internal/query"
	"github.com/ArtalkJS/ArtalkGo/internal/utils"
	"github.com/labstack/echo/v4"
)

type ParamsAdminSiteGet struct {
}

func (a *action) AdminSiteGet(c echo.Context) error {
	var p ParamsAdminSiteGet
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	allSites := query.FindAllSitesCooked()
	sites := allSites

	// 非超级管理员仅显示分配的站点
	if !GetIsSuperAdmin(c) {
		sites = []entity.CookedSite{}
		user := GetUserByReq(c)
		userCooked := query.CookUser(&user)
		for _, s := range allSites {
			if utils.ContainsStr(userCooked.SiteNames, s.Name) {
				sites = append(sites, s)
			}
		}
	}

	return RespData(c, Map{
		"sites": sites,
	})
}
