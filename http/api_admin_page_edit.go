package http

import (
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsAdminPageEdit struct {
	// 查询值
	Key      string `mapstructure:"key" param:"required"`
	SiteName string `mapstructure:"site_name"`
	SiteID   uint

	// 修改值
	Url       string `mapstructure:"url"`
	Title     string `mapstructure:"title"`
	AdminOnly string `mapstructure:"admin_only"`
}

func ActionAdminPageEdit(c echo.Context) error {
	if isOK, resp := AdminOnly(c); !isOK {
		return resp
	}

	var p ParamsAdminPageEdit
	if isOK, resp := ParamsDecode(c, ParamsAdminPageEdit{}, &p); !isOK {
		return resp
	}

	// find site
	if isOK, resp := CheckSite(c, p.SiteName, &p.SiteID); !isOK {
		return resp
	}

	page := model.FindPage(p.Key, p.SiteName)
	if page.IsEmpty() {
		return RespError(c, "page not found")
	}

	// url
	if p.Url != "" {
		page.Url = p.Url
	}

	// title
	if p.Title != "" {
		page.Title = p.Title
	}

	// only_admin
	switch p.AdminOnly {
	case "1":
		page.AdminOnly = true
	case "0":
		page.AdminOnly = false
	}

	if err := model.UpdatePage(&page); err != nil {
		return RespError(c, "page save error")
	}

	return RespData(c, Map{
		"page": page.ToCooked(),
	})
}
