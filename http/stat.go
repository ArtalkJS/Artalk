package http

import (
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsStat struct {
	Type string `mapstructure:"type" param:"required"`

	SiteName string `mapstructure:"site_name" param:"required"`
	PageKeys string `mapstructure:"page_keys"`

	Limit int `mapstructure:"limit"`
}

func (a *action) Stat(c echo.Context) error {
	var p ParamsStat
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	if p.Limit <= 0 {
		p.Limit = 5
	}

	switch p.Type {
	case "latest_comments":
		// 最新评论
		var comments []model.Comment
		a.db.Model(&model.Comment{}).
			Where("site_name = ? AND is_pending = ?", p.SiteName, false).
			Order("created_at DESC").
			Limit(p.Limit).
			Find(&comments)

		return RespData(c, model.CookAllComments(comments))

	case "latest_pages":
		// 最新页面
		var pages []model.Page
		a.db.Model(&model.Page{}).
			Where("site_name = ?", p.SiteName).
			Order("created_at DESC").
			Limit(p.Limit).
			Find(&pages)

		return RespData(c, model.CookAllPages(pages))

	case "pv_most_pages":
		// PV 数最多的页面
		var pages []model.Page
		a.db.Model(&model.Page{}).
			Where("site_name = ?", p.SiteName).
			Order("pv DESC").
			Limit(p.Limit).
			Find(&pages)

		return RespData(c, model.CookAllPages(pages))

	case "comment_most_pages":
		// 评论数最多的页面
		var pages []model.Page
		a.db.Raw("SELECT * FROM pages p WHERE p.site_name = ? ORDER BY (SELECT COUNT(*) FROM comments c WHERE c.page_key = p.key AND c.is_pending = ?) DESC", p.SiteName, false).
			Find(&pages)

		return RespData(c, model.CookAllPages(pages))

	case "page_pv":
		// 查询页面的 PV 数
		keys := lib.SplitAndTrimSpace(p.PageKeys, ",")
		pvs := map[string]int{}
		for _, k := range keys {
			page := model.FindPage(k, p.SiteName)
			if !page.IsEmpty() {
				pvs[page.Key] = page.PV
			}
		}

		return RespData(c, pvs)

	case "site_pv":
		// 全站 PV 数
		var pv int64
		a.db.Model(&model.Page{}).
			Where("site_name = ?", p.SiteName).
			Count(&pv)

		return RespData(c, pv)

	case "page_comment":
		// 查询页面的评论数
		keys := lib.SplitAndTrimSpace(p.PageKeys, ",")
		counts := map[string]int64{}
		for _, k := range keys {
			var count int64
			a.db.Model(&model.Comment{}).
				Where("page_key = ? AND site_name = ? AND is_pending = ?", k, p.SiteName, false).
				Count(&count)

			counts[k] = count
		}

		return RespData(c, counts)

	case "site_comment":
		// 全站评论数
		var count int64
		a.db.Model(&model.Comment{}).
			Where("site_name = ? AND is_pending = ?", p.SiteName, false).
			Count(&count)

		return RespData(c, count)
	}

	return RespError(c, "invalid type")
}
