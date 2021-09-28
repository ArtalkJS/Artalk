package http

import (
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsAdminSiteEdit struct {
	ID   uint   `mapstructure:"id" param:"required"`
	Name string `mapstructure:"name"`
	Url  string `mapstructure:"url"`
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

	modifyName := (p.Name != "" && p.Name != site.Name)
	if modifyName && !model.FindSite(p.Name).IsEmpty() {
		return RespError(c, "site 已存在，请换个名称")
	}

	site.Name = p.Name
	site.Url = p.Url

	err := lib.DB.Save(&site).Error
	if err != nil {
		return RespError(c, "site 保存失败")
	}

	// 同步变更 site_name
	if modifyName {
		var comments []model.Comment
		lib.DB.Where("site_name = ?", p.Name).Find(&comments)

		var pages []model.Page
		lib.DB.Where("site_name = ?", p.Name).Find(&pages)

		tx := lib.DB.Begin()
		for _, c := range comments {
			c.SiteName = site.Name
			tx.Save(&c)
		}
		for _, p := range pages {
			p.SiteName = site.Name
			tx.Save(&p)
		}
		tx.Commit()
	}

	return RespSuccess(c)
}
