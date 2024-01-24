package handler_test

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/ArtalkJS/Artalk/server/handler"
	"github.com/ArtalkJS/Artalk/test"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestVersionApi(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	t.Run("Version", func(t *testing.T) {
		fiberApp := fiber.New()
		handler.Version(app.App, fiberApp)

		config.Version = "v0.0.1"
		config.CommitHash = "0000000"

		req := httptest.NewRequest("GET", "/version", nil)
		resp, _ := fiberApp.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		var data common.ApiVersionData
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &data)

		assert.Equal(t, "artalk", data.App)
		assert.Equal(t, "0.0.1", data.Version)
		assert.Equal(t, "0000000", data.CommitHash)
	})
}
