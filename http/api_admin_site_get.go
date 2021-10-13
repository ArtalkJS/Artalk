package http

import (
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsAdminSiteGet struct {
}

func ActionAdminSiteGet(c echo.Context) error {
	if isOK, resp := AdminOnly(c); !isOK {
		return resp
	}

	var p ParamsAdminSiteGet
	if isOK, resp := ParamsDecode(c, ParamsAdminSiteGet{}, &p); !isOK {
		return resp
	}

	return RespData(c, Map{
		"sites": GetAllCookedSites(),
	})
}

func GetAllCookedSites() []model.CookedSite {
	var sites []model.Site
	lib.DB.Model(&model.Site{}).Find(&sites)

	var cookedSites []model.CookedSite
	for _, s := range sites {
		cookedSites = append(cookedSites, s.ToCooked())
	}

	return cookedSites
}
