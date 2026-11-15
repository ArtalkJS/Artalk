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
	"runtime"
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
		req := multipartUpload(t, "/upload", "test.png", content, token)
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

func multipartUpload(t *testing.T, path, filename string, content []byte, token string) *http.Request {
	t.Helper()
	var body bytes.Buffer
	form := multipart.NewWriter(&body)
	file, err := form.CreateFormFile("file", filename)
	require.NoError(t, err)
	_, err = file.Write(content)
	require.NoError(t, err)
	require.NoError(t, form.Close())
	req := httptest.NewRequest(http.MethodPost, path, &body)
	req.Header.Set("Content-Type", form.FormDataContentType())
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req
}

func TestUpgitUploadCleansUpLocalFiles(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the local Upgit fixture uses a POSIX shell")
	}
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	app.Conf().Captcha.Enabled = false
	upload := &app.Conf().ImgUpload
	upload.Enabled = true
	upload.Upgit.Enabled = true
	upload.Upgit.Exec = "upgit --test"
	router.Use(limiter.ActionLimitMiddleware(app.App, limiter.ActionLimitConf{}))
	handler.Upload(app.App, router)

	bin := t.TempDir()
	t.Setenv("PATH", bin)
	// Use a disposable Upgit stand-in with no external storage.
	const script = `#!/bin/sh
test "$1" = "--test" && test -s "$2" || exit 1
case "$ARTALK_TEST_UPGIT" in
    failed) exit 1 ;;
    invalid) printf '%s\n' 'not a URL' ;;
    *) printf '%s\n' 'https://files.example.org/image.png' ;;
esac
`
	require.NoError(t, os.WriteFile(filepath.Join(bin, "upgit"), []byte(script), 0700))
	var img bytes.Buffer
	require.NoError(t, png.Encode(&img, image.NewRGBA(image.Rect(0, 0, 1, 1))))
	for _, tc := range []struct {
		name      string
		delete    bool
		remaining int
		status    int
	}{
		{"keep", false, 1, http.StatusOK},
		{"delete", true, 0, http.StatusOK},
		{"invalid", false, 0, http.StatusInternalServerError},
		{"failed", false, 0, http.StatusInternalServerError},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("ARTALK_TEST_UPGIT", tc.name)
			upload.Path = t.TempDir()
			upload.Upgit.DelLocal = tc.delete
			req := multipartUpload(t, "/upload", "test.png", img.Bytes(), "")
			resp, err := router.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Equal(t, tc.status, resp.StatusCode)
			entries, err := os.ReadDir(upload.Path)
			require.NoError(t, err)
			require.Len(t, entries, tc.remaining)
			if tc.status == http.StatusOK {
				var result handler.ResponseUpload
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
				require.Equal(t, "https://files.example.org/image.png", result.PublicURL)
			}
		})
	}
}

func TestUploadStorageFailures(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	app.Conf().Captcha.Enabled = false
	app.Conf().ImgUpload.Enabled = true
	router.Use(limiter.ActionLimitMiddleware(app.App, limiter.ActionLimitConf{}))
	handler.Upload(app.App, router)
	var img bytes.Buffer
	require.NoError(t, png.Encode(&img, image.NewRGBA(image.Rect(0, 0, 1, 1))))
	blocked := filepath.Join(t.TempDir(), "occupied")
	require.NoError(t, os.WriteFile(blocked, []byte("keep"), 0600))
	for _, path := range []string{blocked, filepath.Join(blocked, "child")} {
		app.Conf().ImgUpload.Path = path
		response, err := router.Test(multipartUpload(t, "/upload", "test.png", img.Bytes(), ""))
		require.NoError(t, err)
		response.Body.Close()
		require.Equal(t, 500, response.StatusCode)
	}
	contents, err := os.ReadFile(blocked)
	require.NoError(t, err)
	require.Equal(t, "keep", string(contents))
	flowRequest(t, router, "POST", "/upload", map[string]any{}, "", 500)
}
