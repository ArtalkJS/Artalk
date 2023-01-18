package handler

import (
	"os"

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
			return common.RespError(c, i18n.T("No access"))
		}

		dat, err := os.ReadFile(config.GetCfgFileLoaded())
		if err != nil {
			return common.RespError(c, i18n.T("Config file failed to read"))
		}

		return common.RespData(c, string(dat))
	})
}

type ParamsAdminSettingSave struct {
	Data string `form:"data" validate:"required"`
}

// POST /api/admin/setting-save
func AdminSettingSave(router fiber.Router) {
	router.Post("/setting-save", func(c *fiber.Ctx) error {
		if !common.GetIsSuperAdmin(c) {
			return common.RespError(c, i18n.T("No access"))
		}

		var p ParamsAdminSettingSave
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		configFile := config.GetCfgFileLoaded()
		f, err := os.Create(configFile)
		if err != nil {
			return common.RespError(c, i18n.T("Config file failed to read")+", "+err.Error())
		}

		defer f.Close()

		_, err2 := f.WriteString(p.Data)
		if err2 != nil {
			return common.RespError(c, i18n.T("Failed to save")+", "+err2.Error())
		}

		// 重启服务
		workDir, err3 := os.Getwd()
		if err3 != nil {
			return common.RespError(c, i18n.T("Failed to get working directory")+", "+err3.Error())
		}
		core.LoadCore(configFile, workDir)
		common.ReloadCorsAllowOrigins() // 刷新 CORS 可信域名
		logrus.Info(i18n.T("Service restart is complete"))

		return common.RespSuccess(c)
	})
}
