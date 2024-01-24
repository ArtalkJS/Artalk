package comments_get

import (
	"github.com/ArtalkJS/Artalk/internal/dao"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"gorm.io/gorm"
)

// Basic scope for all queries
//
// It will ignore pending comments for non-admin users
func CommonScope(user entity.User) func(*gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		// Ignore pending comments
		if !user.IsAdmin { // If not admin, ignore pending comments
			d.Scopes(NoPending(user.ID))
		}

		return d
	}
}

// Filter root comments
func OnlyRoot() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("rid = 0")
	}
}

func NoPending(allowUserID ...uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 白名单用户 ID
		if len(allowUserID) > 0 && allowUserID[0] != 0 {
			return db.Where("(user_id = ? AND is_pending = ?) OR is_pending = ?", allowUserID[0], true, false)
		}

		return db.Where("is_pending = ?", false)
	}
}

func NoPendingChecker(user entity.User) func(*entity.Comment) bool {
	return func(comment *entity.Comment) bool {
		// Show all comments to admin
		if user.IsAdmin {
			return true
		}

		// Allow self comments even if pending
		if user.ID != 0 && user.ID == comment.UserID {
			return true
		}

		// Prevent pending comments
		if comment.IsPending {
			return false
		}

		return true
	}
}

// Filter by search keywords
func SearchScope(dao *dao.Dao, keywords string) func(d *gorm.DB) *gorm.DB {
	var userIds []uint
	dao.DB().Model(&entity.User{}).Where(
		"LOWER(name) = LOWER(?) OR LOWER(email) = LOWER(?)", keywords, keywords,
	).Pluck("id", &userIds)

	return func(d *gorm.DB) *gorm.DB {
		return d.Where("user_id IN (?) OR content LIKE ? OR page_key = ? OR ip = ? OR ua = ?",
			userIds, "%"+keywords+"%", keywords, keywords, keywords)
	}
}
