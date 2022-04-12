package http

import (
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsPV struct {
	PageKey   string `mapstructure:"page_key" param:"required"`
	PageTitle string `mapstructure:"page_title"`

	SiteName string `mapstructure:"site_name"`

	SiteID  uint
	SiteAll bool
}

func ActionPV(c echo.Context) error {
	var p ParamsPV
	if isOK, resp := ParamsDecode(c, ParamsPV{}, &p); !isOK {
		return resp
	}

	// find site
	if isOK, resp := CheckSite(c, &p.SiteName, &p.SiteID, &p.SiteAll); !isOK {
		return resp
	}

	// find page
	page := model.FindCreatePage(p.PageKey, p.PageTitle, p.SiteName)

	// ip := c.RealIP()
	// ua := c.Request().UserAgent()

	pv := model.FindPV(p.PageKey, p.SiteName)
	if pv.IsEmpty() {
		// create new pv record of page
		pv = model.PV{
			PageKey:  page.Key,
			SiteName: p.SiteName,
			Num:      1,
		}

		lib.DB.Create(&pv)
	} else {
		// +1s
		pv.Num++
		lib.DB.Save(&pv)
	}

	return RespData(c, Map{
		"pv": pv.Num,
	})
}
