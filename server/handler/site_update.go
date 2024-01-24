package handler

import (
	"strings"

	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/dao"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsSiteUpdate struct {
	Name string   `json:"name" validate:"required"` // Updated site name
	Urls []string `json:"urls" validate:"required"` // Updated site urls
}

type ResponseSiteUpdate struct {
	entity.CookedSite
}

// @Id           UpdateSite
// @Summary      Update Site
// @Description  Update a specific site
// @Tags         Site
// @Security     ApiKeyAuth
// @Param        id    path  int               true  "The site ID you want to update"
// @Param        site  body  ParamsSiteUpdate  true  "The site data"
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponseSiteUpdate
// @Router       /sites/{id}  [put]
func SiteUpdate(app *core.App, router fiber.Router) {
	router.Put("/sites/:id", common.AdminGuard(app, func(c *fiber.Ctx) error {
		id, _ := c.ParamsInt("id")

		var p ParamsSiteUpdate
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		site := app.Dao().FindSiteByID(uint(id))
		if site.IsEmpty() {
			return common.RespError(c, 404, i18n.T("{{name}} not found", Map{"name": i18n.T("Site")}))
		}

		if strings.TrimSpace(p.Name) == "" {
			return common.RespError(c, 400, i18n.T("{{name}} cannot be empty", Map{"name": "name"}))
		}

		// 重命名合法性检测
		modifyName := p.Name != site.Name
		if modifyName && !app.Dao().FindSite(p.Name).IsEmpty() {
			return common.RespError(c, 400, i18n.T("{{name}} already exists", Map{"name": i18n.T("Site")}))
		}

		// urls 合法性检测
		if len(p.Urls) > 0 {
			for _, url := range p.Urls {
				if !utils.ValidateURL(url) {
					return common.RespError(c, 400, i18n.T("Contains invalid URL"))
				}
			}
		}

		// 预先删除缓存，防止修改主键原有 site_name 占用问题
		app.Dao().CacheAction(func(cache *dao.DaoCache) {
			cache.SiteCacheDel(&site)
		})

		// 同步变更 site_name
		if modifyName {
			var comments []entity.Comment
			var pages []entity.Page

			app.Dao().DB().Where("site_name = ?", site.Name).Find(&comments)
			app.Dao().DB().Where("site_name = ?", site.Name).Find(&pages)

			for _, comment := range comments {
				comment.SiteName = p.Name
				app.Dao().UpdateComment(&comment)
			}
			for _, page := range pages {
				page.SiteName = p.Name
				app.Dao().UpdatePage(&page)
			}
		}

		// 修改 site
		site.Name = p.Name
		site.Urls = strings.Join(p.Urls, ",")

		err := app.Dao().UpdateSite(&site)
		if err != nil {
			return common.RespError(c, 500, i18n.T("{{name}} save failed", Map{"name": i18n.T("Site")}))
		}

		return common.RespData(c, ResponseSiteUpdate{
			CookedSite: app.Dao().CookSite(&site),
		})
	}))
}
