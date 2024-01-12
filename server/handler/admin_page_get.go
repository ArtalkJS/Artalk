package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsAdminPageGet struct {
	SiteName string `query:"site_name" json:"site_name"` // The site name of your content scope
	Limit    int    `query:"limit" json:"limit"`         // The limit for pagination
	Offset   int    `query:"offset" json:"offset"`       // The offset for pagination
}

type ResponseAdminPageGet struct {
	Total int64               `json:"count"`
	Pages []entity.CookedPage `json:"pages"`
}

// @Summary      Get Page List
// @Description  Get a list of pages by some conditions
// @Tags         Page
// @Security     ApiKeyAuth
// @Param        options  query  ParamsAdminPageGet  true  "The options"
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponseAdminPageGet
// @Failure      403  {object}  Map{msg=string}
// @Router       /pages  [get]
func AdminPageGet(app *core.App, router fiber.Router) {
	router.Get("/pages", common.AdminGuard(app, func(c *fiber.Ctx) error {
		var p ParamsAdminPageGet
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		user := common.GetUserByReq(app, c)

		// 准备 query
		q := app.Dao().DB().Model(&entity.Page{}).Order("created_at DESC")
		if p.SiteName == "" {
			userSites := app.Dao().CookUser(&user).SiteNames
			if !user.IsEmpty() && len(userSites) > 0 {
				q = q.Where("site_name IN (?)", userSites)
			}
		} else {
			if _, ok, resp := common.CheckSiteExist(app, c, p.SiteName); !ok {
				return resp
			}

			if !common.IsAdminHasSiteAccess(app, c, p.SiteName) {
				return common.RespError(c, 403, i18n.T("Access denied"))
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

		return common.RespData(c, ResponseAdminPageGet{
			Pages: cookedPages,
			Total: total,
		})
	}))
}
