package http

import (
	"strings"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsAdminSiteEdit struct {
	// 查询值
	ID uint `mapstructure:"id" param:"required"`

	// 修改值
	Name string `mapstructure:"name"`
	Urls string `mapstructure:"urls"`
}

func ActionAdminSiteEdit(c echo.Context) error {
	if isOK, resp := AdminOnly(c); !isOK {
		return resp
	}

	var p ParamsAdminSiteEdit
	if isOK, resp := ParamsDecode(c, ParamsAdminSiteEdit{}, &p); !isOK {
		return resp
	}

	site := model.FindSiteByID(p.ID)
	if site.IsEmpty() {
		return RespError(c, "site 不存在")
	}

	if strings.TrimSpace(p.Name) == "" {
		return RespError(c, "site 名称不能为空白字符")
	}

	// 重命名合法性检测
	modifyName := p.Name != site.Name
	if modifyName && !model.FindSite(p.Name).IsEmpty() {
		return RespError(c, "site 已存在，请换个名称")
	}

	// urls 合法性检测
	if p.Urls != "" {
		for _, url := range site.ToCooked().Urls {
			if !lib.ValidateURL(url) {
				return RespError(c, "Invalid url exist")
			}
		}
	}

	// 同步变更 site_name
	if modifyName {
		var comments []model.Comment
		lib.DB.Where("site_name = ?", site.Name).Find(&comments)

		var pages []model.Page
		lib.DB.Where("site_name = ?", site.Name).Find(&pages)

		for _, comment := range comments {
			comment.SiteName = p.Name
			lib.DB.Save(&comment)
		}
		for _, page := range pages {
			page.SiteName = p.Name
			lib.DB.Save(&page)
		}
	}

	// 修改 site
	site.Name = p.Name
	site.Urls = p.Urls

	err := lib.DB.Save(&site).Error
	if err != nil {
		return RespError(c, "site 保存失败")
	}

	return RespData(c, Map{
		"site": site.ToCooked(),
	})
}
