package handler

import (
	"strings"

	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/dao"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsAdminPageEdit struct {
	SiteName string `json:"site_name"` // The site name of your content scope

	Key       string `json:"key"`        // Edit page key
	Title     string `json:"title"`      // Edit page title
	AdminOnly bool   `json:"admin_only"` // Edit page admin_only option
}

type ResponseAdminPageEdit struct {
	entity.CookedPage
}

// @Summary      Edit Page
// @Description  Edit a specific page
// @Tags         Page
// @Security     ApiKeyAuth
// @Param        id    path  int                  true "The page ID you want to edit"
// @Param        page  body  ParamsAdminPageEdit  true "The page data"
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponseAdminPageEdit
// @Failure      400  {object}  Map{msg=string}
// @Failure      403  {object}  Map{msg=string}
// @Failure      404  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /pages/{id}  [put]
func AdminPageEdit(app *core.App, router fiber.Router) {
	router.Put("/pages/:id", common.AdminGuard(app, func(c *fiber.Ctx) error {
		id, _ := c.ParamsInt("id")

		var p ParamsAdminPageEdit
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		if strings.TrimSpace(p.Key) == "" {
			return common.RespError(c, 400, i18n.T("{{name}} cannot be empty", Map{"name": "key"}))
		}

		// check site exist
		if _, ok, resp := common.CheckSiteExist(app, c, p.SiteName); !ok {
			return resp
		}

		// find page
		var page = app.Dao().FindPageByID(uint(id))
		if page.IsEmpty() {
			return common.RespError(c, 404, i18n.T("{{name}} not found", Map{"name": i18n.T("Page")}))
		}

		// 重命名合法性检测
		modifyKey := p.Key != page.Key
		if modifyKey && !app.Dao().FindPage(p.Key, page.SiteName).IsEmpty() {
			return common.RespError(c, 400, i18n.T("{{name}} already exists", Map{"name": i18n.T("Page")}))
		}

		// 预先删除缓存，防止修改主键原有 page_key 占用问题
		app.Dao().CacheAction(func(cache *dao.DaoCache) {
			cache.PageCacheDel(&page)
		})

		page.Title = p.Title
		page.AdminOnly = p.AdminOnly
		if modifyKey {
			// 相关性数据修改
			var comments []entity.Comment
			app.Dao().DB().Where("page_key = ?", page.Key).Find(&comments)

			for _, comment := range comments {
				comment.PageKey = p.Key
				app.Dao().UpdateComment(&comment)
			}

			page.Key = p.Key
		}

		if err := app.Dao().UpdatePage(&page); err != nil {
			return common.RespError(c, 500, i18n.T("{{name}} save failed", Map{"name": i18n.T("Page")}))
		}

		return common.RespData(c, ResponseAdminPageEdit{
			CookedPage: app.Dao().CookPage(&page),
		})
	}))
}
