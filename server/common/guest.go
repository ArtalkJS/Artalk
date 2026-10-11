package common

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const GuestTokenHeader = "X-Artalk-Guest"

// GuestTokenHash proves possession of a browser secret, never ownership of a user
// account. Only headers are accepted so credentials do not enter URLs or logs.
func GuestTokenHash(c *fiber.Ctx) string {
	if GetTokenByReq(c) != "" {
		return ""
	} // Invalid/expired login tokens cannot fall back to a guest.
	token := c.Get(GuestTokenHeader)
	if len(token) != 64 {
		return ""
	}
	raw, err := hex.DecodeString(token)
	if err != nil || len(raw) != 32 {
		return ""
	}
	hash := sha256.Sum256(raw)
	return hex.EncodeToString(hash[:])
}

// NotificationQuery applies the same ownership rule to listing and marking read.
// Guest access is limited to replies to comments bearing that browser's proof,
// with an exact recipient match (never an administrator's copy of a notification).
func NotificationQuery(app *core.App, c *fiber.Ctx) *gorm.DB {
	db := app.Dao().DB()
	query := db.Model(&entity.Notify{})
	user, _ := GetUserByReq(app, c)
	if !user.IsEmpty() {
		return query.Where("user_id = ?", user.ID)
	}
	hash := GuestTokenHash(c)
	if hash == "" {
		return query.Where("1 = 0")
	}
	notifyTable := db.NamingStrategy.TableName("Notify")
	ownedParentComments := db.Model(&entity.Comment{}).Select("id").Where("guest_token_hash = ? AND user_id = ?", hash,
		clause.Column{Table: notifyTable, Name: "user_id"})
	visibleRepliesToOwnedComments := db.Model(&entity.Comment{}).Select("1").Where("id = ? AND is_pending = ?",
		clause.Column{Table: notifyTable, Name: "comment_id"}, false).Where("rid IN (?)", ownedParentComments)
	return query.Where("EXISTS (?)", visibleRepliesToOwnedComments)
}

func UnreadNotifications(app *core.App, c *fiber.Ctx) []entity.CookedNotify {
	var notices []entity.Notify
	if err := NotificationQuery(app, c).Where("is_read = ?", false).Find(&notices).Error; err != nil {
		return []entity.CookedNotify{}
	}
	return app.Dao().CookAllNotifies(notices)
}
