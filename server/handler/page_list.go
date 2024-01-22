package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsPageList struct {
	SiteName string `query:"site_name" json:"site_name" validate:"optional"` // The site name of your content scope
	Limit    int    `query:"limit" json:"limit" validate:"optional"`         // The limit for pagination
	Offset   int    `query:"offset" json:"offset" validate:"optional"`       // The offset for pagination
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

		// 准备 query
		q := app.Dao().DB().Model(&entity.Page{}).Order("created_at DESC")
		if p.SiteName != "" {
			if _, ok, resp := common.CheckSiteExist(app, c, p.SiteName); !ok {
				return resp
			}

			q = q.Where("site_name = ?", p.SiteName)
		}

		// 总共条数
		var total int64
		q.Count(&total)

		// 数据分页
		q = q.Scopes(Paginate(p.Offset, p.Limit))

		// 查找
		var pages []entity.Page
		q.Find(&pages)

		var cookedPages []entity.CookedPage
		for _, p := range pages {
			cookedPages = append(cookedPages, app.Dao().CookPage(&p))
		}

		return common.RespData(c, ResponsePageList{
			Pages: cookedPages,
			Total: total,
		})
	}))
}
