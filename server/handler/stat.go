package handler

import (
	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ParamsStat struct {
	SiteName string `query:"site_name" json:"site_name" validate:"optional"` // The site name of your content scope
	PageKeys string `query:"page_keys" json:"page_keys" validate:"optional"` // multiple page keys separated by commas
	Limit    int    `query:"limit" json:"limit" validate:"optional"`         // The limit for pagination
}

type ResponseStat struct {
	Data interface{} `json:"data"`
}

// @Id           GetStats
// @Summary      Statistic
// @Description  Get the statistics of various data analysis
// @Tags         Statistic
// @Param        type        path   string      true   "The type of statistics"  Enums(latest_comments, latest_pages, pv_most_pages, comment_most_pages, page_pv, site_pv, page_comment, site_comment, rand_comments, rand_pages)
// @Param        options     query  ParamsStat  false  "The options"
// @Accept       json
// @Produce      json
// @Success      200  {object}  common.JSONResult
// @Failure      400  {object}  Map{msg=string}
// @Failure      403  {object}  Map{msg=string}
// @Failure      404  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /stats/{type}  [get]
func Stat(app *core.App, router fiber.Router) {
	router.Get("/stats/:type", func(c *fiber.Ctx) error {
		queryType := c.Params("type")

		var p ParamsStat
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// Limit 限定
		if p.Limit <= 0 {
			p.Limit = 5
		}
		if p.Limit > 100 {
			p.Limit = 100
		}

		// 公共查询规则
		QueryPages := func(d *gorm.DB) *gorm.DB {
			return d.Model(&entity.Page{}).Where("site_name = ?", p.SiteName)
		}
		QueryComments := func(d *gorm.DB) *gorm.DB {
			return d.Model(&entity.Comment{}).Where("site_name = ? AND is_pending = ?", p.SiteName, false)
		}
		QueryOrderRand := func(d *gorm.DB) *gorm.DB {
			if app.Conf().DB.Type == config.TypeSQLite {
				return d.Order("RANDOM()") // SQLite case
			} else {
				return d.Order("RAND()")
			}
		}

		switch queryType {
		case "latest_comments":
			// 最新评论
			var comments []*entity.Comment
			app.Dao().DB().Scopes(QueryComments).
				Order("created_at DESC").
				Limit(p.Limit).
				Find(&comments)

			return common.RespData(c, ResponseStat{
				Data: app.Dao().CookAllComments(comments),
			})

		case "latest_pages":
			// 最新页面
			var pages []entity.Page
			app.Dao().DB().Scopes(QueryPages).
				Order("created_at DESC").
				Limit(p.Limit).
				Find(&pages)

			return common.RespData(c, ResponseStat{
				Data: app.Dao().CookAllPages(pages),
			})

		case "pv_most_pages":
			// PV 数最多的页面
			var pages []entity.Page
			app.Dao().DB().Scopes(QueryPages).
				Order("pv DESC").
				Limit(p.Limit).
				Find(&pages)

			return common.RespData(c, ResponseStat{
				Data: app.Dao().CookAllPages(pages),
			})

		case "comment_most_pages":
			// 评论数最多的页面
			var pages []entity.Page
			app.Dao().DB().Raw(
				"SELECT * FROM pages p WHERE p.site_name = ? ORDER BY (SELECT COUNT(*) FROM comments c WHERE c.page_key = p.key AND c.is_pending = ?) DESC LIMIT ?",
				p.SiteName, false, p.Limit,
			).Find(&pages)

			return common.RespData(c, ResponseStat{
				Data: app.Dao().CookAllPages(pages),
			})

		case "page_pv":
			// 查询页面的 PV 数
			keys := utils.SplitAndTrimSpace(p.PageKeys, ",")
			pvs := map[string]int{}
			for _, k := range keys {
				page := app.Dao().FindPage(k, p.SiteName)
				if !page.IsEmpty() {
					pvs[k] = page.PV
				} else {
					pvs[k] = 0
				}
			}

			return common.RespData(c, ResponseStat{
				Data: pvs,
			})

		case "site_pv":
			// 全站 PV 数
			var pv int64
			app.Dao().DB().Raw("SELECT SUM(pv) FROM pages WHERE site_name = ?", p.SiteName).Row().Scan(&pv)

			return common.RespData(c, pv)

		case "page_comment":
			// 查询页面的评论数
			keys := utils.SplitAndTrimSpace(p.PageKeys, ",")
			counts := map[string]int64{}
			for _, k := range keys {
				var count int64
				app.Dao().DB().Scopes(QueryComments).Where("page_key = ?", k).Count(&count)

				counts[k] = count
			}

			return common.RespData(c, ResponseStat{
				Data: counts,
			})

		case "site_comment":
			// 全站评论数
			var count int64
			app.Dao().DB().Scopes(QueryComments).Count(&count)

			return common.RespData(c, ResponseStat{
				Data: count,
			})

		case "rand_comments":
			// 随机评论
			var comments []*entity.Comment
			app.Dao().DB().Scopes(QueryComments, QueryOrderRand).
				Limit(p.Limit).
				Find(&comments)

			return common.RespData(c, ResponseStat{
				Data: app.Dao().CookAllComments(comments),
			})

		case "rand_pages":
			// 随机页面
			var pages []entity.Page
			app.Dao().DB().Scopes(QueryPages, QueryOrderRand).
				Limit(p.Limit).
				Find(&pages)

			return common.RespData(c, ResponseStat{
				Data: app.Dao().CookAllPages(pages),
			})
		}

		return common.RespError(c, 404, i18n.T("Invalid {{name}}", Map{"name": i18n.T("Type")}))
	})
}
