package common

import (
	"strings"

	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/gofiber/fiber/v2"
)

func CheckSiteExist(app *core.App, c *fiber.Ctx, siteName string) (entity.CookedSite, bool, error) {
	siteName = strings.TrimSpace(siteName)

	if siteName == "" {
		return entity.CookedSite{}, false, RespError(c,
			400,
			i18n.T("{{name}} cannot be empty", Map{"name": i18n.T("Site name")}),
			Map{
				"err_no_site": true,
			},
		)
	}

	findSite := app.Dao().FindSite(siteName)
	if findSite.IsEmpty() {
		return entity.CookedSite{}, false, RespError(c,
			404,
			i18n.T("Site `{{name}}` not found. Please create it in control center.", map[string]interface{}{"name": siteName}),
			Map{
				"err_no_site": true,
			},
		)
	}

	return app.Dao().CookSite(&findSite), true, nil
}
