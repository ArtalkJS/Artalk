package handler_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/artalkjs/artalk/v2/server/middleware"
	"github.com/artalkjs/artalk/v2/server/middleware/limiter"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

var guestA = strings.Repeat("a1", 32)
var guestB = strings.Repeat("b2", 32)

func guestHash(token string) string {
	raw, _ := hex.DecodeString(token)
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}
func guestRequest(t *testing.T, router *fiber.App, method, path, guest, login string, body any) *http.Response {
	t.Helper()
	data, err := json.Marshal(body)
	require.NoError(t, err)
	req := httptest.NewRequest(method, path, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(common.GuestTokenHeader, guest)
	if login != "" {
		req.Header.Set("Authorization", "Bearer "+login)
	}
	resp, err := router.Test(req, 10000)
	require.NoError(t, err)
	t.Cleanup(func() { resp.Body.Close() })
	return resp
}

func TestGuestCommentOwnership(t *testing.T) {
	app, router := NewApiTestApp()
	defer app.Cleanup()
	handler.CommentGet(app.App, router)
	handler.CommentList(app.App, router)
	// Same name/email and user ID, but different browsers. Also include a legacy row.
	for id, hash := range map[uint]string{1001: guestHash(guestA), 1002: guestHash(guestB), 1003: ""} {
		require.NoError(t, app.Dao().DB().Model(&entity.Comment{}).Where("id = ?", id).Updates(map[string]any{"is_pending": true, "guest_token_hash": hash}).Error)
	}
	for _, tc := range []struct {
		name, guest, login string
		allowed            bool
	}{
		{"owner", guestA, "", true}, {"same identity other browser", guestB, "", false},
		{"old client", "", "", false}, {"hash is not a credential", guestHash(guestA), "", false},
		{"malformed", "not-a-secret", "", false}, {"invalid login cannot fall back", guestA, "invalid", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			query := "?name=userA&email=user_a@qwqaq.com"
			resp := guestRequest(t, router, "GET", "/comments/1001"+query, tc.guest, tc.login, nil)
			if tc.allowed {
				require.Equal(t, 200, resp.StatusCode)
			} else {
				require.Equal(t, 404, resp.StatusCode)
			}
			for _, flat := range []string{"true", "false"} {
				resp = guestRequest(t, router, "GET", "/comments"+query+"&page_key=/test/1000.html&site_name=Site%20A&limit=100&flat_mode="+flat, tc.guest, tc.login, nil)
				require.Equal(t, 200, resp.StatusCode)
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				require.NotContains(t, string(body), guestHash(guestA))
				require.NotContains(t, string(body), guestA)
				require.NotContains(t, string(body), "guest_token_hash")
				// Serialized nested/linked comments must respect the same visibility rule.
				require.Equal(t, tc.allowed, bytes.Contains(body, []byte("回复测试 1001")))
			}
		})
	}
	// User/pending sidebar scope shows only this browser's pending comments, with correct counts.
	resp := guestRequest(t, router, "GET", "/comments?page_key=/test/1000.html&scope=user&type=pending&flat_mode=true&limit=1&name=userA&email=user_a@qwqaq.com", guestA, "", nil)
	var list handler.ResponseCommentList
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&list))
	require.EqualValues(t, 1, list.Count)
	require.Len(t, list.Comments, 1)
	require.EqualValues(t, 1001, list.Comments[0].ID)
	// Credentials cannot unlock site/admin scope or legacy pending rows.
	require.Equal(t, 404, guestRequest(t, router, "GET", "/comments/1003", guestA, "", nil).StatusCode)
	resp = guestRequest(t, router, "GET", "/comments?page_key=x&scope=site", guestA, "", nil)
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&list))
	require.Empty(t, list.Comments)
	// A proof must not bypass page filters or disclose another browser's parent.
	resp = guestRequest(t, router, "GET", "/comments?page_key=/different-page&site_name=Site%20A&limit=100", guestA, "", nil)
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&list))
	require.Zero(t, list.Count)
	require.Empty(t, list.Comments)
	resp = guestRequest(t, router, "GET", "/comments/1002", guestB, "", nil)
	var detail handler.ResponseCommentGet
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&detail))
	require.Nil(t, detail.ReplyComment)

	// The normal authenticated owner retains access to legacy comments.
	token, err := common.LoginGetUserToken(app.Dao().FindUserByID(1001), app.Conf().AppKey, 60)
	require.NoError(t, err)
	require.Equal(t, 200, guestRequest(t, router, "GET", "/comments/1003", guestB, token, nil).StatusCode)
}

