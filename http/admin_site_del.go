package http

import (
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsAdminSiteDel struct {
	ID uint `mapstructure:"id" param:"required"`
}

func (a *action) AdminSiteDel(c echo.Context) error {
	var p ParamsAdminSiteDel
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	if !GetIsSuperAdmin(c) {
		return RespError(c, "禁止删除站点")
	}

	site := model.FindSiteByID(p.ID)
	if site.IsEmpty() {
		return RespError(c, "site 不存在")
	}

	err := model.DelSite(&site)
	if err != nil {
		return RespError(c, "site 删除失败")
	}

	return RespSuccess(c)
}
