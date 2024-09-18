package handler

import (
	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/i18n"
	"github.com/artalkjs/artalk/v2/internal/utils"
	"github.com/artalkjs/artalk/v2/server/common"
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

		// Limit parameter
		if p.Limit <= 0 {
			p.Limit = 5
		}
		if p.Limit > 100 {
			// The maximum limit is 100 for data security reasons
			p.Limit = 100
		}

		// Reusable query scopes
		// Query Pages by `site_name`
		QueryPages := func(d *gorm.DB) *gorm.DB {
			return d.Model(&entity.Page{}).Where(&entity.Page{SiteName: p.SiteName})
		}
		// Query Comments by `site_name` and `is_pending=false`
		QueryComments := func(d *gorm.DB) *gorm.DB {
			return d.Model(&entity.Comment{}).Where(&entity.Comment{SiteName: p.SiteName, IsPending: false})
		}
		// Query Order by RAND()
		QueryOrderRand := func(d *gorm.DB) *gorm.DB {
			if app.Conf().DB.Type == config.TypeSQLite {
				return d.Order("RANDOM()") // SQLite
			} else {
				return d.Order("RAND()")
			}
		}

		switch queryType {
		case "latest_comments":
			// ------------------------------------
			//  Latest comments
			// ------------------------------------
			var comments []*entity.Comment
			app.Dao().DB().Scopes(QueryComments).
				Order("created_at DESC").
				Limit(p.Limit).
				Find(&comments)

			return common.RespData(c, ResponseStat{
				Data: app.Dao().CookAllComments(comments),
			})

		case "latest_pages":
			// ------------------------------------
			//  Latest pages
			// ------------------------------------
			var pages []entity.Page
			app.Dao().DB().Scopes(QueryPages).
				Order("created_at DESC").
				Limit(p.Limit).
				Find(&pages)

			return common.RespData(c, ResponseStat{
				Data: app.Dao().CookAllPages(pages),
			})

		case "pv_most_pages":
			// ------------------------------------
			//  PV most pages
			// ------------------------------------
			var pages []entity.Page
			app.Dao().DB().Scopes(QueryPages).
				Order("pv DESC").
				Limit(p.Limit).
				Find(&pages)

			return common.RespData(c, ResponseStat{
				Data: app.Dao().CookAllPages(pages),
			})

		case "comment_most_pages":
			// ------------------------------------
			//  Comment most pages
			// ------------------------------------
			tbPages := app.Dao().GetTableName(&entity.Page{})
			tbComments := app.Dao().GetTableName(&entity.Comment{})

			var pages []entity.Page
			app.Dao().DB().Raw(
				"SELECT * FROM "+tbPages+" p WHERE p.site_name = ? ORDER BY ("+
					"SELECT COUNT(*) FROM "+tbComments+" c WHERE c.page_key = p.key AND c.is_pending = ?) DESC LIMIT ?",
				p.SiteName, false, p.Limit,
			).Find(&pages)

			return common.RespData(c, ResponseStat{
				Data: app.Dao().CookAllPages(pages),
			})

		case "page_pv":
			// ------------------------------------
			//  Query Multiple page PV
			// ------------------------------------
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
			// ------------------------------------
			//  Query Site total PV
			// ------------------------------------
			var pv int64
			app.Dao().DB().Model(&entity.Page{}).Where(&entity.Page{SiteName: p.SiteName}).Select("SUM(pv)").Scan(&pv)

			return common.RespData(c, ResponseStat{
				Data: pv,
			})

		case "page_comment":
			// ------------------------------------
			//  Query Multiple page comments
			// ------------------------------------
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
			// ------------------------------------
			//  Query Site total comments
			// ------------------------------------
			var count int64
			app.Dao().DB().Scopes(QueryComments).Count(&count)

			return common.RespData(c, ResponseStat{
				Data: count,
			})

		case "rand_comments":
			// ------------------------------------
			//  Random comments
			// ------------------------------------
			var comments []*entity.Comment
			app.Dao().DB().Scopes(QueryComments, QueryOrderRand).
				Limit(p.Limit).
				Find(&comments)

			return common.RespData(c, ResponseStat{
				Data: app.Dao().CookAllComments(comments),
			})

		case "rand_pages":
			// ------------------------------------
			//  Random pages
			// ------------------------------------
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
