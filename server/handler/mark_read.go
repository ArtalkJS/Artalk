package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsMarkRead struct {
	Name    string `json:"name"`     // The username
	Email   string `json:"email"`    // The user email
	AllRead bool   `json:"all_read"` // The option if mark all user's notify as read
}

// @Summary      Mark Notify as Read
// @Description  Mark specific notification as read for user
// @Tags         Notify
// @Param        comment_id  path  int             true  "The comment id of the notify you want to mark as read"
// @Param        notify_key  path  string          true  "The key of the notify"
// @Param        options     body  ParamsMarkRead  true  "The options"
// @Accept       json
// @Produce      json
// @Success      200  {object}  Map{}
// @Failure      400  {object}  Map{msg=string}
// @Failure      404  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /notifies/{comment_id}/{notify_key}/read  [post]
func MarkRead(app *core.App, router fiber.Router) {
	router.Post("/notifies/:comment_id/:notify_key", func(c *fiber.Ctx) error {
		commentID, _ := c.ParamsInt("comment_id")
		notifyKey := c.Params("notify_key")

		var p ParamsMarkRead
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// all read
		// TODO separate the API `/notifies/:comment_id/:notify_key/read` and `/notifies/read`
		if p.AllRead {
			if p.Name == "" || p.Email == "" {
				return common.RespError(c, 400, "username or email cannot be empty")
			}

			user := app.Dao().FindUser(p.Name, p.Email)
			err := app.Dao().UserNotifyMarkAllAsRead(user.ID)
			if err != nil {
				return common.RespError(c, 500, err.Error())
			}

			return common.RespSuccess(c)
		}

		// find notify
		notify := app.Dao().FindNotifyForComment(uint(commentID), notifyKey)
		if notify.IsEmpty() {
			return common.RespError(c, 404, i18n.T("{{name}} not found", Map{"name": i18n.T("Notify")}))
		}

		if notify.IsRead {
			return common.RespSuccess(c)
		}

		// update notify
		err := app.Dao().NotifySetRead(&notify)
		if err != nil {
			return common.RespError(c, 500, i18n.T("{{name}} save failed", Map{"name": i18n.T("Notify")}))
		}

		return common.RespSuccess(c)
	})
}
