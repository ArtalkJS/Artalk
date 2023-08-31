package handler

import (
	"net/http"
	"os"
	"strings"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ResponseAdminSettingGet struct {
	Custom   string `json:"custom"`
	Template string `json:"template"`
}

// @Summary      Settings Get
// @Description  Get settings from app config file
// @Tags         System
// @Security     ApiKeyAuth
// @Success      200  {object}  common.JSONResult{data=ResponseAdminSettingGet}
// @Router       /admin/setting-get [post]
func AdminSettingGet(app *core.App, router fiber.Router) {
	router.Post("/setting-get", func(c *fiber.Ctx) error {
		if !common.GetIsSuperAdmin(app, c) {
			return common.RespError(c, i18n.T("Access denied"))
		}

		dat, err := os.ReadFile(app.Conf().GetCfgFileLoaded())
		if err != nil {
			return common.RespError(c, i18n.T("Config file read failed"))
		}

		return common.RespData(c, ResponseAdminSettingGet{
			Custom:   string(dat),
			Template: app.ConfTpl(),
		})
	})
}

type ParamsAdminSettingSave struct {
	Data string `form:"data" validate:"required"`
}

// @Summary      Settings Save
// @Description  Save settings to app config file
// @Tags         System
// @Param        data           formData  string  true  "the content of the config file in YAML format"
// @Security     ApiKeyAuth
// @Success      200  {object}  common.JSONResult
// @Router       /admin/setting-save [post]
func AdminSettingSave(app *core.App, router fiber.Router) {
	router.Post("/setting-save", func(c *fiber.Ctx) error {
		if !common.GetIsSuperAdmin(app, c) {
			return common.RespError(c, i18n.T("Access denied"))
		}

		var p ParamsAdminSettingSave
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		configFile := app.Conf().GetCfgFileLoaded()
		f, err := os.Create(configFile)
		if err != nil {
			return common.RespError(c, i18n.T("Config file read failed")+": "+err.Error())
		}

		defer f.Close()

		_, err2 := f.WriteString(p.Data)
		if err2 != nil {
			return common.RespError(c, i18n.T("Save failed")+": "+err2.Error())
		}

		// 重启服务
		if err := app.Restart(); err != nil {
			return common.RespError(c, i18n.T("Restart failed: {{err}}", map[string]interface{}{"err": err.Error()}))
		}

		log.Info(i18n.T("Services restart complete"))

		return common.RespSuccess(c)
	})
}

// @Summary      Settings Template
// @Description  Get config templates in different languages for rendering the settings page in the frontend
// @Tags         System
// @Security     ApiKeyAuth
// @Success      200  {object}  string
// @Router       /admin/setting-tpl  [post]
func AdminSettingTpl(app *core.App, router fiber.Router) {
	router.Post("/setting-tpl", func(c *fiber.Ctx) error {
		locale := strings.TrimSpace(c.FormValue("locale"))
		tpl := config.Template(locale)
		return c.Status(http.StatusOK).SendString(tpl)
	})
}
