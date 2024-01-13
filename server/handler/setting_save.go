package handler

import (
	"os"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsAdminSettingSave struct {
	Yaml string `json:"yaml" validate:"required"` // The content of the config file in YAML format
}

// @Summary      Save Settings
// @Description  Save settings to app config file
// @Tags         System
// @Security     ApiKeyAuth
// @Param        settings  body  ParamsAdminSettingSave  true "The settings"
// @Accept       json
// @Produce      json
// @Success      200  {object}  Map{}
// @Failure      403  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /settings [post]
func AdminSettingSave(app *core.App, router fiber.Router) {
	router.Post("/settings", common.AdminGuard(app, func(c *fiber.Ctx) error {
		var p ParamsAdminSettingSave
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		configFile := app.Conf().GetCfgFileLoaded()
		f, err := os.Create(configFile)
		if err != nil {
			return common.RespError(c, 500, i18n.T("Config file read failed")+": "+err.Error())
		}

		defer f.Close()

		_, err2 := f.WriteString(p.Yaml)
		if err2 != nil {
			return common.RespError(c, 500, i18n.T("Save failed")+": "+err2.Error())
		}

		// 应用新配置文件
		conf, err := config.NewFromFile(configFile)
		if err != nil {
			return common.RespError(c, 500, "Config instance err: "+err.Error())
		}

		app.SetConf(conf)

		// 重启服务
		if err := app.Restart(); err != nil {
			return common.RespError(c, 500, i18n.T("Restart failed: {{err}}", map[string]interface{}{"err": err.Error()}))
		}

		log.Info(i18n.T("Services restart complete"))

		return common.RespSuccess(c)
	}))
}
