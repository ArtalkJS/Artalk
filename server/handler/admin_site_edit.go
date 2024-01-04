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

type ParamsAdminSiteEdit struct {
	Name string `json:"name"` // Edit site name
	Urls string `json:"urls"` // Edit site urls
}

type ResponseAdminSiteEdit struct {
	entity.CookedSite
}

// @Summary      Edit Site
// @Description  Edit a specific site
// @Tags         Site
// @Security     ApiKeyAuth
// @Param        id    path  string               true  "The site ID you want to edit"
// @Param        site  body  ParamsAdminSiteEdit  true  "The site data"
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponseAdminSiteEdit
// @Router       /sites/{id}  [put]
func AdminSiteEdit(app *core.App, router fiber.Router) {
	router.Put("/sites/:id", func(c *fiber.Ctx) error {
		id, _ := c.ParamsInt("id")

		var p ParamsAdminSiteEdit
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		site := app.Dao().FindSiteByID(uint(id))
		if site.IsEmpty() {
			return common.RespError(c, 404, i18n.T("{{name}} not found", Map{"name": i18n.T("Site")}))
		}

		// 站点操作权限检查
		if !common.IsAdminHasSiteAccess(app, c, site.Name) {
			return common.RespError(c, 403, i18n.T("Access denied"))
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
		if p.Urls != "" {
			for _, url := range utils.SplitAndTrimSpace(p.Urls, ",") {
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
		site.Urls = p.Urls

		err := app.Dao().UpdateSite(&site)
		if err != nil {
			return common.RespError(c, 500, i18n.T("{{name}} save failed", Map{"name": i18n.T("Site")}))
		}

		return common.RespData(c, ResponseAdminSiteEdit{
			CookedSite: app.Dao().CookSite(&site),
		})
	})
}
