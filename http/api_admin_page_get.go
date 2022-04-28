package http

import (
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsAdminPageGet struct {
	SiteName string `mapstructure:"site_name"`
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
	if isOK, resp := ParamsDecode(c, ParamsAdminPageGet{}, &p); !isOK {
		return resp
	}

	// find site
	if isOK, resp := AdminSiteInControl(c, &p.SiteName, &p.SiteID, &p.SiteAll); !isOK {
		return resp
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
