package handler_test

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/artalkjs/artalk/v2/server/middleware/limiter"
	"github.com/stretchr/testify/require"
)

func TestImageUpload(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	app.Conf().Captcha.Enabled = false
	upload := &app.Conf().ImgUpload
	upload.Path = t.TempDir()
	upload.MaxSize = 1
	upload.Enabled = false
	router.Use(limiter.ActionLimitMiddleware(app.App, limiter.ActionLimitConf{}))
	handler.Upload(app.App, router)

	var img bytes.Buffer
	require.NoError(t, png.Encode(&img, image.NewRGBA(image.Rect(0, 0, 1, 1))))
	post := func(content []byte, token string, status int) handler.ResponseUpload {
		t.Helper()
		var body bytes.Buffer
		form := multipart.NewWriter(&body)
		file, err := form.CreateFormFile("file", "test.png")
		require.NoError(t, err)
		_, err = file.Write(content)
		require.NoError(t, err)
		require.NoError(t, form.Close())
		req := httptest.NewRequest(http.MethodPost, "/upload", &body)
		req.Header.Set("Content-Type", form.FormDataContentType())
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		res, err := router.Test(req)
		require.NoError(t, err)
		defer res.Body.Close()
		require.Equal(t, status, res.StatusCode)
		var result handler.ResponseUpload
		if status == http.StatusOK {
			require.NoError(t, json.NewDecoder(res.Body).Decode(&result))
		}
		return result
	}

	post(img.Bytes(), "", http.StatusForbidden)
	upload.Enabled = true
	post([]byte("plain text"), "", http.StatusBadRequest)
	post(bytes.Repeat([]byte{'x'}, 1024*1024+1), "", http.StatusBadRequest)
	entries, err := os.ReadDir(upload.Path)
	require.NoError(t, err)
	require.Empty(t, entries)

	upload.Enabled = false
	admin, err := common.LoginGetUserToken(app.Dao().FindUserByID(1000), app.Conf().AppKey, 60)
	require.NoError(t, err)
	upload.PublicPath = "https://images.example.org/comments/"
	result := post(img.Bytes(), admin, http.StatusOK)
	require.Equal(t, "image", result.FileType)
	require.Equal(t, ".png", filepath.Ext(result.FileName))
	require.Equal(t, "https://images.example.org/comments/"+result.FileName, result.PublicURL)
	stored, err := os.ReadFile(filepath.Join(upload.Path, result.FileName))
	require.NoError(t, err)
	require.Equal(t, img.Bytes(), stored)
}
