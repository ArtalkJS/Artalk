package handler_test

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/stretchr/testify/require"
)

func TestConfigPublicViewAndAdminSettings(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	path := filepath.Join(t.TempDir(), "artalk.yml")
	const yaml = "app_key: fixture-secret\nlocale: ja\nimg_upload:\n  enabled: false\n"
	require.NoError(t, os.WriteFile(path, []byte(yaml), 0600))
	conf, err := config.NewFromFile(path)
	require.NoError(t, err)
	conf.Auth.Enabled = true
	app.SetConf(conf)
	handler.Conf(app.App, router)
	handler.SettingGet(app.App, router)
	admin, err := common.LoginGetUserToken(app.Dao().FindUserByID(1000), conf.AppKey, 60)
	require.NoError(t, err)

	public := flowRequest(t, router, http.MethodGet, "/conf", nil, "", http.StatusOK)
	require.NotContains(t, string(public), conf.AppKey)
	var data common.ConfData
	require.NoError(t, json.Unmarshal(public, &data))
	require.Equal(t, false, data.FrontendConf["imgUpload"])
	require.Equal(t, "ja", data.FrontendConf["locale"])
	require.ElementsMatch(t, []any{"dist/plugins/artalk-plugin-auth.js", "dist/i18n/ja.js"}, data.FrontendConf["pluginURLs"])

	require.NoError(t, json.Unmarshal(flowRequest(t, router, http.MethodGet, "/conf", nil, admin, http.StatusOK), &data))
	require.Equal(t, true, data.FrontendConf["imgUpload"])
	require.NoError(t, json.Unmarshal(flowRequest(t, router, http.MethodGet, "/conf", nil, "", http.StatusOK), &data))
	require.Equal(t, false, data.FrontendConf["imgUpload"])

	flowRequest(t, router, http.MethodGet, "/settings", nil, "", http.StatusForbidden)
	var settings handler.ResponseSettingGet
	require.NoError(t, json.Unmarshal(flowRequest(t, router, http.MethodGet, "/settings", nil, admin, http.StatusOK), &settings))
	require.Equal(t, yaml, settings.Yaml)
}
