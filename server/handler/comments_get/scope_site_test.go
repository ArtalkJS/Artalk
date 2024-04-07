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

func TestSiteScopeQuery(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	adminUser := entity.User{IsAdmin: true}
	normalUser := entity.User{IsAdmin: false}

	tests := []struct {
		name    string
		payload SitePayload
		user    entity.User
		want    func(comments []entity.Comment)
	}{
		{
			name: "Only Admin can access",
			payload: SitePayload{
				Type: SiteAll,
			},
			user: normalUser,
			want: func(comments []entity.Comment) {
				assert.Empty(t, comments)
			},
		},
		{
			name: "SiteName empty will return all sites",
			payload: SitePayload{
				Type: SiteAll,
			},
			user: adminUser,
			want: func(comments []entity.Comment) {
				assert.Greater(t, len(comments), 0)

				siteNames := make(map[string]bool)
				for _, c := range comments {
					siteNames[c.SiteName] = true
				}
				assert.Greater(t, len(siteNames), 1, "should have more than 1 site")
			},
		},
		{
			name: "Type=SiteAll and SiteName empty",
			payload: SitePayload{
				Type: SiteAll,
			},
			user: adminUser,
			want: func(comments []entity.Comment) {
				assert.Greater(t, len(comments), 0)

				// should have pending comments and non-pending comments
				hasPending := false
				hasNonPending := false
				for _, c := range comments {
					if c.IsPending {
						hasPending = true
					} else {
						hasNonPending = true
					}
				}
				assert.True(t, hasPending, "should have pending comments")
				assert.True(t, hasNonPending, "should have non-pending comments")
			},
		},
		{
			name: "Type=SitePending and SiteName empty",
			payload: SitePayload{
				Type: SitePending,
			},
			user: adminUser,
			want: func(comments []entity.Comment) {
				assert.Greater(t, len(comments), 0)

				hasNonPending := false
				for _, c := range comments {
					if !c.IsPending {
						hasNonPending = true
					}
				}
				assert.False(t, hasNonPending, "should only have pending comments")
			},
		},
		{
			name: "SiteName not empty",
			payload: SitePayload{
				Type:     SiteAll,
				SiteName: "Site B",
			},
			user: adminUser,
			want: func(comments []entity.Comment) {
				assert.Greater(t, len(comments), 0)

				siteNames := make(map[string]bool)
				for _, c := range comments {
					siteNames[c.SiteName] = true
				}
				assert.Equal(t, 1, len(siteNames), "should only have 1 site")
				assert.Equal(t, "Site B", comments[0].SiteName)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := SiteScopeQuery(tt.payload, tt.user)
			var comments []entity.Comment
			app.Dao().DB().Scopes(ConvertGormScopes(scope)...).Find(&comments)
			tt.want(comments)
		})
	}
}
