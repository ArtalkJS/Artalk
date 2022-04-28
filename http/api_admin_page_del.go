package http

import (
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsAdminPageDel struct {
	Key      string `mapstructure:"key" param:"required"`
	SiteName string `mapstructure:"site_name"`
	SiteID   uint
}

func (a *action) AdminPageDel(c echo.Context) error {
	var p ParamsAdminPageDel
	if isOK, resp := ParamsDecode(c, ParamsAdminPageDel{}, &p); !isOK {
		return resp
	}

	// find site
	if isOK, resp := AdminSiteInControl(c, &p.SiteName, &p.SiteID, nil); !isOK {
		return resp
	}

	page := model.FindPage(p.Key, p.SiteName)
	if page.IsEmpty() {
		return RespError(c, "page not found")
	}

	err := model.DelPage(&page)
	if err != nil {
		return RespError(c, "Page 删除失败")
	}

	return RespSuccess(c)
}
