package http

import (
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsAdminPageGet struct {
	SiteName string `mapstructure:"site_name"`
	SiteID   uint
	SiteAll  bool
}

func ActionAdminPageGet(c echo.Context) error {
	if isOK, resp := AdminOnly(c); !isOK {
		return resp
	}

	var p ParamsAdminPageGet
	if isOK, resp := ParamsDecode(c, ParamsAdminPageGet{}, &p); !isOK {
		return resp
	}

	// find site
	if isOK, resp := CheckSite(c, &p.SiteName, &p.SiteID, &p.SiteAll); !isOK {
		return resp
	}

	var pages []model.Page
	query := lib.DB.Model(&model.Page{})
	if !p.SiteAll { // 不是查的所有站点
		query = query.Where("site_name = ?", p.SiteName)
	}
	query.Find(&pages)

	var cookedPages []model.CookedPage
	for _, p := range pages {
		cookedPages = append(cookedPages, p.ToCooked())
	}

	return RespData(c, Map{
		"pages": cookedPages,
		"sites": GetAllCookedSites(),
	})
}
