package handler

import (
	"os"

	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ResponseAdminSettingGet struct {
	Yaml string `json:"yaml"`
}

// @Summary      Get Settings
// @Description  Get settings from app config file
// @Tags         System
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  ResponseAdminSettingGet
// @Failure      403  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /settings [get]
func AdminSettingGet(app *core.App, router fiber.Router) {
	router.Get("/settings", common.AdminGuard(app, func(c *fiber.Ctx) error {
		dat, err := os.ReadFile(app.Conf().GetCfgFileLoaded())
		if err != nil {
			return common.RespError(c, 500, i18n.T("Config file read failed"))
		}

		return common.RespData(c, ResponseAdminSettingGet{
			Yaml: string(dat),
		})
	}))
}
