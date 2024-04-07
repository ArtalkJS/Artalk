package comments_get

import (
	"github.com/ArtalkJS/Artalk/internal/entity"
)

// User Scope Query Tag
type UserScopeType string

const (
	UserAll      UserScopeType = "all"
	UserMentions UserScopeType = "mentions"
	UserMine     UserScopeType = "mine"
	UserPending  UserScopeType = "pending"
)

type UserScopePayload struct {
	Type UserScopeType
}

type UserScopeOpts struct {
	User            entity.User
	GetUserComments func(userID uint) []uint
}

// User Scope (for message center)
func UserScopeQuery(payload UserScopePayload, opts UserScopeOpts) func(liteDB) liteDB {
	return func(q liteDB) liteDB {
		// If user not found, return empty query
		if opts.User.IsEmpty() {
			return q.Where("1 = 0")
		}

		// Get user all comment ids which someone could reply to him
		userCommentIDs := opts.GetUserComments(opts.User.ID)

		scopes := map[UserScopeType]func(liteDB) liteDB{
			UserAll: func(d liteDB) liteDB {
				return q.Where("user_id = ? OR rid IN (?)", opts.User.ID, userCommentIDs)
			},
			UserMentions: func(d liteDB) liteDB {
				return q.Where("user_id != ? AND rid IN (?)", opts.User.ID, userCommentIDs)
			},
			UserMine: func(d liteDB) liteDB {
				return q.Where("user_id = ?", opts.User.ID)
			},
			UserPending: func(d liteDB) liteDB {
				return q.Where("user_id = ? AND is_pending = ?", opts.User.ID, true)
			},
		}

		if scope, ok := scopes[payload.Type]; ok {
			q.Scopes(scope)
		}

		return q
	}
}
