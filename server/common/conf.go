package common

import (
	"fmt"
	"slices"
	"strings"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/utils"
	"github.com/artalkjs/artalk/v2/server/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

type ApiVersionData struct {
	App        string `json:"app"`
	Version    string `json:"version"`
	CommitHash string `json:"commit_hash"`
}

func GetApiVersionDataMap() ApiVersionData {
	return ApiVersionData{
		App:        "artalk",
		Version:    strings.TrimPrefix(config.Version, "v"),
		CommitHash: config.CommitHash(),
	}
}

type ConfData struct {
	FrontendConf Map            `json:"frontend_conf"`
	Version      ApiVersionData `json:"version"`
}

func GetApiPublicConfDataMap(app *core.App, c *fiber.Ctx) ConfData {
	isAdmin := CheckIsAdminReq(app, c)
	imgUpload := app.Conf().ImgUpload.Enabled
	if isAdmin {
		imgUpload = true // 管理员始终允许上传图片
	}

	frontendConfSrc := app.Conf().Frontend
	if frontendConfSrc == nil {
		frontendConfSrc = make(map[string]interface{})
	}

	frontendConf := make(map[string]interface{})
	utils.CopyStruct(&frontendConfSrc, &frontendConf)

	frontendConf["imgUpload"] = &imgUpload
	if app.Conf().Locale != "" {
		frontendConf["locale"] = app.Conf().Locale
	}

	if _, ok := frontendConf["pluginURLs"].([]any); !ok {
		frontendConf["pluginURLs"] = []any{}
	}
	pluginURLs := frontendConf["pluginURLs"].([]any)

	if app.Conf().Auth.Enabled {
		pluginURLs = append(pluginURLs, "dist/plugins/artalk-plugin-auth.js")
	}

	if !slices.Contains([]string{"en", "zh-CN", ""}, app.Conf().Locale) {
		pluginURLs = append(pluginURLs, fmt.Sprintf("dist/i18n/%s.js", app.Conf().Locale))
	}

	frontendConf["pluginURLs"] = handlePluginURLs(app,
		lo.Map(pluginURLs, func(u any, _ int) string {
			return strings.TrimSpace(fmt.Sprintf("%v", u))
		}))

	return ConfData{
		FrontendConf: frontendConf,
		Version:      GetApiVersionDataMap(),
	}
}

func handlePluginURLs(app *core.App, urls []string) []string {
	return utils.RemoveDuplicates(lo.Filter(urls, func(u string, _ int) bool {
		if strings.TrimSpace(u) == "" {
			return false
		}
		if !utils.ValidateURL(u) {
			return true
		}
		if trusted, _, _ := middleware.CheckURLTrusted(app, u); trusted {
			return true
		}
		return false
	}))
}
