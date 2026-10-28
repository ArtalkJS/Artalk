package handler

import (
	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsNotifyList struct {
	Name  string `query:"name" validate:"optional"`  // The user name
	Email string `query:"email" validate:"optional"` // The user email
}

type ResponseNotifyList struct {
	Notifies []entity.CookedNotify `json:"notifies"`
	Count    int                   `json:"count"`
}

// @Id           GetNotifies
// @Summary      Get Notifies
// @Description  Get unread notifications owned by the authenticated user or anonymous comment credential
// @Security     ApiKeyAuth
// @Tags         Notify
// @Param        name   query  string  true  "The user name"
// @Param        email  query  string  true  "The user email"
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponseNotifyList
// @Failure      400  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /notifies  [get]
func NotifyList(app *core.App, router fiber.Router) {
	router.Get("/notifies", func(c *fiber.Ctx) error {
		var p ParamsNotifyList
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		notifies := common.UnreadNotifications(app, c)

		return common.RespData(c, ResponseNotifyList{
			Notifies: notifies,
			Count:    len(notifies),
		})
	})
}
