package comments_get

import (
	"testing"

	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/test"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestGetQueryScopes(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	normalUser := app.Dao().FindUserByID(1001)
	adminUser := app.Dao().FindUserByID(1000)

	tests := []struct {
		name string
		opts QueryOptions
		want func(comments []entity.Comment)
	}{
		{
			name: "ScopePage",
			opts: QueryOptions{
				User:  normalUser,
				Scope: ScopePage,
				PagePayload: PageScopePayload{
					SiteName: "Site A",
					PageKey:  "/test/1000.html",
				},
			},
			want: func(comments []entity.Comment) {
				assert.Greater(t, len(comments), 0)
			},
		},
		{
			name: "ScopeUser",
			opts: QueryOptions{
				User:  normalUser,
				Scope: ScopeUser,
				UserPayload: UserScopePayload{
					Type: UserAll,
				},
			},
			want: func(comments []entity.Comment) {
				assert.Greater(t, len(comments), 0)
			},
		},
		{
			name: "ScopeSite",
			opts: QueryOptions{
				User:  adminUser,
				Scope: ScopeSite,
				SitePayload: SitePayload{
					Type: SiteAll,
				},
			},
			want: func(comments []entity.Comment) {
				assert.Greater(t, len(comments), 0)
			},
		},
		{
			name: "Search",
			opts: QueryOptions{
				User:  adminUser,
				Scope: ScopeSite,
				SitePayload: SitePayload{
					Type: SiteAll,
				},
				Search: "Hello Artalk",
			},
			want: func(comments []entity.Comment) {
				assert.Greater(t, len(comments), 0)

				for _, c := range comments {
					assert.Contains(t, c.Content, "Hello Artalk")
				}
			},
		},
		{
			name: "Show Pending comments if admin",
			opts: QueryOptions{
				User:  adminUser,
				Scope: ScopePage,
				PagePayload: PageScopePayload{
					SiteName: "Site B",
					PageKey:  "/site_b/1001.html",
				},
			},
			want: func(comments []entity.Comment) {
				assert.Greater(t, len(comments), 0)
				assert.Greater(t, lo.CountBy(comments, func(c entity.Comment) bool {
					return c.IsPending
				}), 0)
			},
		},
		{
			name: "Hide Pending comments if not admin",
			opts: QueryOptions{
				User:  normalUser,
				Scope: ScopePage,
				PagePayload: PageScopePayload{
					SiteName: "Site B",
					PageKey:  "/site_b/1001.html",
				},
			},
			want: func(comments []entity.Comment) {
				assert.Greater(t, len(comments), 0)
				assert.Equal(t, 0, lo.CountBy(comments, func(c entity.Comment) bool {
					return c.IsPending
				}))
			},
		},
		{
			name: "Hide pending comments if empty user",
			opts: QueryOptions{
				User:  entity.User{},
				Scope: ScopePage,
				PagePayload: PageScopePayload{
					SiteName: "Site B",
					PageKey:  "/site_b/1001.html",
				},
			},
			want: func(comments []entity.Comment) {
				assert.Greater(t, len(comments), 0)
				assert.Equal(t, 0, lo.CountBy(comments, func(c entity.Comment) bool {
					return c.IsPending
				}))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scopes := GetQueryScopes(app.Dao(), tt.opts)
			var comments []entity.Comment
			app.Dao().DB().Scopes(ConvertGormScopes(scopes)...).Find(&comments)
			tt.want(comments)
		})
	}
}
