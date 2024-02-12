package common

import (
	"fmt"
	"strings"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/ArtalkJS/Artalk/server/middleware"
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
		CommitHash: config.CommitHash,
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

	if pluginURLs, ok := frontendConf["pluginURLs"].([]any); ok {
		frontendConf["pluginURLs"] = handlePluginURLs(app,
			lo.Map[any, string](pluginURLs, func(u any, _ int) string {
				return strings.TrimSpace(fmt.Sprintf("%v", u))
			}))
	}

	return ConfData{
		FrontendConf: frontendConf,
		Version:      GetApiVersionDataMap(),
	}
}

func handlePluginURLs(app *core.App, urls []string) []string {
	return lo.Filter[string](urls, func(u string, _ int) bool {
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
	})
}
