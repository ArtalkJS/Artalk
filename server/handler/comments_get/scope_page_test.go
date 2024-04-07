// Please note that this test case is dependent on
// the test dataset `comments.yml` in test pkg.
// Modify `comments.yml` may lead to test failure.
package comments_get

import (
	"testing"

	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/test"
	"github.com/stretchr/testify/assert"
)

func TestPageScopeQuery(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	tests := []struct {
		name    string
		payload PageScopePayload
		opts    PageScopeOpts
		want    func(comments []entity.Comment)
	}{
		{
			name: "Normal Page",
			payload: PageScopePayload{
				SiteName: "Site A",
				PageKey:  "/test/1000.html",
			},
			want: func(comments []entity.Comment) {
				assert.Greater(t, len(comments), 0)

				siteNames := make(map[string]bool)
				pageKeys := make(map[string]bool)
				for _, c := range comments {
					siteNames[c.SiteName] = true
					pageKeys[c.PageKey] = true
				}

				assert.Equal(t, 1, len(siteNames))
				assert.Equal(t, 1, len(pageKeys))

				assert.Equal(t, "Site A", comments[0].SiteName)
				assert.Equal(t, "/test/1000.html", comments[0].PageKey)
			},
		},
		{
			name: "AdminOnly",
			payload: PageScopePayload{
				SiteName: "Site A",
				PageKey:  "/test/1000.html",
				Tags:     []PageScopeTag{AdminOnly},
			},
			opts: PageScopeOpts{
				AdminUserIDs: []uint{1002},
			},
			want: func(comments []entity.Comment) {
				assert.Greater(t, len(comments), 0)

				userIDs := make(map[uint]bool)
				for _, c := range comments {
					userIDs[c.UserID] = true
				}

				assert.Equal(t, 1, len(userIDs))
				assert.Contains(t, userIDs, uint(1002))
			},
		},
		{
			name:    "SiteName is empty",
			payload: PageScopePayload{},
			want: func(comments []entity.Comment) {
				assert.Empty(t, comments)
			},
		},
		{
			name:    "PageKey is empty and SiteName is not empty",
			payload: PageScopePayload{SiteName: "Site B"},
			want: func(comments []entity.Comment) {
				assert.Empty(t, comments)
			},
		},
		{
			name:    "SiteName and PageKey are both empty",
			payload: PageScopePayload{},
			want: func(comments []entity.Comment) {
				assert.Empty(t, comments)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := PageScopeQuery(tt.payload, tt.opts)
			var comments []entity.Comment
			app.Dao().DB().Scopes(ConvertGormScopes(scope)...).Find(&comments)
			tt.want(comments)
		})
	}
}
