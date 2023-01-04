package common

import (
	"strings"

	"github.com/ArtalkJS/ArtalkGo/internal/config"
	"github.com/gofiber/fiber/v2"
)

func GetApiVersionDataMap() Map {
	return Map{
		"app":            "artalk-go",
		"version":        strings.TrimPrefix(config.Version, "v"),
		"commit_hash":    config.CommitHash,
		"fe_min_version": strings.TrimPrefix(config.FeMinVersion, "v"),
	}
}

func GetApiPublicConfDataMap(c *fiber.Ctx) Map {
	isAdmin := CheckIsAdminReq(c)
	imgUpload := config.Instance.ImgUpload.Enabled
	if isAdmin {
		imgUpload = true // 管理员始终允许上传图片
	}

	frontendConf := config.Instance.Frontend
	if frontendConf == nil {
		frontendConf = make(map[string]interface{})
	}
	frontendConf["imgUpload"] = &imgUpload

	return Map{
		"img_upload":    imgUpload,
		"frontend_conf": frontendConf,
	}
}
