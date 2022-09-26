package http

import (
	"os"

	"github.com/ArtalkJS/ArtalkGo/lib/core"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func (a *action) AdminSettingGet(c echo.Context) error {
	if !GetIsSuperAdmin(c) {
		return RespError(c, "无权访问")
	}

	dat, err := os.ReadFile(viper.ConfigFileUsed())
	if err != nil {
		return RespError(c, "配置文件读取失败")
	}

	return RespData(c, string(dat))
}

type ParamsAdminSettingSave struct {
	Data string `mapstructure:"data" param:"required"`
}

func (a *action) AdminSettingSave(c echo.Context) error {
	var p ParamsAdminSettingSave
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	if !GetIsSuperAdmin(c) {
		return RespError(c, "无权访问")
	}

	confFilename := viper.ConfigFileUsed()
	f, err := os.Create(confFilename)
	if err != nil {
		return RespError(c, "配置文件读取失败，"+err.Error())
	}

	defer f.Close()

	_, err2 := f.WriteString(p.Data)
	if err2 != nil {
		return RespError(c, "保存失败，"+err2.Error())
	}

	// 重启服务
	workDir, err3 := os.Getwd()
	if err3 != nil {
		return RespError(c, "工作路径获取失败，"+err3.Error())
	}
	core.LoadCore(confFilename, workDir)
	logrus.Info("服务已重启完毕")

	return RespSuccess(c)
}
