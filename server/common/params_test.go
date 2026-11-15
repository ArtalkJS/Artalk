package common

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestRequestParameters(t *testing.T) {
	type params struct {
		ID      int    `params:"id"`
		Name    string `json:"name" query:"name" validate:"required"`
		Limit   int    `json:"limit" query:"limit"`
		Enabled bool   `json:"enabled" query:"enabled" validate:"required"`
	}
	router := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
	router.All("/input/:id?", func(c *fiber.Ctx) error {
		var p params
		if ok, err := ParamsDecode(c, &p); !ok {
			return err
		}
		return c.JSON(p)
	})
	for _, tc := range []struct {
		name, method, path, contentType, body string
		status                                int
		want                                  params
	}{
		{"query", "GET", "/input?name=Reader&enabled=false&limit=3", "", "", 200, params{Name: "Reader", Limit: 3}},
		{"JSON", "POST", "/input", "application/json", `{"name":"Reader","enabled":false,"limit":3}`, 200, params{Name: "Reader", Limit: 3}},
		{"missing required", "GET", "/input?enabled=false", "", "", 400, params{}},
		{"null required", "POST", "/input", "application/json", `{"name":null,"enabled":null}`, 200, params{}},
		{"invalid query type", "GET", "/input?name=Reader&enabled=false&limit=no", "", "", 400, params{}},
		{"invalid path type", "GET", "/input/no?name=Reader&enabled=false", "", "", 400, params{}},
		{"invalid JSON type", "POST", "/input", "application/json", `{"name":"Reader","enabled":false,"limit":"no"}`, 400, params{}},
		{"malformed JSON", "POST", "/input", "application/json", `{`, 400, params{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			if tc.contentType != "" {
				req.Header.Set("Content-Type", tc.contentType)
			}
			resp, err := router.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			data, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, tc.status, resp.StatusCode, string(data))
			if tc.status == 200 {
				var parsed params
				require.NoError(t, json.Unmarshal(data, &parsed))
				require.Equal(t, tc.want, parsed)
			} else {
				var failure map[string]any
				require.NoError(t, json.Unmarshal(data, &failure))
				require.NotEmpty(t, failure["msg"])
			}
		})
	}
}
