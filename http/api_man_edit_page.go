package http

import "github.com/labstack/echo/v4"

type ParamsEditPage struct {
	Key       string `mapstructure:"key" param:"required"`
	Url       string `mapstructure:"url"`
	Title     string `mapstructure:"title"`
	AdminOnly string `mapstructure:"admin_only"`
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

	if err := UpdatePage(&page); err != nil {
		return RespError(c, "page save error.")
	}

	return RespData(c, Map{
		"page": page.ToCooked(),
	})
}
