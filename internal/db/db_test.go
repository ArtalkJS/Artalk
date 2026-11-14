package db

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDB(t *testing.T) {
	db, err := NewDB(config.DBConf{
		Type:        config.TypeSQLite,
		Dsn:         "file::memory:?cache=shared",
		TablePrefix: "atk_",
	})
	defer CloseDB(db)

	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, db)
}

func TestNewTestDB(t *testing.T) {
	db, err := NewTestDB()
	defer CloseDB(db)

	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, db)
}

func TestDatabaseConfigurationFailures(t *testing.T) {
	blocked := filepath.Join(t.TempDir(), "file")
	require.NoError(t, os.WriteFile(blocked, []byte("not a directory"), 0600))
	for _, conf := range []config.DBConf{
		{Type: "unknown"},
		{Type: config.TypeSQLite, Dsn: filepath.Join(blocked, "db.sqlite")},
		{Type: config.TypeMySql, Dsn: "invalid-dsn"},
		{Type: config.TypePostgreSQL, Dsn: "postgres://[invalid"},
		{Type: config.TypeMSSQL, Dsn: "sqlserver://[invalid"},
		{Type: config.TypeMySql, Dsn: "unused", ServerCaPath: "missing.pem"},
	} {
		_, err := NewDB(conf)
		require.Error(t, err, string(conf.Type))
	}
	_, err := OpenSQLite("", nil)
	require.ErrorContains(t, err, "db.file")
}

func TestPreparedStatementConfiguration(t *testing.T) {
	for _, enabled := range []bool{true, false} {
		storage, err := NewDB(config.DBConf{Type: config.TypeSQLite, Dsn: filepath.Join(t.TempDir(), "db.sqlite"), PrepareStmt: &enabled})
		require.NoError(t, err)
		require.Equal(t, enabled, storage.PrepareStmt)
		var value int
		require.NoError(t, storage.Raw("SELECT ?", 42).Scan(&value).Error)
		require.Equal(t, 42, value)
		require.NoError(t, CloseDB(storage))
	}
}
