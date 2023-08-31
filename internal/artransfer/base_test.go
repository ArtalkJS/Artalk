package artransfer

import (
	"encoding/json"
	"testing"

	"github.com/ArtalkJS/Artalk/internal/dao"
	"github.com/ArtalkJS/Artalk/internal/db"
	"github.com/ArtalkJS/Artalk/test"
	"github.com/stretchr/testify/assert"
)

func TestRunExportAndImportArtrans(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	var artransJSON string

	t.Run("Export", func(t *testing.T) {
		var err error
		artransJSON, err = RunExportArtrans(app.Dao(), &ExportParams{
			SiteNameScope: []string{},
		})
		if err != nil {
			t.Fatal(err)
		}
	})

	assert.NotEmpty(t, artransJSON)

	t.Run("Import", func(t *testing.T) {
		ddb, _ := db.NewTestDB()
		defer db.CloseDB(ddb)
		dao := dao.NewDao(ddb)

		RunImportArtrans(dao, &ImportParams{
			JsonData:  artransJSON,
			Assumeyes: true,
		})

		// check equals
		exported, err := RunExportArtrans(dao, &ExportParams{})
		if assert.NoError(t, err) {
			jsonArrCount := func(str string) int {
				var arr []any
				json.Unmarshal([]byte(str), &arr)
				return len(arr)
			}

			assert.Equal(t, jsonArrCount(artransJSON), jsonArrCount(exported))
		}
	})
}
