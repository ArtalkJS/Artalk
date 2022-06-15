package http

import (
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsPV struct {
	PageKey   string `mapstructure:"page_key" param:"required"`
	PageTitle string `mapstructure:"page_title"`

	SiteName string
	SiteID   uint
	SiteAll  bool
}

func (a *action) PV(c echo.Context) error {
	var p ParamsPV
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	// use site
	UseSite(c, &p.SiteName, &p.SiteID, &p.SiteAll)

	// find page
	page := model.FindCreatePage(p.PageKey, p.PageTitle, p.SiteName)

	// ip := c.RealIP()
	// ua := c.Request().UserAgent()

	page.PV++
	model.UpdatePage(&page)

	return RespData(c, Map{
		"pv": page.PV,
	})
}
