package handler

import (
	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsAdminSiteAdd struct {
	Name string `form:"name" validate:"required"`
	Urls string `form:"urls"`
}

// POST /api/admin/site-add
func AdminSiteAdd(router fiber.Router) {
	router.Post("/site-add", func(c *fiber.Ctx) error {
		var p ParamsAdminSiteAdd
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		if !common.GetIsSuperAdmin(c) {
			return common.RespError(c, i18n.T("Access denied"))
		}

		if p.Urls != "" {
			urls := utils.SplitAndTrimSpace(p.Urls, ",")
			for _, url := range urls {
				if !utils.ValidateURL(url) {
					return common.RespError(c, i18n.T("Contains invalid URL"))
				}
			}
		}

		if p.Name == config.ATK_SITE_ALL {
			return common.RespError(c, "Prohibit the use of reserved keywords as names")
		}

		if !query.FindSite(p.Name).IsEmpty() {
			return common.RespError(c, i18n.T("{{name}} already exists", Map{"name": i18n.T("Site")}))
		}

		site := entity.Site{}
		site.Name = p.Name
		site.Urls = p.Urls
		err := query.CreateSite(&site)
		if err != nil {
			return common.RespError(c, i18n.T("{{name}} creation failed", Map{"name": i18n.T("Site")}))
		}

		// 刷新 CORS 可信域名
		common.ReloadCorsAllowOrigins()

		return common.RespData(c, common.Map{
			"site": query.CookSite(&site),
		})
	})
}
