package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsGetNotifies struct {
	Name  string `query:"name" validate:"required"`  // The user name
	Email string `query:"email" validate:"required"` // The user email
}

type ResponseGetNotifies struct {
	Notifies []entity.CookedNotify `json:"notifies"`
	Count    int                   `json:"count"`
}

// @Summary      Get Notifies
// @Description  Get a list of notifies for user
// @Tags         Notify
// @Param        name   query  string  true  "The user name"
// @Param        email  query  string  true  "The user email"
// @Produce      json
// @Success      200  {object}  ResponseGetNotifies
// @Failure      400  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /notifies  [get]
func GetNotifies(app *core.App, router fiber.Router) {
	router.Get("/notifies", func(c *fiber.Ctx) error {
		var p ParamsGetNotifies
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		user := app.Dao().FindUser(p.Name, p.Email)

		var notifies = []entity.CookedNotify{}
		if !user.IsEmpty() {
			notifies = app.Dao().CookAllNotifies(app.Dao().FindUnreadNotifies(user.ID))
		}

		return common.RespData(c, ResponseGetNotifies{
			Notifies: notifies,
			Count:    len(notifies),
		})
	})
}

// Mark all notifies as read
type ParamsMarkAllRead struct {
	Name  string `json:"name" validate:"required"`  // The username
	Email string `json:"email" validate:"required"` // The user email
}

// @Summary      Mark All Notifies as Read
// @Description  Mark all notifies as read for user
// @Tags         Notify
// @Param        name   query  string  true  "The user name"
// @Param        email  query  string  true  "The user email"
// @Produce      json
// @Success      200  {object}  Map{}
// @Failure      400  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /notifies/read  [post]
func MarkAllRead(app *core.App, router fiber.Router) {
	router.Post("/notifies/read", func(c *fiber.Ctx) error {
		var p ParamsMarkAllRead
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

type ParamsMarkRead struct {
	Name  string `json:"name"`  // The username
	Email string `json:"email"` // The user email
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
