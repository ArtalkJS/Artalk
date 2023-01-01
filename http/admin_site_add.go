package http

import (
	"github.com/ArtalkJS/ArtalkGo/internal/config"
	"github.com/ArtalkJS/ArtalkGo/internal/entity"
	"github.com/ArtalkJS/ArtalkGo/internal/query"
	"github.com/ArtalkJS/ArtalkGo/internal/utils"
	"github.com/labstack/echo/v4"
)

type ParamsAdminSiteAdd struct {
	Name string `mapstructure:"name" param:"required"`
	Urls string `mapstructure:"urls"`
}

func (a *action) AdminSiteAdd(c echo.Context) error {
	var p ParamsAdminSiteAdd
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	if !GetIsSuperAdmin(c) {
		return RespError(c, "禁止创建站点")
	}

	if p.Urls != "" {
		urls := utils.SplitAndTrimSpace(p.Urls, ",")
		for _, url := range urls {
			if !utils.ValidateURL(url) {
				return RespError(c, "Invalid url exist")
			}
		}
	}

	if p.Name == config.ATK_SITE_ALL {
		return RespError(c, "禁止使用保留关键字作为名称")
	}

	if !query.FindSite(p.Name).IsEmpty() {
		return RespError(c, "site 已存在")
	}

	site := entity.Site{}
	site.Name = p.Name
	site.Urls = p.Urls
	err := query.CreateSite(&site)
	if err != nil {
		return RespError(c, "site 创建失败")
	}

	return RespData(c, Map{
		"site": query.CookSite(&site),
	})
}
