package http

import (
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsAdminSiteDel struct {
	ID         uint `mapstructure:"id" param:"required"`
	DelContent bool `mapstructure:"del_content"`
}

func ActionAdminSiteDel(c echo.Context) error {
	if isOK, resp := AdminOnly(c); !isOK {
		return resp
	}

	var p ParamsAdminSiteDel
	if isOK, resp := ParamsDecode(c, ParamsAdminSiteDel{}, &p); !isOK {
		return resp
	}

	site := model.FindSiteByID(p.ID)
	if site.IsEmpty() {
		return RespError(c, "site 不存在")
	}

	err := lib.DB.Unscoped().Delete(&site).Error
	if err != nil {
		return RespError(c, "site 删除失败")
	}

	// 删除所有相关内容
	if p.DelContent {
		var comments []model.Comment
		lib.DB.Where("site_name = ?", site.Name).Find(&comments)

		var pages []model.Page
		lib.DB.Where("site_name = ?", site.Name).Find(&pages)

		tx := lib.DB.Begin()
		for _, c := range comments {
			tx.Delete(&c)
		}
		for _, p := range pages {
			tx.Unscoped().Delete(&p)
		}
		tx.Commit()
	}

	return RespSuccess(c)
}
