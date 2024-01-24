package handler

import (
	"strings"

	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsSiteCreate struct {
	Name string   `json:"name" validate:"required"` // The site name
	Urls []string `json:"urls" validate:"required"` // The site urls
}

type ResponseSiteCreate struct {
	entity.CookedSite
}

// @Id           CreateSite
// @Summary      Create Site
// @Description  Create a new site
// @Tags         Site
// @Security     ApiKeyAuth
// @Param        site  body  ParamsSiteCreate  true  "The site data"
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponseSiteCreate
// @Failure      400  {object}  Map{msg=string}
// @Failure      403  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /sites  [post]
func SiteCreate(app *core.App, router fiber.Router) {
	router.Post("/sites", common.AdminGuard(app, func(c *fiber.Ctx) error {
		var p ParamsSiteCreate
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		if len(p.Urls) > 0 {
			for _, url := range p.Urls {
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
		site.Urls = strings.Join(p.Urls, ",")
		err := app.Dao().CreateSite(&site)
		if err != nil {
			return common.RespError(c, 500, i18n.T("{{name}} creation failed", Map{"name": i18n.T("Site")}))
		}

		return common.RespData(c, ResponseSiteCreate{
			CookedSite: app.Dao().CookSite(&site),
		})
	}))
}
