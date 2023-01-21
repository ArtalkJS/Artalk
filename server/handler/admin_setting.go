package handler

import (
	"net/http"
	"os"
	"strings"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// POST /api/admin/setting-get
func AdminSettingGet(router fiber.Router) {
	router.Post("/setting-get", func(c *fiber.Ctx) error {
		if !common.GetIsSuperAdmin(c) {
			return common.RespError(c, i18n.T("Access denied"))
		}

		dat, err := os.ReadFile(config.GetCfgFileLoaded())
		if err != nil {
			return common.RespError(c, i18n.T("Config file read failed"))
		}

		return common.RespData(c, Map{
			"custom":   string(dat),
			"template": core.GetConfTpl(),
		})
	})
}

type ParamsAdminSettingSave struct {
	Data string `form:"data" validate:"required"`
}

// POST /api/admin/setting-save
func AdminSettingSave(router fiber.Router) {
	router.Post("/setting-save", func(c *fiber.Ctx) error {
		if !common.GetIsSuperAdmin(c) {
			return common.RespError(c, i18n.T("Access denied"))
		}

		var p ParamsAdminSettingSave
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		configFile := config.GetCfgFileLoaded()
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
		workDir, err3 := os.Getwd()
		if err3 != nil {
			return common.RespError(c, i18n.T("Working directory retrieval failed")+": "+err3.Error())
		}
		core.LoadCore(configFile, workDir)
		common.ReloadCorsAllowOrigins() // 刷新 CORS 可信域名
		logrus.Info(i18n.T("Services restart complete"))

		return common.RespSuccess(c)
	})
}

// POST /api/admin/setting-tpl
func AdminSettingTpl(router fiber.Router) {
	router.Post("/setting-tpl", func(c *fiber.Ctx) error {
		locale := strings.TrimSpace(c.FormValue("locale"))
		tpl := core.GetConfTpl(locale)
		return c.Status(http.StatusOK).SendString(tpl)
	})
}
