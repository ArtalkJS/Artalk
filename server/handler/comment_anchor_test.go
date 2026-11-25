package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/artalkjs/artalk/v2/test"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func newAnnotationRouter(t *testing.T) (*test.TestApp, *fiber.App) {
	t.Helper()
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	handler.CommentCreate(app.App, router)
	handler.CommentList(app.App, router)
	handler.CommentAnchor(app.App, router)
	handler.CommentDelete(app.App, router)
	return app, router
}

func userToken(t *testing.T, app *test.TestApp, userID uint) string {
	t.Helper()
	token, err := common.LoginGetUserToken(app.Dao().FindUserByID(userID), app.Conf().AppKey, 60)
	require.NoError(t, err)
	return token
}

func annotationBody(pageKey, exact string, rid uint) map[string]any {
	return map[string]any{
		"name":      "userA",
		"email":     "user_a@qwqaq.com",
		"content":   "annotation " + exact,
		"page_key":  pageKey,
		"site_name": "Site A",
		"rid":       rid,
		"anchor": map[string]any{
			"quote": map[string]any{
				"type":   "TextQuoteSelector",
				"exact":  exact,
				"prefix": "pre",
				"suffix": "suf",
			},
			"position": map[string]any{"start": 1, "end": 1 + len(exact)},
			"section":  "Intro",
		},
	}
}

func createAnnotation(t *testing.T, router *fiber.App, token, pageKey, exact string) entity.CommentDTO {
	t.Helper()
	var dto entity.CommentDTO
	require.NoError(t, json.Unmarshal(
		flowRequest(t, router, http.MethodPost, "/comments", annotationBody(pageKey, exact, 0), token, http.StatusOK),
		&dto,
	))
	return dto
}

func TestCommentAnchorCreate(t *testing.T) {
	app, router := newAnnotationRouter(t)
	token := userToken(t, app, 1001)

	dto := createAnnotation(t, router, token, "/test/1000.html", "Hello Artalk")
	require.True(t, dto.IsAnnotation)
	require.NotNil(t, dto.Anchor)
	require.Equal(t, "Hello Artalk", dto.Anchor.Quote.Exact)
	require.Equal(t, "pre", dto.Anchor.Quote.Prefix)
	require.Equal(t, "suf", dto.Anchor.Quote.Suffix)
	require.Equal(t, "Intro", dto.Anchor.Section)
	require.NotNil(t, dto.Anchor.Position)
	require.Equal(t, 1, dto.Anchor.Position.Start)
	require.False(t, dto.AnchorResolved)

	replyBody := map[string]any{
		"name":      "userA",
		"email":     "user_a@qwqaq.com",
		"content":   "reply to annotation",
		"page_key":  "/test/1000.html",
		"site_name": "Site A",
		"rid":       dto.ID,
	}
	var reply entity.CommentDTO
	require.NoError(t, json.Unmarshal(
		flowRequest(t, router, http.MethodPost, "/comments", replyBody, token, http.StatusOK),
		&reply,
	))
	require.False(t, reply.IsAnnotation)

	var list handler.ResponseCommentList
	require.NoError(t, json.Unmarshal(
		flowRequest(t, router, http.MethodGet, "/comments?page_key=/test/1000.html&site_name=Site%20A&has_anchor=true&flat_mode=true&limit=100", nil, "", http.StatusOK),
		&list,
	))
	require.Equal(t, int64(2), list.Count)
	require.Equal(t, int64(1), list.RootsCount)
	require.Len(t, list.Comments, 2)
}

func TestCommentAnchorValidation(t *testing.T) {
	app, router := newAnnotationRouter(t)
	token := userToken(t, app, 1001)

	cases := []struct {
		name string
		body map[string]any
	}{
		{
			name: "empty exact",
			body: annotationBody("/test/1000.html", "   ", 0),
		},
		{
			name: "too long exact",
			body: annotationBody("/test/1000.html", strings.Repeat("a", 2001), 0),
		},
		{
			name: "invalid position",
			body: func() map[string]any {
				b := annotationBody("/test/1000.html", "valid", 0)
				b["anchor"].(map[string]any)["position"] = map[string]any{"start": 5, "end": 1}
				return b
			}(),
		},
		{
			name: "anchor on reply",
			body: annotationBody("/test/1000.html", "valid", 1000),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			flowRequest(t, router, http.MethodPost, "/comments", tc.body, token, http.StatusBadRequest)
		})
	}
}

