package handler

import (
	"strings"

	"github.com/ArtalkJS/Artalk/internal/cache"
	"github.com/ArtalkJS/Artalk/internal/db"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsAdminPageEdit struct {
	// 查询值
	ID       uint `form:"id"`
	SiteName string
	SiteID   uint

	// 修改值
	Key       string `form:"key"`
	Title     string `form:"title"`
	AdminOnly bool   `form:"admin_only"`
}

type ResponseAdminPageEdit struct {
	Page entity.CookedPage `json:"page"`
}

// @Summary      Page Edit
// @Description  Edit a specific page
// @Tags         Page
// @Param        id             formData  string  true   "the page ID you want to edit"
// @Param        site_name      formData  string  false  "the site name of your content scope"
// @Param        key            formData  string  false  "edit page key"
// @Param        title          formData  string  false  "edit page title"
// @Param        admin_only     formData  bool    false  "edit page admin_only option"
// @Security     ApiKeyAuth
// @Success      200  {object}  common.JSONResult{data=ResponseAdminPageEdit}
// @Router       /admin/page-edit  [post]
func AdminPageEdit(router fiber.Router) {
	router.Post("/page-edit", func(c *fiber.Ctx) error {
		var p ParamsAdminPageEdit
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		if strings.TrimSpace(p.Key) == "" {
			return common.RespError(c, i18n.T("{{name}} cannot be empty", Map{"name": "key"}))
		}

		// use site
		common.UseSite(c, &p.SiteName, &p.SiteID, nil)

		// find page
		var page = query.FindPageByID(p.ID)
		if page.IsEmpty() {
			return common.RespError(c, i18n.T("{{name}} not found", Map{"name": i18n.T("Page")}))
		}

		if !common.IsAdminHasSiteAccess(c, page.SiteName) {
			return common.RespError(c, i18n.T("Access denied"))
		}

		// 重命名合法性检测
		modifyKey := p.Key != page.Key
		if modifyKey && !query.FindPage(p.Key, p.SiteName).IsEmpty() {
			return common.RespError(c, i18n.T("{{name}} already exists", Map{"name": i18n.T("Page")}))
		}

		// 预先删除缓存，防止修改主键原有 page_key 占用问题
		cache.PageCacheDel(&page)

		page.Title = p.Title
		page.AdminOnly = p.AdminOnly
		if modifyKey {
			// 相关性数据修改
			var comments []entity.Comment
			db.DB().Where("page_key = ?", page.Key).Find(&comments)

			for _, comment := range comments {
				comment.PageKey = p.Key
				query.UpdateComment(&comment)
			}

			page.Key = p.Key
		}

		if err := query.UpdatePage(&page); err != nil {
			return common.RespError(c, i18n.T("{{name}} save failed", Map{"name": i18n.T("Page")}))
		}

		return common.RespData(c, ResponseAdminPageEdit{
			Page: query.CookPage(&page),
		})
	})
}
