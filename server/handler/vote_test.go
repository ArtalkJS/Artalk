package handler_test

import (
	"bytes"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestVote(t *testing.T) {
	tests := []struct {
		description  string
		method       string
		url          string
		body         string
		ip           string
		expectedCode int
		expectedBody func(t *testing.T, body string)
	}{
		{
			description:  "Get comment vote status (original up)",
			method:       "GET",
			url:          "/votes/comment/1000",
			expectedCode: 200,
			ip:           "127.0.0.1",
			expectedBody: func(t *testing.T, body string) {
				assert.NotEmpty(t, body)
				assert.Equal(t, `{"up":4,"down":2,"is_up":true,"is_down":false}`, body)
			},
		},
		{
			description:  "Get page vote status (original down)",
			method:       "GET",
			url:          "/votes/page/1001",
			expectedCode: 200,
			ip:           "127.0.0.1",
			expectedBody: func(t *testing.T, body string) {
				assert.NotEmpty(t, body)
				assert.Equal(t, `{"up":1,"down":2,"is_up":false,"is_down":true}`, body)
			},
		},
		{
			description:  "Get comment vote status (original null)",
			method:       "GET",
			url:          "/votes/comment/1001",
			expectedCode: 200,
			ip:           "127.0.0.1",
			expectedBody: func(t *testing.T, body string) {
				assert.NotEmpty(t, body)
				assert.Equal(t, `{"up":0,"down":0,"is_up":false,"is_down":false}`, body)
			},
		},
		{
			description:  "Get page vote status (original null)",
			method:       "GET",
			url:          "/votes/page/1002",
			expectedCode: 200,
			ip:           "127.0.0.1",
			expectedBody: func(t *testing.T, body string) {
				assert.NotEmpty(t, body)
				assert.Equal(t, `{"up":0,"down":0,"is_up":false,"is_down":false}`, body)
			},
		},
		{
			description:  "Reject malformed vote target ID",
			method:       "GET",
			url:          "/votes/comment/not-a-number",
			expectedCode: 400,
			expectedBody: func(t *testing.T, body string) {
				assert.JSONEq(t, `{"msg":"invalid vote target id"}`, body)
			},
		},
		{
			description:  "Reject non-positive vote target ID",
			method:       "GET",
			url:          "/votes/comment/0",
			expectedCode: 400,
			expectedBody: func(t *testing.T, body string) {
				assert.JSONEq(t, `{"msg":"invalid vote target id"}`, body)
			},
		},
		{
			description:  "Reject unknown vote target",
			method:       "GET",
			url:          "/votes/unknown/1000",
			expectedCode: 404,
			expectedBody: func(t *testing.T, body string) {
				assert.JSONEq(t, `{"msg":"unknown vote target name"}`, body)
			},
		},
		{
			description:  "Return not found for missing comment",
			method:       "GET",
			url:          "/votes/comment/999999",
			expectedCode: 404,
			expectedBody: func(t *testing.T, body string) {
				assert.Contains(t, body, "not found")
			},
		},
		{
			description:  "Return not found for missing page",
			method:       "GET",
			url:          "/votes/page/999999",
			expectedCode: 404,
			expectedBody: func(t *testing.T, body string) {
				assert.Contains(t, body, "not found")
			},
		},
		{
			description:  "Create comment up vote (original null, set to up)",
			method:       "POST",
			url:          "/votes/comment/1000/up",
			body:         "{}",
			expectedCode: 200,
			ip:           "192.168.1.1",
			expectedBody: func(t *testing.T, body string) {
				assert.NotEmpty(t, body)
				assert.Equal(t, `{"up":5,"down":2,"is_up":true,"is_down":false}`, body)
			},
		},
		{
			description:  "Create comment down vote (original null, set to down)",
			method:       "POST",
			url:          "/votes/comment/1000/down",
			body:         "{}",
			expectedCode: 200,
			ip:           "192.168.1.2",
			expectedBody: func(t *testing.T, body string) {
				assert.NotEmpty(t, body)
				assert.Equal(t, `{"up":4,"down":3,"is_up":false,"is_down":true}`, body)
			},
		},
		{
			description:  "Create page up vote (original null, set to up)",
			method:       "POST",
			url:          "/votes/page/1001/up",
			body:         "{}",
			expectedCode: 200,
			ip:           "192.168.1.3",
			expectedBody: func(t *testing.T, body string) {
				assert.NotEmpty(t, body)
				assert.Equal(t, `{"up":2,"down":2,"is_up":true,"is_down":false}`, body)
			},
		},
		{
			description:  "Create page down vote (original null, set to down)",
			method:       "POST",
			url:          "/votes/page/1001/down",
			body:         "{}",
			expectedCode: 200,
			ip:           "192.168.1.4",
			expectedBody: func(t *testing.T, body string) {
				assert.NotEmpty(t, body)
				assert.Equal(t, `{"up":1,"down":3,"is_up":false,"is_down":true}`, body)
			},
		},
		{
			description:  "Un-vote comment comment (original up, revoke up)",
			method:       "POST",
			url:          "/votes/comment/1000/up",
			body:         "{}",
			expectedCode: 200,
			ip:           "127.0.0.1",
			expectedBody: func(t *testing.T, body string) {
				assert.NotEmpty(t, body)
				assert.Equal(t, `{"up":3,"down":2,"is_up":false,"is_down":false}`, body)
			},
		},
		{
			description:  "Un-vote page comment (original down, revoke down)",
			method:       "POST",
			url:          "/votes/page/1001/down",
			body:         "{}",
			expectedCode: 200,
			ip:           "127.0.0.1",
			expectedBody: func(t *testing.T, body string) {
				assert.NotEmpty(t, body)
				assert.Equal(t, `{"up":1,"down":1,"is_up":false,"is_down":false}`, body)
			},
		},
		{
			description:  "Opposite-vote comment comment (original up, set to down)",
			method:       "POST",
			url:          "/votes/comment/1000/down",
			body:         "{}",
			expectedCode: 200,
			ip:           "127.0.0.1",
			expectedBody: func(t *testing.T, body string) {
				assert.NotEmpty(t, body)
				assert.Equal(t, `{"up":3,"down":3,"is_up":false,"is_down":true}`, body)
			},
		},
		{
			description:  "Opposite-vote page comment (original down, set to up)",
			method:       "POST",
			url:          "/votes/page/1001/up",
			body:         "{}",
			expectedCode: 200,
			ip:           "127.0.0.1",
			expectedBody: func(t *testing.T, body string) {
				assert.NotEmpty(t, body)
				assert.Equal(t, `{"up":2,"down":1,"is_up":true,"is_down":false}`, body)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			app, fiber := NewApiTestApp()
			defer app.Cleanup()

			handler.VoteGet(app.App, fiber)
			handler.VoteCreate(app.App, fiber)

			req := httptest.NewRequest(tt.method, tt.url, bytes.NewReader([]byte(tt.body)))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Forwarded-For", tt.ip) // mock IP
			resp, err := fiber.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			tt.expectedBody(t, string(body))
		})
	}
}

func TestVoteTransactionRollback(t *testing.T) {
	app, fiber := NewApiTestApp()
	defer app.Cleanup()

	handler.VoteCreate(app.App, fiber)

	db := app.Dao().DB()
	var originalComment entity.Comment
	require.NoError(t, db.First(&originalComment, 1000).Error)

	const callbackName = "test:fail_vote_count_update"
	require.NoError(t, db.Callback().Update().Before("gorm:update").Register(callbackName, func(tx *gorm.DB) {
		if tx.Statement.Table == "atk_comments" {
			tx.AddError(errors.New("injected vote count update failure"))
		}
	}))
	defer db.Callback().Update().Remove(callbackName)

	req := httptest.NewRequest("POST", "/votes/comment/1000/up", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", "192.168.100.1")
	resp, err := fiber.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)

	var voteCount int64
	require.NoError(t, db.Model(&entity.Vote{}).Where("target_id = ? AND ip = ?", 1000, "192.168.100.1").Count(&voteCount).Error)
	assert.Zero(t, voteCount, "Vote insert should be rolled back")

	var comment entity.Comment
	require.NoError(t, db.First(&comment, 1000).Error)
	assert.Equal(t, originalComment.VoteUp, comment.VoteUp, "Comment counters should remain unchanged")
	assert.Equal(t, originalComment.VoteDown, comment.VoteDown, "Comment counters should remain unchanged")
}
