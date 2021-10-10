package http

import (
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsAdminSiteAdd struct {
	Name string `mapstructure:"name" param:"required"`
	Urls string `mapstructure:"urls"`
}

func ActionAdminSiteAdd(c echo.Context) error {
	if isOK, resp := AdminOnly(c); !isOK {
		return resp
	}

	var p ParamsAdminSiteAdd
	if isOK, resp := ParamsDecode(c, ParamsAdminSiteAdd{}, &p); !isOK {
		return resp
	}

	if p.Urls != "" && !lib.ValidateURL(p.Urls) {
		return RespError(c, "Invalid url")
	}

	if p.Name == lib.ATK_SITE_ALL {
		return RespError(c, "禁止使用保留关键字作为名称")
	}

	if !model.FindSite(p.Name).IsEmpty() {
		return RespError(c, "site 已存在")
	}

	site := model.Site{}
	site.Name = p.Name
	site.Urls = p.Urls
	err := lib.DB.Create(&site).Error
	if err != nil {
		return RespError(c, "site 创建失败")
	}

	return RespSuccess(c)
}
