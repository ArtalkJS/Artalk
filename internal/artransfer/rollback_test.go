package artransfer

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/internal/db"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/stretchr/testify/require"
)

func TestSkippedAncestorsKeepImportedRepliesVisible(t *testing.T) {
	storage, err := db.NewDB(config.DBConf{Type: config.TypeSQLite, Dsn: filepath.Join(t.TempDir(), "ancestors.db")})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.CloseDB(storage)) })
	d := dao.NewDao(storage)
	rows := []entity.Artran{
		{ID: "missing-root", PageKey: "/post"},
		{ID: "parent", Rid: "missing-root", SiteName: "Blog", PageKey: "/post", Content: "Parent"},
		{ID: "child", Rid: "parent", SiteName: "Blog", PageKey: "/post", Content: "Child"},
		{ID: "root", SiteName: "Blog", PageKey: "/post", Content: "Root"},
		{ID: "missing-middle", Rid: "root", SiteName: "Blog"},
		{ID: "other-parent", Rid: "missing-middle", SiteName: "Blog", PageKey: "/post", Content: "Other parent"},
		{ID: "other-child", Rid: "other-parent", SiteName: "Blog", PageKey: "/post", Content: "Other child"},
	}
	data, err := json.Marshal(rows)
	require.NoError(t, err)
	require.NoError(t, RunImportArtrans(d, &ImportParams{JsonData: string(data), Assumeyes: true}, func(string) {}))
	var saved []entity.Comment
	require.NoError(t, storage.Order("id").Find(&saved).Error)
	require.Len(t, saved, 5)
	for _, pair := range [][2]int{{0, 1}, {3, 4}} {
		parent, child := saved[pair[0]], saved[pair[1]]
		require.Zero(t, parent.Rid)
		require.Zero(t, parent.RootID)
		require.Equal(t, parent.ID, child.Rid)
		require.Equal(t, parent.ID, child.RootID)
	}
	require.Zero(t, saved[2].Rid)
}
