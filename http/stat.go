package http

import (
	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ParamsStat struct {
	Type string `mapstructure:"type" param:"required"`

	SiteName string
	PageKeys string `mapstructure:"page_keys"`

	Limit int `mapstructure:"limit"`

	SiteID  uint
	SiteAll bool
}

func (a *action) Stat(c echo.Context) error {
	var p ParamsStat
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	// use site
	UseSite(c, &p.SiteName, &p.SiteID, &p.SiteAll)

	// Limit 限定
	if p.Limit <= 0 {
		p.Limit = 5
	}
	if p.Limit > 100 {
		p.Limit = 100
	}

	// 公共查询规则
	QueryPages := func(d *gorm.DB) *gorm.DB {
		return d.Model(&model.Page{}).Where("site_name = ?", p.SiteName)
	}
	QueryComments := func(d *gorm.DB) *gorm.DB {
		return d.Model(&model.Comment{}).Where("site_name = ? AND is_pending = ?", p.SiteName, false)
	}
	QueryOrderRand := func(d *gorm.DB) *gorm.DB {
		if config.Instance.DB.Type == config.TypeSQLite {
			return d.Order("RANDOM()") // SQLite case
		} else {
			return d.Order("RAND()")
		}
	}

	switch p.Type {
	case "latest_comments":
		// 最新评论
		var comments []model.Comment
		a.db.Scopes(QueryComments).
			Order("created_at DESC").
			Limit(p.Limit).
			Find(&comments)

		return RespData(c, model.CookAllComments(comments))

	case "latest_pages":
		// 最新页面
		var pages []model.Page
		a.db.Scopes(QueryPages).
			Order("created_at DESC").
			Limit(p.Limit).
			Find(&pages)

		return RespData(c, model.CookAllPages(pages))

	case "pv_most_pages":
		// PV 数最多的页面
		var pages []model.Page
		a.db.Scopes(QueryPages).
			Order("pv DESC").
			Limit(p.Limit).
			Find(&pages)

		return RespData(c, model.CookAllPages(pages))

	case "comment_most_pages":
		// 评论数最多的页面
		var pages []model.Page
		a.db.Raw(
			"SELECT * FROM pages p WHERE p.site_name = ? ORDER BY (SELECT COUNT(*) FROM comments c WHERE c.page_key = p.key AND c.is_pending = ?) DESC LIMIT ?",
			p.SiteName, false, p.Limit,
		).Find(&pages)

		return RespData(c, model.CookAllPages(pages))

	case "page_pv":
		// 查询页面的 PV 数
		keys := lib.SplitAndTrimSpace(p.PageKeys, ",")
		pvs := map[string]int{}
		for _, k := range keys {
			page := model.FindPage(k, p.SiteName)
			if !page.IsEmpty() {
				pvs[k] = page.PV
			} else {
				pvs[k] = 0
			}
		}

		return RespData(c, pvs)

	case "site_pv":
		// 全站 PV 数
		var pv int64
		a.db.Raw("SELECT SUM(pv) FROM pages WHERE site_name = ?", p.SiteName).Row().Scan(&pv)

		return RespData(c, pv)

	case "page_comment":
		// 查询页面的评论数
		keys := lib.SplitAndTrimSpace(p.PageKeys, ",")
		counts := map[string]int64{}
		for _, k := range keys {
			var count int64
			a.db.Scopes(QueryComments).Where("page_key = ?", k).Count(&count)

			counts[k] = count
		}

		return RespData(c, counts)

	case "site_comment":
		// 全站评论数
		var count int64
		a.db.Scopes(QueryComments).Count(&count)

		return RespData(c, count)

	case "rand_comments":
		// 随机评论
		var comments []model.Comment
		a.db.Scopes(QueryComments, QueryOrderRand).
			Limit(p.Limit).
			Find(&comments)

		return RespData(c, model.CookAllComments(comments))

	case "rand_pages":
		// 随机页面
		var pages []model.Page
		a.db.Scopes(QueryPages, QueryOrderRand).
			Limit(p.Limit).
			Find(&pages)

		return RespData(c, model.CookAllPages(pages))
	}

	return RespError(c, "invalid type")
}
