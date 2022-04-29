package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindComment(t *testing.T) {
	reloadTestDatabase()

	type args struct {
		id uint
	}
	type wants struct {
		id        uint
		rid       uint
		user_id   uint
		page_key  string
		site_name string
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{name: "评论 ID=1000", args: args{id: 1000}, wants: wants{id: 1000, rid: 0, user_id: 1000, page_key: "/test/1000.html", site_name: "Site A"}},
		{name: "评论 ID=1001", args: args{id: 1001}, wants: wants{id: 1001, rid: 1000, user_id: 1001, page_key: "/test/1000.html", site_name: "Site A"}},
		{name: "评论 ID=1006", args: args{id: 1006}, wants: wants{id: 1006, rid: 0, user_id: 1001, page_key: "/site_b/1001.html", site_name: "Site B"}},
		{name: "不存在的评论", args: args{id: 9999}, wants: wants{id: 0, rid: 0, user_id: 0, page_key: "", site_name: ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindComment(tt.args.id)
			assert.Equal(t, tt.wants.id, got.ID)
			assert.Equal(t, tt.wants.rid, got.Rid)
			assert.Equal(t, tt.wants.user_id, got.UserID)
			assert.Equal(t, tt.wants.page_key, got.PageKey)
			assert.Equal(t, tt.wants.site_name, got.SiteName)
			if tt.name != "不存在的评论" {
				assert.NotEmpty(t, got.Content)
			}
		})
	}
}

func TestFindCommentChildren(t *testing.T) {
	t.Run("Children Found", func(t *testing.T) {
		result := FindCommentChildren(1001)
		assert.Equal(t, 2, len(result))
	})

	t.Run("No Children Found", func(t *testing.T) {
		result := FindCommentChildren(1005)
		assert.Empty(t, result)
	})
}
