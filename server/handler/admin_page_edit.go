package handler

import (
	"strings"

	"github.com/ArtalkJS/ArtalkGo/internal/cache"
	"github.com/ArtalkJS/ArtalkGo/internal/db"
	"github.com/ArtalkJS/ArtalkGo/internal/entity"
	"github.com/ArtalkJS/ArtalkGo/internal/query"
	"github.com/ArtalkJS/ArtalkGo/server/common"
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

// POST /api/admin/page-edit
func AdminPageEdit(router fiber.Router) {
	router.Post("/page-edit", func(c *fiber.Ctx) error {
		var p ParamsAdminPageEdit
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		if strings.TrimSpace(p.Key) == "" {
			return common.RespError(c, "page key 不能为空白字符")
		}

		// use site
		common.UseSite(c, &p.SiteName, &p.SiteID, nil)

		// find page
		var page = query.FindPageByID(p.ID)
		if page.IsEmpty() {
			return common.RespError(c, "page not found")
		}

		if !common.IsAdminHasSiteAccess(c, page.SiteName) {
			return common.RespError(c, "无权操作")
		}

		// 重命名合法性检测
		modifyKey := p.Key != page.Key
		if modifyKey && !query.FindPage(p.Key, p.SiteName).IsEmpty() {
			return common.RespError(c, "page 已存在，请换个 key")
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
			return common.RespError(c, "page save error")
		}

		return common.RespData(c, common.Map{
			"page": query.CookPage(&page),
		})
	})
}
