package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/server/common"
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
// @Description  Get a list of notifies for user
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

		var user entity.User
		if p.Name != "" && p.Email != "" {
			user = app.Dao().FindUser(p.Name, p.Email)
		}

		var notifies = []entity.CookedNotify{}
		if !user.IsEmpty() {
			notifies = app.Dao().CookAllNotifies(app.Dao().FindUnreadNotifies(user.ID))
		}

		return common.RespData(c, ResponseNotifyList{
			Notifies: notifies,
			Count:    len(notifies),
		})
	})
}
