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
	// 查询值
	ID uint `form:"id" validate:"required"`

	// 修改值
	Name string `form:"name"`
	Urls string `form:"urls"`
}

type ResponseAdminSiteEdit struct {
	Site entity.CookedSite `json:"site"`
}

// @Summary      Site Edit
// @Description  Edit a specific site
// @Tags         Site
// @Param        id             formData  string  true   "the site ID you want to edit"
// @Param        name           formData  string  false  "edit site name"
// @Param        urls           formData  string  false  "edit site urls"
// @Security     ApiKeyAuth
// @Success      200  {object}  common.JSONResult{data=ResponseAdminSiteEdit}
// @Router       /admin/site-edit  [post]
func AdminSiteEdit(app *core.App, router fiber.Router) {
	router.Post("/site-edit", func(c *fiber.Ctx) error {
		var p ParamsAdminSiteEdit
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		site := app.Dao().FindSiteByID(p.ID)
		if site.IsEmpty() {
			return common.RespError(c, i18n.T("{{name}} not found", Map{"name": i18n.T("Site")}))
		}

		// 站点操作权限检查
		if !common.IsAdminHasSiteAccess(app, c, site.Name) {
			return common.RespError(c, i18n.T("Access denied"))
		}

		if strings.TrimSpace(p.Name) == "" {
			return common.RespError(c, i18n.T("{{name}} cannot be empty", Map{"name": "name"}))
		}

		// 重命名合法性检测
		modifyName := p.Name != site.Name
		if modifyName && !app.Dao().FindSite(p.Name).IsEmpty() {
			return common.RespError(c, i18n.T("{{name}} already exists", Map{"name": i18n.T("Site")}))
		}

		// urls 合法性检测
		if p.Urls != "" {
			for _, url := range utils.SplitAndTrimSpace(p.Urls, ",") {
				if !utils.ValidateURL(url) {
					return common.RespError(c, i18n.T("Contains invalid URL"))
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
			return common.RespError(c, i18n.T("{{name}} save failed", Map{"name": i18n.T("Site")}))
		}

		return common.RespData(c, ResponseAdminSiteEdit{
			Site: app.Dao().CookSite(&site),
		})
	})
}
