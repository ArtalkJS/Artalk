package dao_test

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/internal/db"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestSecurityMigrationStopsOnDatabaseFailure(t *testing.T) {
	if os.Getenv("ARTALK_TEST_MIGRATION_FAILURE") == "1" {
		database, err := db.NewTestDB()
		require.NoError(t, err)
		require.NoError(t, db.CloseDB(database))
		// A server must not continue when it cannot install the security schema.
		dao.NewDao(database)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestSecurityMigrationStopsOnDatabaseFailure$")
	cmd.Env = append(os.Environ(), "ARTALK_TEST_MIGRATION_FAILURE=1")
	output, err := cmd.CombinedOutput()
	var exit *exec.ExitError
	require.ErrorAs(t, err, &exit, string(output))
	require.Equal(t, 1, exit.ExitCode())
	require.Contains(t, string(output), "Failed to migrate schema")
}

// Model the old schema without either new security field.
type LegacySecurityComment struct {
	gorm.Model
	Content   string
	PageKey   string
	SiteName  string
	UserID    uint
	RootID    uint
	IsPending bool
}
type legacySecurityChallenge struct {
	gorm.Model
	Email     string
	Code      string
	ExpiresAt time.Time
	IP        string
	UA        string
}

func TestSecurityFieldsUpgradeLegacyDatabase(t *testing.T) {
	database, err := db.NewDB(config.DBConf{Type: config.TypeSQLite, Dsn: filepath.Join(t.TempDir(), "legacy.db"), TablePrefix: "atk_"})
	require.NoError(t, err)
	defer db.CloseDB(database)
	require.NoError(t, database.Table("atk_comments").AutoMigrate(&LegacySecurityComment{}))
	require.NoError(t, database.Table("atk_user_email_verifies").AutoMigrate(&legacySecurityChallenge{}))
	oldComment := LegacySecurityComment{Content: "legacy pending", PageKey: "/legacy", SiteName: "Legacy", UserID: 1, IsPending: true}
	oldCode := legacySecurityChallenge{Email: "guest@example.com", Code: "654321", ExpiresAt: time.Now().Add(time.Minute)}
	require.NoError(t, database.Table("atk_comments").Create(&oldComment).Error)
	require.NoError(t, database.Table("atk_user_email_verifies").Create(&oldCode).Error)
	require.False(t, database.Migrator().HasColumn("atk_comments", "guest_token_hash"))
	require.False(t, database.Migrator().HasColumn("atk_user_email_verifies", "attempts"))
	d := dao.NewDao(database)
	var comment entity.Comment
	var challenge entity.UserEmailVerify
	require.NoError(t, database.First(&comment, oldComment.ID).Error)
	require.NoError(t, database.First(&challenge, oldCode.ID).Error)
	require.Empty(t, comment.GuestTokenHash)
	require.Equal(t, oldComment.Content, comment.Content)
	require.True(t, comment.IsPending)
	require.Zero(t, challenge.Attempts)
	require.Equal(t, oldCode.Code, challenge.Code)
	// Repeated upgrades preserve proof and attempt budgets; old writes can omit both fields.
	hash := strings.Repeat("ab", 32)
	require.NoError(t, database.Model(&comment).Update("guest_token_hash", hash).Error)
	require.NoError(t, database.Model(&challenge).Update("attempts", 4).Error)
	d.MigrateModels()
	require.NoError(t, database.First(&comment, comment.ID).Error)
	require.NoError(t, database.First(&challenge, challenge.ID).Error)
	require.Equal(t, hash, comment.GuestTokenHash)
	require.Equal(t, 4, challenge.Attempts)
	oldComment.ID = 0
	require.NoError(t, database.Table("atk_comments").Create(&oldComment).Error)
	var inserted entity.Comment
	require.NoError(t, database.First(&inserted, oldComment.ID).Error)
	require.Empty(t, inserted.GuestTokenHash)
}

func TestGuestProofCacheAndPublicJSONCompatibility(t *testing.T) {
	c := newTestDaoCache(t)
	defer c.Close()
	comment := entity.Comment{Model: gorm.Model{ID: 123}, Content: "pending", IsPending: true, GuestTokenHash: strings.Repeat("ab", 32)}
	require.NoError(t, c.CommentCacheSave(&comment))
	var cached entity.Comment
	require.NoError(t, c.FindCache(fmt.Sprintf(dao.CommentByIDKey, 123), &cached))
	require.Equal(t, comment.GuestTokenHash, cached.GuestTokenHash)
	data, err := json.Marshal(cached)
	require.NoError(t, err)
	require.NotContains(t, string(data), comment.GuestTokenHash)
	// Pre-upgrade cache entries have no proof and must not grant guest access.
	require.NoError(t, c.StoreCache(LegacySecurityComment{Model: gorm.Model{ID: 124}, Content: "old", IsPending: true}, fmt.Sprintf(dao.CommentByIDKey, 124)))
	var legacy entity.Comment
	require.NoError(t, c.FindCache(fmt.Sprintf(dao.CommentByIDKey, 124), &legacy))
	require.Empty(t, legacy.GuestTokenHash)
	require.True(t, legacy.IsPending)
}

// This name exists only to model an earlier unreleased database/cache layout.
type legacyGuestProofComment struct {
	LegacySecurityComment
	GuestKeyHash string `gorm:"index;size:64;not null;default:''" json:"-"`
}

func TestGuestTokenRenamePreservesOwnership(t *testing.T) {
	database, err := db.NewDB(config.DBConf{Type: config.TypeSQLite, Dsn: filepath.Join(t.TempDir(), "draft.db"), TablePrefix: "atk_"})
	require.NoError(t, err)
	defer db.CloseDB(database)
	require.NoError(t, database.Table("atk_comments").AutoMigrate(&legacyGuestProofComment{}))
	old := legacyGuestProofComment{
		LegacySecurityComment: LegacySecurityComment{Content: "pending from draft build", IsPending: true},
		GuestKeyHash:          strings.Repeat("ab", 32),
	}
	require.NoError(t, database.Table("atk_comments").Create(&old).Error)
	c := newTestDaoCache(t)
	defer c.Close()
	require.NoError(t, c.StoreCache(old, fmt.Sprintf("comment#id=%d", old.ID)))
	d := dao.NewDao(database)
	d.SetCache(c)
	require.False(t, database.Migrator().HasColumn(&entity.Comment{}, "guest_key_hash"))
	require.True(t, database.Migrator().HasColumn(&entity.Comment{}, "guest_token_hash"))
	require.False(t, database.Migrator().HasIndex(&entity.Comment{}, "idx_atk_comments_guest_key_hash"))
	require.True(t, database.Migrator().HasIndex(&entity.Comment{}, "idx_atk_comments_guest_token_hash"))
	// Both the first database read and a subsequent cache hit retain the proof.
	for range 2 {
		comment := d.FindComment(old.ID)
		require.Equal(t, old.GuestKeyHash, comment.GuestTokenHash)
		require.True(t, comment.IsPending)
	}
	d.MigrateModels()
	var stored entity.Comment
	require.NoError(t, database.First(&stored, old.ID).Error)
	require.Equal(t, old.GuestKeyHash, stored.GuestTokenHash)
}
