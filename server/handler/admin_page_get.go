package handler

import (
	"github.com/ArtalkJS/Artalk/internal/db"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsAdminPageGet struct {
	SiteName string
	SiteID   uint
	SiteAll  bool
	Limit    int `form:"limit"`
	Offset   int `form:"offset"`
}

type ResponseAdminPageGet struct {
	Total int64               `json:"total"`
	Pages []entity.CookedPage `json:"pages"`
}

// @Summary      Page List
// @Description  Get a list of pages by some conditions
// @Tags         Page
// @Param        site_name      formData  string  false  "the site name of your content scope"
// @Param        limit          formData  int     false  "the limit for pagination"
// @Param        offset         formData  int     false  "the offset for pagination"
// @Security     ApiKeyAuth
// @Success      200  {object}  common.JSONResult{data=ResponseAdminPageGet}
// @Router       /admin/page-get  [post]
func AdminPageGet(router fiber.Router) {
	router.Post("/page-get", func(c *fiber.Ctx) error {
		var p ParamsAdminPageGet
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// use site
		common.UseSite(c, &p.SiteName, &p.SiteID, &p.SiteAll)

		if !common.IsAdminHasSiteAccess(c, p.SiteName) {
			return common.RespError(c, i18n.T("Access denied"))
		}

		// 准备 query
		q := db.DB().Model(&entity.Page{}).Order("created_at DESC")
		if !p.SiteAll { // 不是查的所有站点
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
			cookedPages = append(cookedPages, query.CookPage(&p))
		}

		return common.RespData(c, ResponseAdminPageGet{
			Pages: cookedPages,
			Total: total,
		})
	})
}
