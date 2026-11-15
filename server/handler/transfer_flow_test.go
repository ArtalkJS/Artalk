package handler_test

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/artransfer"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/stretchr/testify/require"
)

func TestAdminTransferUploadImportAndExport(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	t.Setenv("TMPDIR", t.TempDir())
	handler.Transfer(app.App, router)
	admin, err := common.LoginGetUserToken(app.Dao().FindUserByID(1000), app.Conf().AppKey, 60)
	require.NoError(t, err)

	input, err := json.Marshal([]entity.Artran{{
		ID: "1", Rid: "0", Content: "Transferred comment", SiteName: "Transfer site",
		PageKey: "/transferred", PageTitle: "Transferred page", Nick: "Reader", Email: "transfer@example.org",
	}})
	require.NoError(t, err)
	upload := func(token string) *http.Response {
		t.Helper()
		resp, err := router.Test(multipartUpload(t, "/transfer/upload", "test.artrans", input, token))
		require.NoError(t, err)
		t.Cleanup(func() { resp.Body.Close() })
		return resp
	}
	require.Equal(t, http.StatusForbidden, upload("").StatusCode)
	resp := upload(admin)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var file handler.ResponseTransferUpload
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&file))
	saved, err := os.ReadFile(file.Filename)
	require.NoError(t, err)
	require.Equal(t, input, saved)

	params := handler.ParamsTransferImport{ImportParams: artransfer.ImportParams{JsonFile: file.Filename}}
	flowRequest(t, router, http.MethodPost, "/transfer/import", params, admin, http.StatusOK)
	require.NotZero(t, app.Dao().FindPage("/transferred", "Transfer site").ID)

	flowRequest(t, router, http.MethodGet, "/transfer/export", nil, "", http.StatusForbidden)
	var exported handler.ResponseTransferExport
	require.NoError(t, json.Unmarshal(flowRequest(t, router, http.MethodGet, "/transfer/export", nil, admin, http.StatusOK), &exported))
	var rows []entity.Artran
	require.NoError(t, json.Unmarshal([]byte(exported.Artrans), &rows))
	for _, row := range rows {
		if row.SiteName == "Transfer site" && row.PageKey == "/transferred" {
			require.Equal(t, "Transferred comment", row.Content)
			require.Equal(t, "Transferred page", row.PageTitle)
			return
		}
	}
	t.Fatal("imported comment missing from the export")
}
