package comments_get

import (
	"github.com/ArtalkJS/Artalk/internal/dao"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"gorm.io/gorm"
)

// Basic scope for all queries
//
// It will ignore pending comments for non-admin users
func CommonScope(user entity.User) func(liteDB) liteDB {
	return func(d liteDB) liteDB {
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

func NoPending(allowUserID ...uint) func(db liteDB) liteDB {
	return func(db liteDB) liteDB {
		// 白名单用户 ID
		if len(allowUserID) > 0 && allowUserID[0] != 0 {
			return db.Where("(user_id = ? AND is_pending = ?) OR is_pending = ?", allowUserID[0], true, false)
		}

		return db.Where("is_pending = ?", false)
	}
}

// Filter by search keywords
func SearchScope(dao *dao.Dao, keywords string) func(d liteDB) liteDB {
	var userIds []uint
	dao.DB().Model(&entity.User{}).Where(
		"LOWER(name) = LOWER(?) OR LOWER(email) = LOWER(?)", keywords, keywords,
	).Pluck("id", &userIds)

	return func(d liteDB) liteDB {
		return d.Where("user_id IN (?) OR content LIKE ? OR page_key = ? OR ip = ? OR ua = ?",
			userIds, "%"+keywords+"%", keywords, keywords, keywords)
	}
}
