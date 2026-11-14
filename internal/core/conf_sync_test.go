package core

import (
	"path/filepath"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/internal/db"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/stretchr/testify/require"
)

func TestConfigSyncUpdatesExistingSiteAndAdministrator(t *testing.T) {
	storage, err := db.NewDB(config.DBConf{Type: config.TypeSQLite, Dsn: filepath.Join(t.TempDir(), "sync.db")})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.CloseDB(storage)) })
	cfg := &config.Config{SiteDefault: "Blog", SiteURL: "https://old.example.org", AdminUsers: []config.AdminUserConf{{Name: "Admin", Email: "admin@example.org", Password: "encrypted-password"}}}
	app := NewApp(cfg)
	app.SetDao(dao.NewDao(storage))
	app.syncFromConf()
	user := app.Dao().FindUser("Admin", "admin@example.org")
	require.True(t, user.IsAdmin)
	require.True(t, user.IsInConf)
	require.True(t, user.ReceiveEmail)
	site := app.Dao().FindSite("Blog")
	require.Equal(t, cfg.SiteURL, site.Urls)
	disabled := false
	cfg.SiteURL = "https://new.example.org"
	cfg.AdminUsers[0].ReceiveEmail = &disabled
	cfg.AdminUsers[0].Password = "updated-password-hash"
	cfg.AdminUsers[0].BadgeName = "Owner"
	cfg.AdminUsers[0].BadgeColor = "#123456"
	cfg.AdminUsers[0].Link = "https://owner.example.org"
	for range 2 {
		app.syncFromConf()
	}
	saved := app.Dao().FindUserByID(user.ID)
	require.Equal(t, cfg.AdminUsers[0].Password, saved.Password)
	require.Equal(t, "Owner", saved.BadgeName)
	require.Equal(t, cfg.AdminUsers[0].Link, saved.Link)
	require.False(t, saved.ReceiveEmail)
	require.Equal(t, cfg.SiteURL, app.Dao().FindSiteByID(site.ID).Urls)
	var users int64
	require.NoError(t, storage.Model(&entity.User{}).Count(&users).Error)
	require.EqualValues(t, 1, users)
}
