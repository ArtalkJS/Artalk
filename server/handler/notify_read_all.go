package handler

import (
	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/gofiber/fiber/v2"
	"time"
)

// Mark all notifies as read
type ParamsNotifyReadAll struct {
	Name  string `json:"name" validate:"required"`  // The username
	Email string `json:"email" validate:"required"` // The user email
}

// @Id           MarkAllNotifyRead
// @Summary      Mark All Notifies as Read
// @Description  Mark owned notifications as read using a login or anonymous comment credential
// @Security     ApiKeyAuth
// @Tags         Notify
// @Param        options  body  ParamsNotifyReadAll  true  "The options"
// @Accept       json
// @Produce      json
// @Success      200  {object}  Map{}
// @Failure      401  {object}  Map{msg=string}
// @Failure      400  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /notifies/read  [post]
func NotifyReadAll(app *core.App, router fiber.Router) {
	router.Post("/notifies/read", func(c *fiber.Ctx) error {
		user, _ := common.GetUserByReq(app, c)
		if user.IsEmpty() && common.GuestTokenHash(c) == "" {
			return common.RespError(c, 401, "Notification credential required")
		}
		result := common.NotificationQuery(app, c).Where("is_read = ?", false).
			Updates(map[string]any{"is_read": true, "read_at": time.Now()})
		if result.Error != nil {
			return common.RespError(c, 500, "Failed to update notifications")
		}
		return common.RespSuccess(c)
	})
}
