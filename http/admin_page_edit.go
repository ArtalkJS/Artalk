package http

import (
	"strings"

	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsAdminPageEdit struct {
	// 查询值
	ID       uint `mapstructure:"id"`
	SiteName string
	SiteID   uint

	// 修改值
	Key       string `mapstructure:"key"`
	Title     string `mapstructure:"title"`
	AdminOnly bool   `mapstructure:"admin_only"`
}

func (a *action) AdminPageEdit(c echo.Context) error {
	var p ParamsAdminPageEdit
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	if strings.TrimSpace(p.Key) == "" {
		return RespError(c, "page key 不能为空白字符")
	}

	// use site
	UseSite(c, &p.SiteName, &p.SiteID, nil)

	// find page
	var page = model.FindPageByID(p.ID)
	if page.IsEmpty() {
		return RespError(c, "page not found")
	}

	if !IsAdminHasSiteAccess(c, page.SiteName) {
		return RespError(c, "无权操作")
	}

	// 重命名合法性检测
	modifyKey := p.Key != page.Key
	if modifyKey && !model.FindPage(p.Key, p.SiteName).IsEmpty() {
		return RespError(c, "page 已存在，请换个 key")
	}

	// 预先删除缓存，防止修改主键原有 page_key 占用问题
	model.PageCacheDel(&page)

	page.Title = p.Title
	page.AdminOnly = p.AdminOnly
	if modifyKey {
		// 相关性数据修改
		var comments []model.Comment
		a.db.Where("page_key = ?", page.Key).Find(&comments)

		for _, comment := range comments {
			comment.PageKey = p.Key
			model.UpdateComment(&comment)
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
