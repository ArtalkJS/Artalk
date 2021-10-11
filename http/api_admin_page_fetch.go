package http

import (
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsAdminPageFetch struct {
	ID uint `mapstructure:"id" param:"required"`
}

func ActionAdminPageFetch(c echo.Context) error {
	if isOK, resp := AdminOnly(c); !isOK {
		return resp
	}

	var p ParamsAdminPageFetch
	if isOK, resp := ParamsDecode(c, ParamsAdminPageFetch{}, &p); !isOK {
		return resp
	}

	page := model.FindPageByID(p.ID)
	if page.IsEmpty() {
		return RespError(c, "page not found")
	}

	if err := page.FetchURL(); err != nil {
		return RespError(c, "page fetch error: "+err.Error())
	}

	return RespData(c, Map{
		"page": page.ToCooked(),
	})
}
