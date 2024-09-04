package dao_test

import (
	"testing"

	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/test"
	"github.com/stretchr/testify/assert"
)

func TestGetTableName(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	assert.Equal(t, "atk_pages", app.Dao().GetTableName(&entity.Page{}))
	assert.Equal(t, "atk_comments", app.Dao().GetTableName(&entity.Comment{}))
}
