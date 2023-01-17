package handler

import (
	"os"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// POST /api/admin/setting-get
func AdminSettingGet(router fiber.Router) {
	router.Post("/setting-get", func(c *fiber.Ctx) error {
		if !common.GetIsSuperAdmin(c) {
			return common.RespError(c, "无权访问")
		}

		dat, err := os.ReadFile(config.GetCfgFileLoaded())
		if err != nil {
			return common.RespError(c, "配置文件读取失败")
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
			return common.RespError(c, "无权访问")
		}

		var p ParamsAdminSettingSave
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		configFile := config.GetCfgFileLoaded()
		f, err := os.Create(configFile)
		if err != nil {
			return common.RespError(c, "配置文件读取失败，"+err.Error())
		}

		defer f.Close()

		_, err2 := f.WriteString(p.Data)
		if err2 != nil {
			return common.RespError(c, "保存失败，"+err2.Error())
		}

		// 重启服务
		workDir, err3 := os.Getwd()
		if err3 != nil {
			return common.RespError(c, "工作路径获取失败，"+err3.Error())
		}
		core.LoadCore(configFile, workDir)
		common.ReloadCorsAllowOrigins() // 刷新 CORS 可信域名
		logrus.Info("服务已重启完毕")

		return common.RespSuccess(c)
	})
}
