package http

import (
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

	allSites := model.GetAllCookedSites()
	sites := allSites

	if !GetIsSuperAdmin(c) {
		// 非超级管理员仅显示分配的站点
		sites = []model.CookedSite{}
		user := GetUserByReq(c).ToCooked()
		for _, s := range allSites {
			hasAccess := false
			for _, us := range user.SiteNames {
				if us == s.Name {
					hasAccess = true
					break
				}
			}
			if hasAccess {
				sites = append(sites, s)
			}
		}
	}

	return RespData(c, Map{
		"sites": sites,
	})
}
