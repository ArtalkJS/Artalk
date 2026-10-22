package handler

import (
	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ParamsPageList struct {
	HasComments bool   `query:"has_comments" json:"has_comments" validate:"optional"` // Only include pages with comments (including pending comments)
	SiteName    string `query:"site_name" json:"site_name" validate:"optional"`       // The site name of your content scope
	Limit       int    `query:"limit" json:"limit" validate:"optional"`               // The limit for pagination
	Offset      int    `query:"offset" json:"offset" validate:"optional"`             // The offset for pagination
	Search      string `query:"search" json:"search" validate:"optional"`             // Search keywords
}

type ResponsePageList struct {
	Total int64               `json:"count"`
	Pages []entity.CookedPage `json:"pages"`
}

// @Id           GetPages
// @Summary      Get Page List
// @Description  Get a list of pages by some conditions
// @Tags         Page
// @Security     ApiKeyAuth
// @Param        options  query  ParamsPageList  true  "The options"
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponsePageList
// @Failure      403  {object}  Map{msg=string}
// @Router       /pages  [get]
func PageList(app *core.App, router fiber.Router) {
	router.Get("/pages", common.AdminGuard(app, func(c *fiber.Ctx) error {
		var p ParamsPageList
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// Prepare query
		q := app.Dao().DB().Model(&entity.Page{}).Order("created_at DESC")
		if p.SiteName != "" {
			if _, ok, resp := common.CheckSiteExist(app, c, p.SiteName); !ok {
				return resp
			}

			q = q.Where("site_name = ?", p.SiteName)
		}

		// Filter before count and pagination. A page key is only unique within a site.
		if p.HasComments {
			pagesTable := app.Dao().GetTableName(&entity.Page{})
			commentsTable := app.Dao().GetTableName(&entity.Comment{})
			comments := app.Dao().DB().Model(&entity.Comment{}).Select("1").
				Where(commentsTable + ".page_key = " + pagesTable + ".key").
				Where(commentsTable + ".site_name = " + pagesTable + ".site_name")
			q = q.Where("EXISTS (?)", comments)
		}

		// Search
		if p.Search != "" {
			q = q.Scopes(func(d *gorm.DB) *gorm.DB {
				// Because historical reasons, the naming of this field named `key` does not follow best practices.
				// In some database, directly use the field name `key` will cause an error.
				// So must keep the table name before the field name.
				tbPages := app.Dao().GetTableName(&entity.Page{})
				return d.Where("LOWER("+tbPages+".key) LIKE LOWER(?) OR LOWER(title) LIKE LOWER(?)",
					"%"+p.Search+"%", "%"+p.Search+"%")
			})
		}

		// Total count
		var total int64
		q.Count(&total)

		// Pagination
		q = q.Scopes(Paginate(p.Offset, p.Limit))

		// Find database
		var pages []entity.Page
		q.Find(&pages)

		cookedPages := make([]entity.CookedPage, 0, len(pages))
		for _, p := range pages {
			cookedPages = append(cookedPages, app.Dao().CookPage(&p))
		}

		return common.RespData(c, ResponsePageList{
			Pages: cookedPages,
			Total: total,
		})
	}))
}
