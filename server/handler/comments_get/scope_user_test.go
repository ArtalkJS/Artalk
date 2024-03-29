// Please note that this test case is dependent on
// the test dataset `comments.yml` in test pkg.
// Modify `comments.yml` may lead to test failure.
package comments_get

import (
	"testing"

	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/test"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
)

func TestUserScopeQuery(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	tests := []struct {
		name    string
		payload UserScopePayload
		opts    UserScopeOpts
		want    func(comments []entity.Comment)
	}{
		{
			name: "Type=UserAll",
			payload: UserScopePayload{
				Type: UserAll,
			},
			opts: UserScopeOpts{
				User: app.Dao().FindUserByID(1001),
				GetUserComments: func(userID uint) []uint {
					return app.Dao().GetUserAllCommentIDs(userID)
				},
			},
			want: func(comments []entity.Comment) {
				assert.Greater(t, len(comments), 0)

				m := make(map[uint]bool)

				for _, c := range comments {
					if c.Rid == 0 {
						// Root comment
						m[c.ID] = true
						assert.Equal(t, uint(1001), c.UserID, "root comment user_id should be user self")
					} else {
						// Reply comment
						if c.UserID == 1001 {
							m[c.ID] = true
						} else {
							m[c.Rid] = true
						}
					}
				}

				userCommentIDs := app.Dao().GetUserAllCommentIDs(1001)
				assert.ElementsMatch(t, userCommentIDs, maps.Keys(m))
			},
		},
		{
			name: "Type=UserMentions",
			payload: UserScopePayload{
				Type: UserMentions,
			},
			opts: UserScopeOpts{
				User: app.Dao().FindUserByID(1001),
				GetUserComments: func(userID uint) []uint {
					return app.Dao().GetUserAllCommentIDs(userID)
				},
			},
			want: func(comments []entity.Comment) {
				assert.Greater(t, len(comments), 0)

				m := make(map[uint]bool)

				for _, c := range comments {
					assert.NotEqual(t, uint(1001), c.UserID, "mentioned comment user_id should not be user self")
					m[c.Rid] = true
				}

				userCommentIDs := app.Dao().GetUserAllCommentIDs(1001)
				assert.Subset(t, userCommentIDs, maps.Keys(m), "mentioned comment rid should be in user's comment list")
			},
		},
		{
			name: "Type=UserMine",
			payload: UserScopePayload{
				Type: UserMine,
			},
			opts: UserScopeOpts{
				User: app.Dao().FindUserByID(1001),
				GetUserComments: func(userID uint) []uint {
					return app.Dao().GetUserAllCommentIDs(userID)
				},
			},
			want: func(comments []entity.Comment) {
				assert.Greater(t, len(comments), 0)

				for _, c := range comments {
					assert.Equal(t, uint(1001), c.UserID, "user_id should be user self")
				}
			},
		},
		{
			name: "Type=Pending",
			payload: UserScopePayload{
				Type: UserPending,
			},
			opts: UserScopeOpts{
				User: app.Dao().FindUserByID(1002),
				GetUserComments: func(userID uint) []uint {
					return app.Dao().GetUserAllCommentIDs(userID)
				},
			},
			want: func(comments []entity.Comment) {
				assert.Greater(t, len(comments), 0)

				for _, c := range comments {
					assert.Equal(t, uint(1002), c.UserID, "user_id should be user self")
					assert.True(t, c.IsPending, "should be pending comment")
				}
			},
		},
		{
			name: "User not found",
			payload: UserScopePayload{
				Type: UserAll,
			},
			opts: UserScopeOpts{
				User: entity.User{},
				GetUserComments: func(userID uint) []uint {
					return app.Dao().GetUserAllCommentIDs(userID)
				},
			},
			want: func(comments []entity.Comment) {
				assert.Empty(t, comments)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := UserScopeQuery(tt.payload, tt.opts)
			var comments []entity.Comment
			app.Dao().DB().Scopes(ConvertGormScopes(scope)...).Find(&comments)
			tt.want(comments)
		})
	}
}
