package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsAdminSiteAdd struct {
	Name string `json:"name" validate:"required"` // The site name
	Urls string `json:"urls"`                     // The site urls
}

type ResponseAdminSiteAdd struct {
	entity.CookedSite
}

// @Summary      Create Site
// @Description  Create a new site
// @Tags         Site
// @Security     ApiKeyAuth
// @Param        site  body  ParamsAdminSiteAdd  true  "The site data"
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponseAdminSiteAdd
// @Failure      400  {object}  Map{msg=string}
// @Failure      403  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /sites  [post]
func AdminSiteAdd(app *core.App, router fiber.Router) {
	router.Post("/sites", common.AdminGuard(app, func(c *fiber.Ctx) error {
		var p ParamsAdminSiteAdd
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		if !common.GetIsSuperAdmin(app, c) {
			return common.RespError(c, 403, i18n.T("Access denied"))
		}

		if p.Urls != "" {
			urls := utils.SplitAndTrimSpace(p.Urls, ",")
			for _, url := range urls {
				if !utils.ValidateURL(url) {
					return common.RespError(c, 400, i18n.T("Contains invalid URL"))
				}
			}
		}

		if !app.Dao().FindSite(p.Name).IsEmpty() {
			return common.RespError(c, 400, i18n.T("{{name}} already exists", Map{"name": i18n.T("Site")}))
		}

		site := entity.Site{}
		site.Name = p.Name
		site.Urls = p.Urls
		err := app.Dao().CreateSite(&site)
		if err != nil {
			return common.RespError(c, 500, i18n.T("{{name}} creation failed", Map{"name": i18n.T("Site")}))
		}

		return common.RespData(c, ResponseAdminSiteAdd{
			CookedSite: app.Dao().CookSite(&site),
		})
	}))
}