func TestCommentListHasAnchor(t *testing.T) {
	app, router := newAnnotationRouter(t)
	token := userToken(t, app, 1001)
	createAnnotation(t, router, token, "/test/1000.html", "list filter")

	var all, onlyAnnotations, onlyNormal handler.ResponseCommentList
	require.NoError(t, json.Unmarshal(
		flowRequest(t, router, http.MethodGet, "/comments?page_key=/test/1000.html&site_name=Site%20A&flat_mode=true&limit=100", nil, "", http.StatusOK),
		&all,
	))
	require.NoError(t, json.Unmarshal(
		flowRequest(t, router, http.MethodGet, "/comments?page_key=/test/1000.html&site_name=Site%20A&has_anchor=true&flat_mode=true&limit=100", nil, "", http.StatusOK),
		&onlyAnnotations,
	))
	require.NoError(t, json.Unmarshal(
		flowRequest(t, router, http.MethodGet, "/comments?page_key=/test/1000.html&site_name=Site%20A&has_anchor=false&flat_mode=true&limit=100", nil, "", http.StatusOK),
		&onlyNormal,
	))

	require.Equal(t, all.Count, onlyAnnotations.Count+onlyNormal.Count)
	require.Equal(t, int64(1), onlyAnnotations.Count)
	require.Equal(t, int64(1), onlyAnnotations.RootsCount)
	require.Equal(t, int64(6), onlyNormal.Count)
	require.Equal(t, int64(2), onlyNormal.RootsCount)

	flowRequest(t, router, http.MethodGet, "/comments?page_key=/test/1000.html&site_name=Site%20A&has_anchor=maybe", nil, "", http.StatusBadRequest)
}

func TestCommentAnchorPermission(t *testing.T) {
	app, router := newAnnotationRouter(t)
	userAToken := userToken(t, app, 1001)
	userBToken := userToken(t, app, 1002)
	adminToken := userToken(t, app, 1000)

	dto := createAnnotation(t, router, userAToken, "/test/1000.html", "permission")

	path := fmt.Sprintf("/comments/%d/anchor", dto.ID)
	flowRequest(t, router, http.MethodPut, path, map[string]any{"resolved": true}, userBToken, http.StatusForbidden)

	var updated entity.CommentDTO
	require.NoError(t, json.Unmarshal(
		flowRequest(t, router, http.MethodPut, path, map[string]any{"resolved": true}, userAToken, http.StatusOK),
		&updated,
	))
	require.True(t, updated.AnchorResolved)

	require.NoError(t, json.Unmarshal(
		flowRequest(t, router, http.MethodPut, path, map[string]any{"resolved": false}, adminToken, http.StatusOK),
		&updated,
	))
	require.False(t, updated.AnchorResolved)

	require.NoError(t, json.Unmarshal(
		flowRequest(t, router, http.MethodPut, path, map[string]any{
			"anchor": map[string]any{
				"quote": map[string]any{"type": "TextQuoteSelector", "exact": "relocated"},
			},
		}, userAToken, http.StatusOK),
		&updated,
	))
	require.Equal(t, "relocated", updated.Anchor.Quote.Exact)

	flowRequest(t, router, http.MethodDelete, fmt.Sprintf("/comments/%d", dto.ID), nil, userBToken, http.StatusForbidden)
	flowRequest(t, router, http.MethodDelete, fmt.Sprintf("/comments/%d", dto.ID), nil, userAToken, http.StatusOK)

	// Normal comments remain admin-only for deletion.
	flowRequest(t, router, http.MethodDelete, "/comments/1001", nil, userAToken, http.StatusForbidden)
	flowRequest(t, router, http.MethodDelete, "/comments/1001", nil, adminToken, http.StatusOK)
}

func TestCommentAnchorGuest(t *testing.T) {
	_, router := newAnnotationRouter(t)

	body := annotationBody("/test/1000.html", "guest anchor", 0)
	body["name"] = "guest"
	body["email"] = "guest@example.org"

	resp := guestRequest(t, router, http.MethodPost, "/comments", guestA, "", body)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var dto entity.CommentDTO
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&dto))
	require.True(t, dto.IsAnnotation)

	path := fmt.Sprintf("/comments/%d/anchor", dto.ID)
	require.Equal(t, http.StatusForbidden, guestRequest(t, router, http.MethodPut, path, guestB, "", map[string]any{"resolved": true}).StatusCode)
	require.Equal(t, http.StatusOK, guestRequest(t, router, http.MethodPut, path, guestA, "", map[string]any{"resolved": true}).StatusCode)
	require.Equal(t, http.StatusOK, guestRequest(t, router, http.MethodDelete, fmt.Sprintf("/comments/%d", dto.ID), guestA, "", nil).StatusCode)
}
