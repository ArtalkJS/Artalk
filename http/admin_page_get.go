package http

import (
	"github.com/ArtalkJS/ArtalkGo/model"
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
	Total int64              `json:"total"`
	Pages []model.CookedPage `json:"pages"`
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
	query := a.db.Model(&model.Page{}).Order("created_at DESC")
	if !p.SiteAll { // 不是查的所有站点
		query = query.Where("site_name = ?", p.SiteName)
	}

	// 总共条数
	var total int64
	query.Count(&total)

	// 数据分页
	query = query.Scopes(Paginate(p.Offset, p.Limit))

	// 查找
	var pages []model.Page
	query.Find(&pages)

	var cookedPages []model.CookedPage
	for _, p := range pages {
		cookedPages = append(cookedPages, p.ToCooked())
	}

	return RespData(c, ResponseAdminPageGet{
		Pages: cookedPages,
		Total: total,
	})
}
