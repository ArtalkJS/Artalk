package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestManagementRejectsIncompleteRequests(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	for _, register := range []func(*core.App, fiber.Router){
		handler.UserCreate, handler.UserUpdate, handler.UserInfoUpdate, handler.SiteCreate, handler.SiteUpdate,
		handler.PageUpdate, handler.SettingApply, handler.EmailSend, handler.AuthMergeApply,
		handler.AuthEmailRegister, handler.AuthEmailLogin, handler.AuthEmailSend, handler.CaptchaVerify,
	} {
		register(app.App, router)
	}
	app.Conf().Auth.Email.Enabled = true
	app.Conf().Captcha.Enabled = false
	admin, err := common.LoginGetUserToken(app.Dao().FindUserByID(1000), app.Conf().AppKey, 60)
	require.NoError(t, err)
	var before int64
	require.NoError(t, app.Dao().DB().Model(&entity.User{}).Count(&before).Error)
	for _, tc := range []struct{ method, path string }{
		{"POST", "/users"}, {"PUT", "/users/1001"}, {"POST", "/user"},
		{"POST", "/sites"}, {"PUT", "/sites/1000"}, {"PUT", "/pages/1000"},
		{"PUT", "/settings"}, {"POST", "/send_email"}, {"POST", "/auth/merge"},
		{"POST", "/auth/email/register"}, {"POST", "/auth/email/login"}, {"POST", "/auth/email/send"}, {"POST", "/captcha/verify"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			flowRequest(t, router, tc.method, tc.path, map[string]any{}, admin, http.StatusBadRequest)
		})
	}
	var after int64
	require.NoError(t, app.Dao().DB().Model(&entity.User{}).Count(&after).Error)
	require.Equal(t, before, after)
}

func TestManagementTemplatesAndStatus(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	handler.SettingTemplate(app.App, router)
	handler.UserStatus(app.App, router)
	handler.ConfDomain(app.App, router)
	handler.EmailSend(app.App, router)
	handler.CacheFlush(app.App, router)
	handler.CacheWarmUp(app.App, router)
	user := app.Dao().FindUserByID(1000)
	admin, err := common.LoginGetUserToken(user, app.Conf().AppKey, 60)
	require.NoError(t, err)
	for _, path := range []string{"/settings/template", "/settings/template/en", "/settings/template/unknown"} {
		flowRequest(t, router, "GET", path, nil, "", 403)
		var response handler.ResponseSettingTemplate
		require.NoError(t, json.Unmarshal(flowRequest(t, router, "GET", path, nil, admin, 200), &response))
		require.Contains(t, response.Yaml, "site_default:")
	}
	for _, token := range []string{"", admin} {
		var status handler.ResponseUserStatus
		require.NoError(t, json.Unmarshal(flowRequest(t, router, "GET", "/user/status?name="+user.Name+"&email="+user.Email, nil, token, 200), &status))
		require.True(t, status.IsAdmin)
		require.Equal(t, token != "", status.IsLogin)
	}
	app.Conf().TrustedDomains = []string{"https://trusted.example.org"}
	for _, path := range []string{"/conf/domain", "/conf/domain?url=https://trusted.example.org/post", "/conf/domain?url=https://untrusted.example.org"} {
		var domain handler.ResponseConfDomain
		require.NoError(t, json.Unmarshal(flowRequest(t, router, "GET", path, nil, "", 200), &domain))
		require.Equal(t, path == "/conf/domain?url=https://trusted.example.org/post", domain.IsTrusted)
	}
	app.Conf().Email.Enabled = false
	flowRequest(t, router, "POST", "/send_email", handler.ParamsEmailSend{ToAddr: "reader@example.org", Subject: "Test", Body: "Body"}, admin, 200)
	app.Conf().Cache.Enabled = false
	for _, path := range []string{"/cache/flush", "/cache/warm_up"} {
		flowRequest(t, router, "POST", path, nil, admin, 400)
	}
}
