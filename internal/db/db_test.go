package db

import (
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/stretchr/testify/assert"
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
