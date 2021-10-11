package http

import (
	"strings"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsAdminPageEdit struct {
	// 查询值
	ID       uint   `mapstructure:"id" param:"required"`
	SiteName string `mapstructure:"site_name"`
	SiteID   uint

	// 修改值
	Key       string `mapstructure:"key"`
	Title     string `mapstructure:"title"`
	AdminOnly bool   `mapstructure:"admin_only"`
}

func ActionAdminPageEdit(c echo.Context) error {
	if isOK, resp := AdminOnly(c); !isOK {
		return resp
	}

	var p ParamsAdminPageEdit
	if isOK, resp := ParamsDecode(c, ParamsAdminPageEdit{}, &p); !isOK {
		return resp
	}

	if strings.TrimSpace(p.Key) == "" {
		return RespError(c, "page key 不能为空白字符")
	}

	// find site
	if isOK, resp := CheckSite(c, &p.SiteName, &p.SiteID, nil); !isOK {
		return resp
	}

	page := model.FindPageByID(p.ID)
	if page.IsEmpty() {
		return RespError(c, "page not found")
	}

	// 重命名合法性检测
	modifyKey := p.Key != page.Key
	if modifyKey && !model.FindPage(p.Key, p.SiteName).IsEmpty() {
		return RespError(c, "page 已存在，请换个 key")
	}

	page.Title = p.Title
	page.AdminOnly = p.AdminOnly
	if modifyKey {
		// 相关性数据修改
		var comments []model.Comment
		lib.DB.Where("page_key = ?", page.Key).Find(&comments)

		for _, comment := range comments {
			comment.PageKey = p.Key
			lib.DB.Save(&comment)
		}

		page.Key = p.Key
	}

	if err := model.UpdatePage(&page); err != nil {
		return RespError(c, "page save error")
	}

	return RespData(c, Map{
		"page": page.ToCooked(),
	})
}
