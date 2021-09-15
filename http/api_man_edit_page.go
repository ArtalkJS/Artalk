package http

import "github.com/labstack/echo/v4"

type ParamsEditPage struct {
	Key       string `mapstructure:"key" param:"required"`
	OnlyAdmin string `mapstructure:"only_admin"`
}

func ActionManagerEditPage(c echo.Context) error {
	if isOK, resp := AdminOnly(c); !isOK {
		return resp
	}

	var p ParamsEditPage
	if isOK, resp := ParamsDecode(c, ParamsEditPage{}, &p); !isOK {
		return resp
	}

	page := FindPage(p.Key)
	if page.IsEmpty() {
		return RespError(c, "page not found.")
	}

	switch p.OnlyAdmin {
	case "1":
		page.AdminOnly = true
	case "0":
		page.AdminOnly = false
	}

	if err := UpdatePage(&page); err != nil {
		return RespError(c, "page save error.")
	}

	return RespData(c, Map{
		"page": page.ToCooked(),
	})
}
