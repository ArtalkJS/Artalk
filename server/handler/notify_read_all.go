package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

// Mark all notifies as read
type ParamsNotifyReadAll struct {
	Name  string `json:"name" validate:"required"`  // The username
	Email string `json:"email" validate:"required"` // The user email
}

// @Id           MarkAllNotifyRead
// @Summary      Mark All Notifies as Read
// @Description  Mark all notifies as read for user
// @Tags         Notify
// @Param        options  body  ParamsNotifyReadAll  true  "The options"
// @Accept       json
// @Produce      json
// @Success      200  {object}  Map{}
// @Failure      400  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /notifies/read  [post]
func NotifyReadAll(app *core.App, router fiber.Router) {
	router.Post("/notifies/read", func(c *fiber.Ctx) error {
		var p ParamsNotifyReadAll
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		user := app.Dao().FindUser(p.Name, p.Email)
		if user.IsEmpty() {
			return common.RespSuccess(c)
		}

		err := app.Dao().UserNotifyMarkAllAsRead(user.ID)
		if err != nil {
			return common.RespError(c, 500, err.Error())
		}

		return common.RespSuccess(c)
	})
}