func TestGuestNotificationsAndReadAll(t *testing.T) {
	app, router := NewApiTestApp()
	defer app.Cleanup()
	handler.NotifyList(app.App, router)
	handler.NotifyReadAll(app.App, router)
	handler.UserInfo(app.App, router)
	// Parent 1001 belongs to browser A; parent 1002 to browser B (same user ID).
	require.NoError(t, app.Dao().DB().Model(&entity.Comment{}).Where("id = ?", 1001).Update("guest_token_hash", guestHash(guestA)).Error)
	require.NoError(t, app.Dao().DB().Model(&entity.Comment{}).Where("id = ?", 1002).Update("guest_token_hash", guestHash(guestB)).Error)
	own := app.Dao().NewNotify(1001, 1002)   // reply to A
	other := app.Dao().NewNotify(1001, 1004) // reply to B
	admin := app.Dao().NewNotify(1000, 1002) // admin's copy of the same event
	for _, tc := range []struct {
		token string
		ids   []uint
	}{{guestA, []uint{own.ID}}, {guestB, []uint{other.ID}}, {"", nil}} {
		resp := guestRequest(t, router, "GET", "/notifies?name=admin&email=admin@qwqaq.com", tc.token, "", nil)
		var notices handler.ResponseNotifyList
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&notices))
		var ids []uint
		for _, n := range notices.Notifies {
			ids = append(ids, n.ID)
			require.EqualValues(t, 1001, n.UserID)
		}
		require.Equal(t, tc.ids, ids)
		resp = guestRequest(t, router, "GET", "/user?name=userA&email=user_a@qwqaq.com", tc.token, "", nil)
		var info handler.ResponseUserInfo
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&info))
		require.False(t, info.IsLogin)
		require.Equal(t, len(tc.ids), info.NotifiesCount)
	}
	resp := guestRequest(t, router, "POST", "/notifies/read", guestA, "", map[string]string{"name": "admin", "email": "admin@qwqaq.com"})
	require.Equal(t, 200, resp.StatusCode)
	for _, tc := range []struct {
		id   uint
		read bool
	}{{own.ID, true}, {other.ID, false}, {admin.ID, false}, {1002, false}} {
		var n entity.Notify
		require.NoError(t, app.Dao().DB().First(&n, tc.id).Error)
		require.Equal(t, tc.read, n.IsRead)
	}
	// No keys for pending replies or deleted parents, even if a notification exists.
	require.NoError(t, app.Dao().DB().Model(&entity.Comment{}).Where("id = ?", 1004).Update("is_pending", true).Error)
	resp = guestRequest(t, router, "GET", "/notifies", guestB, "", nil)
	var notices handler.ResponseNotifyList
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&notices))
	require.Empty(t, notices.Notifies)
	require.NoError(t, app.Dao().DB().Model(&entity.Comment{}).Where("id = ?", 1004).Update("is_pending", false).Error)
	require.NoError(t, app.Dao().DB().Delete(&entity.Comment{}, 1002).Error)
	resp = guestRequest(t, router, "GET", "/notifies", guestB, "", nil)
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&notices))
	require.Empty(t, notices.Notifies)
}

func TestGuestCreationAndLegacyCompatibility(t *testing.T) {
	app, router := NewApiTestApp()
	defer app.Cleanup()
	app.Conf().Captcha.Enabled = false
	app.Conf().Auth.Enabled = false
	app.Conf().Moderator.PendingDefault = true
	app.Conf().Email.Enabled = false
	router.Use(limiter.ActionLimitMiddleware(app.App, limiter.ActionLimitConf{}))
	handler.CommentCreate(app.App, router)
	handler.CommentGet(app.App, router)
	handler.PagePV(app.App, router)
	for i, guest := range []string{guestA, guestB, ""} {
		page := fmt.Sprintf("/guest-new-%d", i)
		resp := guestRequest(t, router, "POST", "/comments", guest, "", handler.ParamsCommentCreate{
			Name: "guest", Email: "guest@example.com", Content: "pending test", SiteName: "Site A", PageKey: page, PageTitle: "Guest test",
		})
		require.Equal(t, 200, resp.StatusCode)
		var created handler.ResponseCommentCreate
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&created))
		comment := app.Dao().FindComment(created.ID)
		if guest != "" {
			require.Equal(t, guestHash(guest), comment.GuestTokenHash)
		} else {
			require.Empty(t, comment.GuestTokenHash)
		}
		require.False(t, comment.IsVerified)
		require.True(t, comment.IsPending)
		require.False(t, app.Dao().FindPage(page, "Site A").IsEmpty())
		for pv := 1; pv <= 2; pv++ {
			resp = guestRequest(t, router, "POST", "/pages/pv", guest, "", handler.ParamsPagePV{PageKey: page, SiteName: "Site A"})
			require.Equal(t, 200, resp.StatusCode)
			var result handler.ResponsePagePV
			require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
			require.Equal(t, pv, result.PV)
		}
		expected := 200
		if guest == "" {
			expected = 404
		}
		require.Equal(t, expected, guestRequest(t, router, "GET", fmt.Sprintf("/comments/%d", created.ID), guest, "", nil).StatusCode)
		require.Equal(t, 404, guestRequest(t, router, "GET", fmt.Sprintf("/comments/%d?guest_token=%s", created.ID, url.QueryEscape(guest)), "", "", nil).StatusCode)
	}
}

func TestGuestHeaderCors(t *testing.T) {
	app, router := NewApiTestApp()
	defer app.Cleanup()
	app.Conf().TrustedDomains = []string{"https://blog.example.com"}
	router.Use(middleware.CorsMiddleware(app.App))
	req := httptest.NewRequest("OPTIONS", "/comments", nil)
	req.Header.Set("Origin", "https://blog.example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "content-type,x-artalk-guest")
	resp, err := router.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, 204, resp.StatusCode)
	require.Contains(t, strings.ToLower(resp.Header.Get("Access-Control-Allow-Headers")), "x-artalk-guest")
}
