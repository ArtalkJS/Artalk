package artransfer

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/internal/db"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/test"
	"github.com/stretchr/testify/assert"
)

func TestRunExportAndImportArtrans(t *testing.T) {
	getRootIDSnapshot := func(dao *dao.Dao) string {
		var allComments []entity.Comment
		dao.DB().Find(&allComments)

		type item struct {
			ID     uint `json:"id"`
			RID    uint `json:"rid"`
			RootID uint `json:"root_id"`
		}
		data := []item{}
		for _, c := range allComments {
			data = append(data, item{c.ID, c.Rid, c.RootID})
		}

		jsonByte, _ := json.Marshal(data)
		return string(jsonByte)
	}

	var artransJSON string
	var idSnapshotBeforeImport string

	t.Run("Export", func(t *testing.T) {
		app, _ := test.NewTestApp()
		defer app.Cleanup()

		idSnapshotBeforeImport = getRootIDSnapshot(app.Dao())

		var err error
		if artransJSON, err = RunExportArtrans(app.Dao(), &ExportParams{
			SiteNameScope: []string{},
		}); err != nil {
			t.Fatal(err, "Export failed")
		}
	})

	assert.NotEmpty(t, artransJSON)

	t.Run("Import", func(t *testing.T) {
		// init a new clean db without any records
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
			assert.JSONEq(t, rebuildAutoIncrementIDs(artransJSON), rebuildAutoIncrementIDs(exported), "Exported data should be the same after import")
		}

		// check root_id
		idSnapshotAfterImport := getRootIDSnapshot(dao)

		handledIDSnapshotBefore := rebuildAutoIncrementIDs(idSnapshotBeforeImport)
		handledIDSnapshotAfter := rebuildAutoIncrementIDs(idSnapshotAfterImport)

		fmt.Println("Handled Before:", handledIDSnapshotBefore)
		fmt.Println("Handled After:", handledIDSnapshotAfter)

		assert.JSONEq(t, handledIDSnapshotBefore, handledIDSnapshotAfter, "RootID should be the same after import")
	})
}

func Test_RunImportArtrans(t *testing.T) {
	t.Run("Read JSON file test", func(t *testing.T) {
		// init a new clean db without any records
		ddb, _ := db.NewTestDB()
		defer db.CloseDB(ddb)
		dao := dao.NewDao(ddb)

		// create temp json file
		tmp, _ := os.CreateTemp("", "artrans.*.json")
		defer os.Remove(tmp.Name())

		_, _ = tmp.WriteString(`[{ "content": "TestContent", "page_key": "/test_page_key.html", "site_name": "test_site" }]`)
		_ = tmp.Close()

		// run import
		err := RunImportArtrans(dao, &ImportParams{
			JsonFile:  tmp.Name(),
			Assumeyes: true,
		})
		assert.NoError(t, err)
	})

	t.Run("Read JSON data test", func(t *testing.T) {
		// init a new clean db without any records
		ddb, _ := db.NewTestDB()
		defer db.CloseDB(ddb)
		dao := dao.NewDao(ddb)

		// run import
		err := RunImportArtrans(dao, &ImportParams{
			JsonData:  `[{"content": "TestContent", "page_key": "/test_page_key.html", "site_name": "test_site"}]`,
			Assumeyes: true,
		})
		assert.NoError(t, err)
	})

	t.Run("No data test", func(t *testing.T) {
		// init a new clean db without any records
		ddb, _ := db.NewTestDB()
		defer db.CloseDB(ddb)
		dao := dao.NewDao(ddb)

		// run import
		err := RunImportArtrans(dao, &ImportParams{
			Assumeyes: true,
		})
		assert.Error(t, err)
	})

	t.Run("Set console output test", func(t *testing.T) {
		// init a new clean db without any records
		ddb, _ := db.NewTestDB()
		defer db.CloseDB(ddb)
		dao := dao.NewDao(ddb)

		// run import
		called := false
		err := RunImportArtrans(dao, &ImportParams{
			JsonData:  `[{"content": "TestContent", "page_key": "/test_page_key.html", "site_name": "test_site"}]`,
			Assumeyes: true,
		}, func(s string) {
			called = true
			fmt.Println(s)
		})
		assert.NoError(t, err)
		assert.True(t, called, "Output function should be called")
	})
}

func Test_rebuildAutoIncrementIDs(t *testing.T) {
	jsonStr := `[
		{"id": "100", "rid": "", "root_id": 0},
		{"id": 200, "rid": "100", "root_id": "0"},
		{"id": "300", "rid": 200, "root_id": "100"},
		{"id": 400, "rid": "300", "root_id": "200"},
		{"id": "500", "rid": 400, "root_id": 300}
	]`

	expected := `[
		{"id":"1","rid":"0","root_id":"0"},
		{"id":"2","rid":"1","root_id":"0"},
		{"id":"3","rid":"2","root_id":"1"},
		{"id":"4","rid":"3","root_id":"2"},
		{"id":"5","rid":"4","root_id":"3"}
	]`

	assert.JSONEq(t, expected, rebuildAutoIncrementIDs(jsonStr))
}

func rebuildAutoIncrementIDs(jsonStr string) string {
	// current id
	id := 1
	// map to store the old id and new id
	idMap := map[string]string{}
	// use regex to replace ID values
	re := regexp.MustCompile(`"id":\s*"?(\d+)?"?`)
	jsonStr = re.ReplaceAllStringFunc(jsonStr, func(s string) string {
		newID := re.FindStringSubmatch(s)[1]
		idMap[newID] = fmt.Sprint(id)
		v := fmt.Sprintf(`"id":"%d"`, id)
		id++
		return v
	})

	// update rid and root_id values
	update := func(key string) {
		re := regexp.MustCompile(fmt.Sprintf(`"%s":\s*"?(\d+)?"?`, key))
		jsonStr = re.ReplaceAllStringFunc(jsonStr, func(s string) string {
			newID := re.FindStringSubmatch(s)[1]
			if newID == "0" || newID == "" {
				return fmt.Sprintf(`"%s":"0"`, key)
			}
			return fmt.Sprintf(`"%s":"%s"`, key, idMap[newID])
		})
	}
	update("rid")
	update("root_id")

	return jsonStr
}
