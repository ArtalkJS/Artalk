package http

import (
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsAdminSiteGet struct {
}

func (a *action) AdminSiteGet(c echo.Context) error {
	var p ParamsAdminSiteGet
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	allSites := model.FindAllSitesCooked()
	sites := allSites

	// 非超级管理员仅显示分配的站点
	if !GetIsSuperAdmin(c) {
		sites = []model.CookedSite{}
		user := GetUserByReq(c).ToCooked()
		for _, s := range allSites {
			if lib.ContainsStr(user.SiteNames, s.Name) {
				sites = append(sites, s)
			}
		}
	}

	return RespData(c, Map{
		"sites": sites,
	})
}
