package artransfer

import (
	"fmt"

	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/i18n"
	"gorm.io/gorm"
)

func RunExportArtrans(dao *dao.Dao, params *ExportParams) (string, error) {
	return exportArtrans(dao.DB(), params)
}

func RunImportArtrans(dao *dao.Dao, params *ImportParams, outputFunc ...func(string)) error {
	console := NewConsole()
	if len(outputFunc) > 0 {
		console.SetOutputFunc(outputFunc[0])
	}
	params.SetConsole(console)

	// Read JSON
	if params.JsonData == "" {
		if params.JsonFile == "" {
			err := fmt.Errorf(i18n.T("{{name}} is required", map[string]any{"name": "json_file:<JSON file path>"}))
			console.Error(err)
			return err
		}

		var err error
		params.JsonData, err = readJsonFile(params.JsonFile)
		if err != nil {
			console.Error(err)
			return err
		}
	}

	// Json to Artrans
	comments := []*entity.Artran{}
	if err := jsonDecodeFAS(params.JsonData, &comments); err != nil {
		console.Error(err)
		return err
	}

	// Execute import
	err := dao.DB().Transaction(func(tx *gorm.DB) error {
		return importArtrans(tx, params, comments)
	})

	if err != nil {
		console.Error("[Artransfer] ", i18n.T("Import failed"), ": ", err)
	} else {
		console.Info("[Artransfer] ", i18n.T("Import completed"))
	}

	return err
}
