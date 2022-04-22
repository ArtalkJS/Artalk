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

func (a *action) AdminSiteAdd(c echo.Context) error {
	if isOK, resp := AdminOnly(c); !isOK {
		return resp
	}

	var p ParamsAdminSiteAdd
	if isOK, resp := ParamsDecode(c, ParamsAdminSiteAdd{}, &p); !isOK {
		return resp
	}

	if p.Urls != "" {
		urls := lib.SplitAndTrimSpace(p.Urls, ",")
		for _, url := range urls {
			if !lib.ValidateURL(url) {
				return RespError(c, "Invalid url exist")
			}
		}
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
	err := model.CreateSite(&site)
	if err != nil {
		return RespError(c, "site 创建失败")
	}

	return RespData(c, Map{
		"site": site.ToCooked(),
	})
}
