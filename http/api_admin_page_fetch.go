package http

import (
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsPageFetch struct {
	Key      string `mapstructure:"key" param:"required"`
	SiteName string `mapstructure:"site_name"`
	SiteID   uint
}

func ActionAdminPageFetch(c echo.Context) error {
	if isOK, resp := AdminOnly(c); !isOK {
		return resp
	}

	var p ParamsPageFetch
	if isOK, resp := ParamsDecode(c, ParamsPageFetch{}, &p); !isOK {
		return resp
	}

	// find site
	if isOK, resp := CheckSite(c, p.SiteName, &p.SiteID); !isOK {
		return resp
	}

	page := model.FindPage(p.Key, p.SiteName)
	if page.IsEmpty() {
		return RespError(c, "page not found.")
	}

	if err := page.FetchURL(); err != nil {
		return RespError(c, "page fetch error: "+err.Error())
	}

	return RespData(c, Map{
		"page": page.ToCooked(),
	})
}
