package http

import (
	"strings"

	"github.com/ArtalkJS/ArtalkGo/internal/cache"
	"github.com/ArtalkJS/ArtalkGo/internal/entity"
	"github.com/ArtalkJS/ArtalkGo/internal/query"
	"github.com/ArtalkJS/ArtalkGo/internal/utils"
	"github.com/labstack/echo/v4"
)

type ParamsAdminSiteEdit struct {
	// 查询值
	ID uint `mapstructure:"id" param:"required"`

	// 修改值
	Name string `mapstructure:"name"`
	Urls string `mapstructure:"urls"`
}

func (a *action) AdminSiteEdit(c echo.Context) error {
	var p ParamsAdminSiteEdit
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	site := query.FindSiteByID(p.ID)
	if site.IsEmpty() {
		return RespError(c, "site 不存在")
	}

	// 站点操作权限检查
	if !IsAdminHasSiteAccess(c, site.Name) {
		return RespError(c, "无权操作")
	}

	if strings.TrimSpace(p.Name) == "" {
		return RespError(c, "site 名称不能为空白字符")
	}

	// 重命名合法性检测
	modifyName := p.Name != site.Name
	if modifyName && !query.FindSite(p.Name).IsEmpty() {
		return RespError(c, "site 已存在，请换个名称")
	}

	// urls 合法性检测
	if p.Urls != "" {
		for _, url := range utils.SplitAndTrimSpace(p.Urls, ",") {
			if !utils.ValidateURL(url) {
				return RespError(c, "Invalid url exist")
			}
		}
	}

	// 预先删除缓存，防止修改主键原有 site_name 占用问题
	cache.SiteCacheDel(&site)

	// 同步变更 site_name
	if modifyName {
		var comments []entity.Comment
		var pages []entity.Page

		a.db.Where("site_name = ?", site.Name).Find(&comments)
		a.db.Where("site_name = ?", site.Name).Find(&pages)

		for _, comment := range comments {
			comment.SiteName = p.Name
			query.UpdateComment(&comment)
		}
		for _, page := range pages {
			page.SiteName = p.Name
			query.UpdatePage(&page)
		}
	}

	// 修改 site
	site.Name = p.Name
	site.Urls = p.Urls

	err := query.UpdateSite(&site)
	if err != nil {
		return RespError(c, "site 保存失败")
	}

	return RespData(c, Map{
		"site": query.CookSite(&site),
	})
}
