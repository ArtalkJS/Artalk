package handler

import (
	"github.com/ArtalkJS/ArtalkGo/internal/query"
	"github.com/ArtalkJS/ArtalkGo/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsAdminSiteDel struct {
	ID uint `form:"id" validate:"required"`
}

// POST /api/admin/site-del
func AdminSiteDel(router fiber.Router) {
	router.Post("/site-del", func(c *fiber.Ctx) error {
		var p ParamsAdminSiteDel
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		if !common.GetIsSuperAdmin(c) {
			return common.RespError(c, "禁止删除站点")
		}

		site := query.FindSiteByID(p.ID)
		if site.IsEmpty() {
			return common.RespError(c, "site 不存在")
		}

		err := query.DelSite(&site)
		if err != nil {
			return common.RespError(c, "site 删除失败")
		}

		// 刷新 CORS 可信域名
		common.ReloadCorsAllowOrigins()

		return common.RespSuccess(c)
	})
}
