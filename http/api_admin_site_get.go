package http

import (
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsAdminSiteGet struct {
}

func (a *action) AdminSiteGet(c echo.Context) error {
	var p ParamsAdminSiteGet
	if isOK, resp := ParamsDecode(c, ParamsAdminSiteGet{}, &p); !isOK {
		return resp
	}

	return RespData(c, Map{
		"sites": model.GetAllCookedSites(),
	})
}
