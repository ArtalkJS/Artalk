package http

import (
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsAdminPageDel struct {
	Key      string `mapstructure:"key" param:"required"`
	SiteName string `mapstructure:"site_name"`
	SiteID   uint
}

func ActionAdminPageDel(c echo.Context) error {
	if isOK, resp := AdminOnly(c); !isOK {
		return resp
	}

	var p ParamsAdminPageDel
	if isOK, resp := ParamsDecode(c, ParamsAdminPageDel{}, &p); !isOK {
		return resp
	}

	// find site
	if isOK, resp := CheckSite(c, &p.SiteName, &p.SiteID, nil); !isOK {
		return resp
	}

	page := model.FindPage(p.Key, p.SiteName)
	if page.IsEmpty() {
		return RespError(c, "page not found")
	}

	err := DelPage(&page)
	if err != nil {
		return RespError(c, "Page 删除失败")
	}

	return RespSuccess(c)
}

func DelPage(page *model.Page) error {
	err := lib.DB.Unscoped().Delete(page).Error
	if err != nil {
		return err
	}

	// 删除所有相关内容
	var comments []model.Comment
	lib.DB.Where("page_key = ? AND site_name = ?", page.Key, page.SiteName).Find(&comments)

	for _, c := range comments {
		DelComment(c.ID)
	}

	// 删除 vote
	lib.DB.Unscoped().Where(
		"target_id = ? AND (type = ? OR type = ?)",
		page.ID,
		string(model.VoteTypePageUp),
		string(model.VoteTypePageDown),
	).Delete(&model.Vote{})

	return nil
}
