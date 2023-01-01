package http

import (
	"github.com/ArtalkJS/ArtalkGo/internal/entity"
	"github.com/ArtalkJS/ArtalkGo/internal/query"
	"github.com/labstack/echo/v4"
)

type ParamsAdminPageGet struct {
	SiteName string
	SiteID   uint
	SiteAll  bool
	Limit    int `mapstructure:"limit"`
	Offset   int `mapstructure:"offset"`
}

type ResponseAdminPageGet struct {
	Total int64               `json:"total"`
	Pages []entity.CookedPage `json:"pages"`
}

func (a *action) AdminPageGet(c echo.Context) error {
	var p ParamsAdminPageGet
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	// use site
	UseSite(c, &p.SiteName, &p.SiteID, &p.SiteAll)

	if !IsAdminHasSiteAccess(c, p.SiteName) {
		return RespError(c, "无权操作")
	}

	// 准备 query
	db := a.db.Model(&entity.Page{}).Order("created_at DESC")
	if !p.SiteAll { // 不是查的所有站点
		db = db.Where("site_name = ?", p.SiteName)
	}

	// 总共条数
	var total int64
	db.Count(&total)

	// 数据分页
	db = db.Scopes(Paginate(p.Offset, p.Limit))

	// 查找
	var pages []entity.Page
	db.Find(&pages)

	var cookedPages []entity.CookedPage
	for _, p := range pages {
		cookedPages = append(cookedPages, query.CookPage(&p))
	}

	return RespData(c, ResponseAdminPageGet{
		Pages: cookedPages,
		Total: total,
	})
}
