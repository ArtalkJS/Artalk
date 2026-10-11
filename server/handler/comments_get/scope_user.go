package comments_get

import (
	"github.com/artalkjs/artalk/v2/internal/entity"
	"gorm.io/gorm"
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
	User             entity.User
	GuestTokenHash   string
	GetGuestComments func() *gorm.DB
	GetUserComments  func(userID uint) []uint
}

// User Scope (for message center)
func UserScopeQuery(payload UserScopePayload, opts UserScopeOpts) func(liteDB) liteDB {
	return func(q liteDB) liteDB {
		if opts.User.IsEmpty() && opts.GuestTokenHash != "" && opts.GetGuestComments != nil {
			switch payload.Type {
			case UserMine:
				return q.Where("guest_token_hash = ?", opts.GuestTokenHash)
			case UserPending:
				return q.Where("guest_token_hash = ? AND is_pending = ?", opts.GuestTokenHash, true)
			case UserMentions:
				return q.Where("rid IN (?) AND guest_token_hash <> ?", opts.GetGuestComments(), opts.GuestTokenHash)
			default:
				return q.Where("guest_token_hash = ? OR rid IN (?)", opts.GuestTokenHash, opts.GetGuestComments())
			}
		}
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
