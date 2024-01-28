package artransfer

import (
	"github.com/ArtalkJS/Artalk/internal/dao"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
)

func RunExportArtrans(dao *dao.Dao, params *ExportParams) (string, error) {
	return exportArtrans(dao, params)
}

func RunImportArtrans(dao *dao.Dao, params *ImportParams) {
	if !params.UrlResolver {
		logWarn("Target site URL resolver disabled")
	}

	// 读取 JSON
	if params.JsonData == "" {
		if params.JsonFile == "" {
			logFatal(i18n.T("{{name}} is required", map[string]any{"name": "json_file:<JSON file path>"}))
			return
		}

		var err error
		params.JsonData, err = readJsonFile(params.JsonFile)
		if err != nil {
			logFatal(err)
			return
		}
	}

	// Json 转 Artran 实例列表
	comments := []*entity.Artran{}
	if err := jsonDecodeFAS(params.JsonData, &comments); err != nil {
		logFatal(err)
		return
	}

	// 执行导入数据
	importArtrans(dao, params, comments)

	print("\n")
	logInfo(i18n.T("Import complete"))
}

func ArrToImportParams(arr []string) *ImportParams {
	params := ImportParams{}

	params.UrlResolver = false // 默认关闭

	getParamsFrom(arr).To(map[string]any{
		"target_site_name": &params.TargetSiteName,
		"target_site_url":  &params.TargetSiteUrl,
		"url_resolver":     &params.UrlResolver,
		"json_file":        &params.JsonFile,
		"json_data":        &params.JsonData,
		"assumeyes":        &params.Assumeyes,
	})

	return &params
}
