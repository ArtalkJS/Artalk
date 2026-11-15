package handler_test

import (
	"bytes"
	"errors"
	"github.com/artalkjs/artalk/v2/internal/cache"
	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/dao"
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

func TestVoteCreateReturnsNotFoundWhenTransactionTargetDisappears(t *testing.T) {
	tests := []struct {
		name         string
		target       string
		targetID     string
		table        string
		expectedBody string
	}{
		{
			name:         "comment",
			target:       "comment",
			targetID:     "1000",
			table:        "atk_comments",
			expectedBody: `{"msg":"Comment not found"}`,
		},
		{
			name:         "page",
			target:       "page",
			targetID:     "1001",
			table:        "atk_pages",
			expectedBody: `{"msg":"Page not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, fiber := NewApiTestApp()
			defer app.Cleanup()

			handler.VoteCreate(app.App, fiber)

			db := app.Dao().DB()
			const callbackName = "test:missing_vote_target_lock"
			injected := false
			require.NoError(t, db.Callback().Query().Before("gorm:query").Register(callbackName, func(tx *gorm.DB) {
				_, hasLock := tx.Statement.Clauses["FOR"]
				if tx.Statement.Table == tt.table && hasLock {
					injected = true
					tx.AddError(gorm.ErrRecordNotFound)
				}
			}))
			defer db.Callback().Query().Remove(callbackName)

			url := "/votes/" + tt.target + "/" + tt.targetID + "/up"
			req := httptest.NewRequest("POST", url, bytes.NewReader([]byte("{}")))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Forwarded-For", "192.168.100.2")
			resp, err := fiber.Test(req)
			require.NoError(t, err)
			assert.True(t, injected, "The missing target error should be injected while locking the target")
			assert.Equal(t, 404, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.JSONEq(t, tt.expectedBody, string(body))
		})
	}
}

func TestVoteStorageFailuresPreservePreviousChoice(t *testing.T) {
	for _, stage := range []string{"find", "delete", "create", "count_up", "count_down", "page_counter"} {
		t.Run(stage, func(t *testing.T) {
			app, router := NewApiTestApp()
			t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
			handler.VoteCreate(app.App, router)
			db := app.Dao().DB()
			var before []entity.Vote
			require.NoError(t, db.Order("id").Find(&before).Error)
			comment := app.Dao().FindComment(1000)
			page := app.Dao().FindPageByID(1001)
			const callback = "test:vote_storage_failure"
			fail := func(tx *gorm.DB) { tx.AddError(errors.New("vote storage unavailable")) }
			var remove func() error
			switch stage {
			case "delete":
				require.NoError(t, db.Callback().Delete().Before("gorm:delete").Register(callback, fail))
				remove = func() error { return db.Callback().Delete().Remove(callback) }
			case "create":
				require.NoError(t, db.Callback().Create().Before("gorm:create").Register(callback, fail))
				remove = func() error { return db.Callback().Create().Remove(callback) }
			case "page_counter":
				require.NoError(t, db.Callback().Update().Before("gorm:update").Register(callback, fail))
				remove = func() error { return db.Callback().Update().Remove(callback) }
			default:
				query := 0
				require.NoError(t, db.Callback().Query().Before("gorm:query").Register(callback, func(tx *gorm.DB) {
					if tx.Statement.Table == "atk_votes" {
						query++
						if query == map[string]int{"find": 1, "count_up": 2, "count_down": 3}[stage] {
							fail(tx)
						}
					}
				}))
				remove = func() error { return db.Callback().Query().Remove(callback) }
			}
			t.Cleanup(func() { require.NoError(t, remove()) })
			url := "/votes/comment/1000/down"
			if stage == "page_counter" {
				url = "/votes/page/1001/up"
			}
			response, err := securityRequest(router, "POST", url, map[string]string{}, "", "127.0.0.1")
			require.NoError(t, err)
			response.Body.Close()
			require.Equal(t, 500, response.StatusCode)
			require.NoError(t, remove())
			var after []entity.Vote
			require.NoError(t, db.Order("id").Find(&after).Error)
			require.Equal(t, before, after)
			require.Equal(t, comment, app.Dao().FindComment(comment.ID))
			require.Equal(t, page, app.Dao().FindPageByID(page.ID))
		})
	}
}

func TestVotingRefreshesCachedCounters(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	cached, err := cache.New(config.CacheConf{Type: config.CacheTypeBuiltin})
	require.NoError(t, err)
	t.Cleanup(cached.Close)
	app.Dao().SetCache(dao.NewCacheAdaptor(cached))
	handler.VoteCreate(app.App, router)
	comment := app.Dao().FindComment(1001)
	page := app.Dao().FindPageByID(1002)
	require.Zero(t, comment.VoteUp)
	require.Zero(t, page.VoteUp)
	flowRequest(t, router, "POST", "/votes/comment/1001/up", map[string]string{}, "", 200)
	flowRequest(t, router, "POST", "/votes/page/1002/up", map[string]string{}, "", 200)
	require.Equal(t, 1, app.Dao().FindComment(comment.ID).VoteUp)
	require.Equal(t, 1, app.Dao().FindPageByID(page.ID).VoteUp)
}

func TestVoteValidationAndVoterStorageFailure(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	handler.VoteCreate(app.App, router)
	voter := handler.ParamsVoteCreate{Name: "New voter", Email: "voter@example.org"}
	for _, url := range []string{
		"/votes/comment/1001/unknown", "/votes/unknown/1001/up",
		"/votes/comment/999999/up", "/votes/page/999999/up",
	} {
		flowRequest(t, router, "POST", url, voter, "", 404)
	}
	flowRequest(t, router, "POST", "/votes/comment/1001/up", map[string]any{"name": 42}, "", 400)
	require.True(t, app.Dao().FindUser(voter.Name, voter.Email).IsEmpty())
	const callback = "test:voter_storage_failure"
	require.NoError(t, app.Dao().DB().Callback().Create().Before("gorm:create").Register(callback, func(tx *gorm.DB) {
		if tx.Statement.Table == "atk_users" {
			tx.AddError(errors.New("user storage unavailable"))
		}
	}))
	t.Cleanup(func() { require.NoError(t, app.Dao().DB().Callback().Create().Remove(callback)) })
	flowRequest(t, router, "POST", "/votes/comment/1001/up", voter, "", 500)
	require.Zero(t, app.Dao().FindComment(1001).VoteUp)
	require.True(t, app.Dao().FindUser(voter.Name, voter.Email).IsEmpty())
	require.NoError(t, app.Dao().DB().Callback().Create().Remove(callback))
	flowRequest(t, router, "POST", "/votes/comment/1001/up", voter, "", 200)
	user := app.Dao().FindUser(voter.Name, voter.Email)
	require.NotZero(t, user.ID)
	var vote entity.Vote
	require.NoError(t, app.Dao().DB().Where("target_id = ? AND type = ?", 1001, entity.VoteTypeCommentUp).First(&vote).Error)
	require.Equal(t, user.ID, vote.UserID)
}
