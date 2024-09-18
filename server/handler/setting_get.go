package handler

import (
	"os"
	"strings"

	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/i18n"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

type ResponseSettingGet struct {
	Yaml string   `json:"yaml"`
	Envs []string `json:"envs"`
}

// @Id           GetSettings
// @Summary      Get Settings
// @Description  Get settings from app config file
// @Tags         System
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  ResponseSettingGet
// @Failure      403  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /settings [get]
func SettingGet(app *core.App, router fiber.Router) {
	router.Get("/settings", common.AdminGuard(app, func(c *fiber.Ctx) error {
		dat, err := os.ReadFile(app.Conf().GetCfgFileLoaded())
		if err != nil {
			return common.RespError(c, 500, i18n.T("Config file read failed"))
		}

		// get all environment variables which start with ATK_
		envs := lo.Filter(os.Environ(), func(v string, _ int) bool {
			return strings.HasPrefix(v, "ATK_")
		})

		return common.RespData(c, ResponseSettingGet{
			Yaml: string(dat),
			Envs: envs,
		})
	}))
}
