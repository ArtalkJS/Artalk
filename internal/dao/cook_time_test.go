package dao_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCookedTimestampsAcrossTimezones(t *testing.T) {
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

			for _, response := range []any{d.CookComment(comment), d.CookPage(page)} {
				data, err := json.Marshal(response)
				require.NoError(t, err)
				var fields map[string]json.RawMessage
				require.NoError(t, json.Unmarshal(data, &fields))
				var unix int64
				var utc, legacy string
				require.NoError(t, json.Unmarshal(fields["date_unix"], &unix))
				require.NoError(t, json.Unmarshal(fields["date_utc"], &utc))
				require.NoError(t, json.Unmarshal(fields["date"], &legacy))
				require.Equal(t, instant.Unix(), unix)
				require.Equal(t, "2025-01-02T03:04:05Z", utc)
				require.Equal(t, instant.In(zone).Format(dao.CommonDateTimeFormat), legacy)
			}
		})
	}
}
