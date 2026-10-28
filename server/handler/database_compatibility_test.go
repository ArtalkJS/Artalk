package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/internal/db"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Some legacy query paths return data even when a preload fails. Count SQL
// errors as well as checking responses so compatibility failures cannot hide.
type compatibilitySQLLogger struct {
	logger.Interface
	failures atomic.Int32
}

func (l *compatibilitySQLLogger) Trace(ctx context.Context, begin time.Time, query func() (string, int64), err error) {
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		l.failures.Add(1)
	}
	l.Interface.Trace(ctx, begin, query, err)
}

// Minimal pre-upgrade tables, without the new security fields.
type compatibilityComment struct {
	gorm.Model
	Content        string
	PageKey        string `gorm:"size:255"`
	SiteName       string `gorm:"size:255"`
	UserID, RootID uint
	IsPending      bool
}

type compatibilityChallenge struct {
	gorm.Model
	Email     string `gorm:"index;size:255"`
	Code      string
	ExpiresAt time.Time
	IP, UA    string
}

// Opt in with ARTALK_TEST_DB_TYPE and ARTALK_TEST_DB_DSN against a disposable
// database. A unique table prefix isolates each run; cleanup touches only it.
func TestDatabaseCompatibility(t *testing.T) {
	kind, dsn := os.Getenv("ARTALK_TEST_DB_TYPE"), os.Getenv("ARTALK_TEST_DB_DSN")
	if kind == "" {
		t.Skip("external database not configured")
	}
	require.NotEmpty(t, dsn)
	prefix := fmt.Sprintf("compat_%x_", time.Now().UnixNano())
	database, err := db.NewDB(config.DBConf{Type: config.DBType(kind), Dsn: dsn, TablePrefix: prefix})
	require.NoError(t, err)
	sqlLog := &compatibilitySQLLogger{Interface: database.Logger}
	database.Logger = sqlLog
	t.Cleanup(func() {
		for _, model := range []any{&entity.Notify{}, &entity.Vote{}, &entity.Comment{}, &entity.UserEmailVerify{}, &entity.AuthIdentity{}, &entity.Page{}, &entity.Site{}, &entity.User{}} {
			require.NoError(t, database.Migrator().DropTable(model))
		}
		require.NoError(t, db.CloseDB(database))
	})
	t.Cleanup(func() { require.Zero(t, sqlLog.failures.Load(), "unexpected SQL errors") })
	commentTable, challengeTable := prefix+"comments", prefix+"user_email_verifies"
	require.NoError(t, database.Table(commentTable).AutoMigrate(&compatibilityComment{}))
	require.NoError(t, database.Table(challengeTable).AutoMigrate(&compatibilityChallenge{}))
	legacy := compatibilityComment{Content: "Legacy pending", IsPending: true}
	challenge := compatibilityChallenge{Email: "Guest@example.com", Code: "123456", ExpiresAt: time.Now().Add(5 * time.Minute)}
	require.NoError(t, database.Table(commentTable).Create(&legacy).Error)
	require.NoError(t, database.Table(challengeTable).Create(&challenge).Error)
	d := dao.NewDao(database)
	app := core.NewApp(&config.Config{AppKey: "database-compatibility-test", LoginTimeout: 3600})
	app.SetDao(d)
	router := fiber.New()
	handler.CommentGet(app, router)
	handler.CommentList(app, router)
	handler.NotifyList(app, router)
	handler.NotifyReadAll(app, router)
	handler.PageList(app, router)

	t.Run("migration", func(t *testing.T) {
		var comment entity.Comment
		var code entity.UserEmailVerify
		require.NoError(t, database.First(&comment, legacy.ID).Error)
		require.NoError(t, database.First(&code, challenge.ID).Error)
		require.Empty(t, comment.GuestTokenHash)
		require.True(t, comment.IsPending)
		require.Zero(t, code.Attempts)
		require.NoError(t, database.Model(&comment).Update("guest_token_hash", guestHash(guestA)).Error)
		require.NoError(t, database.Model(&code).Update("attempts", 3).Error)
		// Exercise the draft-column rename and repeat migration on every dialect.
		migrator := database.Migrator()
		index := database.NamingStrategy.IndexName(commentTable, "guest_token_hash")
		oldIndex := database.NamingStrategy.IndexName(commentTable, "guest_key_hash")
		require.NoError(t, migrator.RenameColumn(&entity.Comment{}, "guest_token_hash", "guest_key_hash"))
		require.NoError(t, migrator.RenameIndex(&entity.Comment{}, index, oldIndex))
		d.MigrateModels()
		d.MigrateModels()
		require.True(t, migrator.HasIndex(&entity.Comment{}, index))
		require.False(t, migrator.HasColumn(&entity.Comment{}, "guest_key_hash"))
		require.NoError(t, database.First(&comment, legacy.ID).Error)
		require.NoError(t, database.First(&code, challenge.ID).Error)
		require.Equal(t, guestHash(guestA), comment.GuestTokenHash)
		require.Equal(t, 3, code.Attempts)
		// Restore an unclaimed historical comment for the ownership checks below.
		require.NoError(t, database.Model(&comment).Update("guest_token_hash", "").Error)
	})

	site := d.NewSite("Compat", "https://example.test")
	user, err := d.NewUser("Visitor", "guest@example.com", "")
	require.NoError(t, err)
	admin, err := d.NewUser("Admin", "admin@example.com", "")
	require.NoError(t, err)
	require.NoError(t, database.Model(&admin).Update("is_admin", true).Error)
	adminToken, err := common.LoginGetUserToken(admin, app.Conf().AppKey, 60)
	require.NoError(t, err)
	page := d.NewPage("/post", "", site.Name)
	require.False(t, d.NewPage("/empty", "", site.Name).IsEmpty())
	parentA := entity.Comment{Content: "A", PageKey: page.Key, SiteName: site.Name, UserID: user.ID, IsPending: true, GuestTokenHash: guestHash(guestA)}
	parentB := parentA
	parentB.GuestTokenHash = guestHash(guestB)
	require.NoError(t, database.Create(&parentA).Error)
	require.NoError(t, database.Create(&parentB).Error)
	reply := entity.Comment{Content: "Reply", PageKey: page.Key, SiteName: site.Name, UserID: admin.ID, Rid: parentA.ID}
	require.NoError(t, database.Create(&reply).Error)
	owned := entity.Notify{UserID: user.ID, CommentID: reply.ID, Key: "owner"}
	adminCopy := entity.Notify{UserID: admin.ID, CommentID: reply.ID, Key: "admin"}
	require.NoError(t, database.Create(&owned).Error)
	require.NoError(t, database.Create(&adminCopy).Error)

	t.Run("guest ownership and notification updates", func(t *testing.T) {
		path := fmt.Sprintf("/comments/%d", parentA.ID)
		require.Equal(t, 200, guestRequest(t, router, "GET", path, guestA, "", nil).StatusCode)
		require.Equal(t, 404, guestRequest(t, router, "GET", path, guestB, "", nil).StatusCode)
		require.Equal(t, 404, guestRequest(t, router, "GET", fmt.Sprintf("/comments/%d", legacy.ID), guestA, "", nil).StatusCode)
		resp := guestRequest(t, router, "GET", "/comments?scope=user&type=mine&page_key=_&flat_mode=true&limit=20", guestA, "", nil)
		require.Equal(t, 200, resp.StatusCode)
		var comments handler.ResponseCommentList
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&comments))
		require.Len(t, comments.Comments, 1)
		require.Equal(t, parentA.ID, comments.Comments[0].ID)
		require.Equal(t, "https://example.test/post", comments.Comments[0].PageURL)
		for _, token := range []string{guestA, guestB} {
			resp := guestRequest(t, router, "GET", "/notifies", token, "", nil)
			require.Equal(t, 200, resp.StatusCode)
			var notices handler.ResponseNotifyList
			require.NoError(t, json.NewDecoder(resp.Body).Decode(&notices))
			if token == guestA {
				require.Equal(t, 1, notices.Count)
			} else {
				require.Zero(t, notices.Count)
			}
		}
		require.Equal(t, 200, guestRequest(t, router, "POST", "/notifies/read", guestA, "", nil).StatusCode)
		require.NoError(t, database.First(&owned, owned.ID).Error)
		require.NoError(t, database.First(&adminCopy, adminCopy.ID).Error)
		require.True(t, owned.IsRead)
		require.False(t, adminCopy.IsRead)
	})

	t.Run("page filtering and atomic PV", func(t *testing.T) {
		for _, filter := range []string{"", "?has_comments=true", "?has_comments=true&search=POST"} {
			resp := guestRequest(t, router, "GET", "/pages"+filter, "", adminToken, nil)
			require.Equal(t, 200, resp.StatusCode)
			var pages handler.ResponsePageList
			require.NoError(t, json.NewDecoder(resp.Body).Decode(&pages))
			if filter == "" {
				require.EqualValues(t, 2, pages.Total)
			} else {
				require.EqualValues(t, 1, pages.Total)
			}
		}
		require.NoError(t, database.Model(&page).Updates(map[string]any{"pv": nil, "admin_only": true, "title": "Admin title"}).Error)
		var callers sync.WaitGroup
		errs := make(chan error, 8)
		for range 8 {
			callers.Go(func() { snapshot := page; errs <- d.IncrementPagePV(&snapshot) })
		}
		callers.Wait()
		close(errs)
		for err := range errs {
			require.NoError(t, err)
		}
		require.NoError(t, database.First(&page, page.ID).Error)
		require.Equal(t, 8, page.PV)
		require.True(t, page.AdminOnly)
		require.Equal(t, "Admin title", page.Title)
	})

	t.Run("email challenge attempts and single use", func(t *testing.T) {
		router.Post("/verify", func(c *fiber.Ctx) error {
			ok, resp := handler.CheckEmailVerifyCode(app, c, "guest@example.com", c.Query("code"))
			if !ok {
				return resp
			}
			return c.SendStatus(200)
		})
		verify := func(code string) int {
			resp, err := router.Test(httptest.NewRequest(http.MethodPost, "/verify?code="+code, nil), 10000)
			if err != nil {
				return 0
			}
			defer resp.Body.Close()
			return resp.StatusCode
		}
		reset := func() {
			require.NoError(t, database.Model(&entity.UserEmailVerify{}).Where("id = ?", challenge.ID).
				Updates(map[string]any{"code": "123456", "attempts": 0}).Error)
		}
		reset()
		for range 5 {
			require.Equal(t, 400, verify("000000"))
		}
		require.Equal(t, 400, verify("123456"))
		reset()
		var callers sync.WaitGroup
		statuses := make(chan int, 8)
		for range 8 {
			callers.Go(func() { statuses <- verify("123456") })
		}
		callers.Wait()
		close(statuses)
		accepted := 0
		for status := range statuses {
			require.Contains(t, []int{200, 400}, status)
			if status == 200 {
				accepted++
			}
		}
		require.Equal(t, 1, accepted)
	})

	t.Run("title updates preserve edits and deletion", func(t *testing.T) {
		target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("<title>Fetched title</title>"))
		}))
		defer target.Close()
		page := d.NewPage(target.URL+"/post", "", site.Name)
		require.NoError(t, database.Model(&entity.Page{}).Where("id = ?", page.ID).
			Updates(map[string]any{"title": "Manual title", "pv": 70, "admin_only": true}).Error)
		require.NoError(t, d.TryFetchPageFromURL(&page, target.URL))
		require.Equal(t, "Manual title", d.FindPageByID(page.ID).Title)
		for range 2 {
			require.NoError(t, d.FetchPageFromURL(&page, target.URL))
			require.Equal(t, "Fetched title", page.Title)
			require.Equal(t, 70, page.PV)
			require.True(t, page.AdminOnly)
		}
		require.NoError(t, database.Delete(&entity.Page{}, page.ID).Error)
		require.Error(t, d.FetchPageFromURL(&page, target.URL))
		var deleted entity.Page
		require.NoError(t, database.Unscoped().First(&deleted, page.ID).Error)
		require.True(t, deleted.DeletedAt.Valid)
	})
}
