package handler_test

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/stretchr/testify/require"
)

func TestPrivateNotificationsRequireIdentity(t *testing.T) {
	app, router := NewApiTestApp()
	defer app.Cleanup()
	handler.UserInfo(app.App, router)
	handler.NotifyList(app.App, router)
	handler.NotifyReadAll(app.App, router)
	victim := app.Dao().FindUserByID(1001)
	query := "?name=" + url.QueryEscape(victim.Name) + "&email=" + url.QueryEscape(victim.Email)
	for _, userID := range []uint{0, 1002, 1000, 1001} {
		t.Run(fmt.Sprint(userID), func(t *testing.T) {
			token := ""
			if userID != 0 {
				var err error
				token, err = common.LoginGetUserToken(app.Dao().FindUserByID(userID), app.Conf().AppKey, 60)
				require.NoError(t, err)
			}
			resp, err := securityRequest(router, "GET", "/notifies"+query, nil, token, "")
			require.NoError(t, err)
			var notices handler.ResponseNotifyList
			require.NoError(t, json.NewDecoder(resp.Body).Decode(&notices))
			resp.Body.Close()
			require.Equal(t, 200, resp.StatusCode)
			for _, notice := range notices.Notifies {
				require.Equal(t, userID, notice.UserID)
			}
			if userID == 0 {
				require.Empty(t, notices.Notifies)
			}
			if userID == victim.ID {
				require.NotEmpty(t, notices.Notifies)
			}
			resp, err = securityRequest(router, "GET", "/user"+query, nil, token, "")
			require.NoError(t, err)
			var info handler.ResponseUserInfo
			require.NoError(t, json.NewDecoder(resp.Body).Decode(&info))
			resp.Body.Close()
			if userID == 0 {
				require.Empty(t, info.Notifies)
				require.Zero(t, info.NotifiesCount)
				require.NotNil(t, info.User) // retain public link autocomplete and admin login hints
				require.Empty(t, info.User.Email)
				require.False(t, info.User.ReceiveEmail)
			} else {
				require.Equal(t, userID, info.User.ID)
			}
			resp, err = securityRequest(router, "POST", "/notifies/read", map[string]string{"name": victim.Name, "email": victim.Email}, token, "")
			require.NoError(t, err)
			resp.Body.Close()
			if userID == 0 {
				require.Equal(t, 401, resp.StatusCode)
			} else {
				require.Equal(t, 200, resp.StatusCode)
			}
			var notify entity.Notify
			require.NoError(t, app.Dao().DB().First(&notify, 1002).Error)
			require.Equal(t, userID == victim.ID, notify.IsRead, "request fields must not select the notification owner")
		})
	}
	// The report's name-only/email-only enumeration does not work in current code.
	for _, query := range []string{"?name=" + victim.Name, "?email=" + url.QueryEscape(victim.Email)} {
		resp, err := securityRequest(router, "GET", "/user"+query, nil, "", "")
		require.NoError(t, err)
		var info handler.ResponseUserInfo
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&info))
		resp.Body.Close()
		require.Nil(t, info.User)
	}
}

func TestNotificationCapabilityLinks(t *testing.T) {
	app, router := NewApiTestApp()
	defer app.Cleanup()
	handler.NotifyRead(app.App, router)
	notify := entity.Notify{UserID: 1001, CommentID: 1005}
	notify.GenerateKey()
	require.GreaterOrEqual(t, len(notify.Key), 32)
	require.NoError(t, app.Dao().DB().Create(&notify).Error)
	for _, tc := range []struct {
		key    string
		status int
	}{{"wrong", 404}, {notify.Key, 200}, {notify.Key, 200}} {
		resp, err := securityRequest(router, "POST", "/notifies/1005/"+tc.key, nil, "", "")
		require.NoError(t, err)
		resp.Body.Close()
		require.Equal(t, tc.status, resp.StatusCode)
	}
	// Previously sent email links stay usable.
	resp, err := securityRequest(router, "POST", "/notifies/1000/ZMTM4", nil, "", "")
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, 200, resp.StatusCode)
}

func TestCommentListCannotImpersonatePendingAuthor(t *testing.T) {
	app, router := NewApiTestApp()
	defer app.Cleanup()
	handler.CommentList(app.App, router)
	require.NoError(t, app.Dao().DB().Model(&entity.Comment{}).Where("id = ?", 1001).Update("is_pending", true).Error)
	victim := app.Dao().FindUserByID(1001)
	for _, userID := range []uint{0, 1002, 1001, 1000} {
		token := ""
		if userID != 0 {
			var err error
			token, err = common.LoginGetUserToken(app.Dao().FindUserByID(userID), app.Conf().AppKey, 60)
			require.NoError(t, err)
		}
		for _, scope := range []string{"page", "user"} {
			query := url.Values{"page_key": {"/test/1000.html"}, "site_name": {"Site A"}, "name": {victim.Name}, "email": {victim.Email}, "scope": {scope}, "flat_mode": {"true"}, "limit": {"100"}}
			resp, err := securityRequest(router, "GET", "/comments?"+query.Encode(), nil, token, "")
			require.NoError(t, err)
			require.Equal(t, 200, resp.StatusCode)
			var result handler.ResponseCommentList
			require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
			resp.Body.Close()
			found := false
			for _, comment := range result.Comments {
				if comment.ID == 1001 {
					found = true
				}
			}
			if userID == 0 || userID == 1002 {
				require.False(t, found)
			}
			if userID == 1001 || (userID == 1000 && scope == "page") {
				require.True(t, found)
			}
		}
	}
}

func TestRevokedTokenCannotReadPrivateData(t *testing.T) {
	app, router := NewApiTestApp()
	defer app.Cleanup()
	handler.UserInfo(app.App, router)
	handler.NotifyList(app.App, router)
	handler.NotifyReadAll(app.App, router)
	handler.CommentList(app.App, router)
	user := app.Dao().FindUserByID(1000)
	token, err := common.LoginGetUserToken(user, app.Conf().AppKey, 60)
	require.NoError(t, err)
	require.NoError(t, user.TokenValidFrom.Scan(time.Now().Add(time.Minute)))
	require.NoError(t, app.Dao().UpdateUser(&user))
	require.NoError(t, app.Dao().DB().Model(&entity.Comment{}).Where("id = ?", 1001).Update("is_pending", true).Error)
	for _, path := range []string{"/user", "/notifies"} {
		resp, err := securityRequest(router, "GET", path, nil, token, "")
		require.NoError(t, err)
		var result map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		resp.Body.Close()
		require.Empty(t, result["notifies"])
		if path == "/user" {
			require.Nil(t, result["user"])
			require.Equal(t, false, result["is_login"])
		}
	}
	resp, err := securityRequest(router, "POST", "/notifies/read", map[string]string{"name": user.Name, "email": user.Email}, token, "")
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, 401, resp.StatusCode)
	resp, err = securityRequest(router, "GET", "/comments?page_key=/test/1000.html&site_name=Site%20A&flat_mode=true", nil, token, "")
	require.NoError(t, err)
	var result handler.ResponseCommentList
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	resp.Body.Close()
	for _, comment := range result.Comments {
		require.False(t, comment.IsPending)
	}
}
