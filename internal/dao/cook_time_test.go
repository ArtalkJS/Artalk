package dao_test

import (
	"testing"
	"time"

	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCookedDateUnixAcrossTimezones(t *testing.T) {
	originalLocal := time.Local
	t.Cleanup(func() { time.Local = originalLocal })
	instant := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)
	d := &dao.Dao{}

	for _, zone := range []*time.Location{
		time.UTC,
		time.FixedZone("UTC+8", 8*60*60),
		time.FixedZone("UTC-5", -5*60*60),
	} {
		t.Run(zone.String(), func(t *testing.T) {
			time.Local = zone
			model := gorm.Model{CreatedAt: instant.In(zone)}
			page := &entity.Page{Model: model, Key: "https://example.test/page"}
			comment := &entity.Comment{Model: model, Page: page, User: &entity.User{}}
			legacy := instant.In(zone).Format(dao.CommonDateTimeFormat)

			cookedComment := d.CookComment(comment)
			assert.Equal(t, instant.Unix(), cookedComment.DateUnix)
			assert.Equal(t, legacy, cookedComment.Date)

			cookedPage := d.CookPage(page)
			assert.Equal(t, instant.Unix(), cookedPage.DateUnix)
			assert.Equal(t, legacy, cookedPage.Date)
		})
	}
}
