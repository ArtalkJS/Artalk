package handler

import (
	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ResponseAdminSettingTpl struct {
	Yaml string `json:"yaml"`
}

// @Summary      Get Settings Template
// @Description  Get config templates in different languages for rendering the settings page in the frontend
// @Tags         System
// @Security     ApiKeyAuth
// @Param        locale  path  string  false  "The locale of the settings template you want to get"
// @Produce      json
// @Success      200  {object}  ResponseAdminSettingTpl
// @Router       /settings/template/{locale}  [get]
func AdminSettingTpl(app *core.App, router fiber.Router) {
	router.Get("/settings/template/:locale?", common.AdminGuard(app, func(c *fiber.Ctx) error {
		var tpl string

		locale := c.Params("locale")
		if locale == "" {
			tpl = app.ConfTpl()
		} else {
			tpl = config.Template(locale)
		}

		return common.RespData(c, ResponseAdminSettingTpl{
			Yaml: tpl,
		})
	}))
}
