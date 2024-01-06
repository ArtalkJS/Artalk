package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsUserGet struct {
	Name  string `query:"name" json:"name"`   // The username
	Email string `query:"email" json:"email"` // The user email
}

type ResponseUserGet struct {
	User          *entity.CookedUser    `json:"user"`
	IsLogin       bool                  `json:"is_login"`
	Notifies      []entity.CookedNotify `json:"notifies"`
	NotifiesCount int                   `json:"notifies_count"`
}

// @Summary      Get User Info
// @Description  Get user info to prepare for login or check current user status
// @Tags         Account
// @Security     ApiKeyAuth
// @Param        user  query  ParamsUserGet  true  "The user to query"
// @Produce      json
// @Success      200  {object}  ResponseUserGet
// @Failure      400  {object}  Map{msg=string}
// @Router       /user/info  [get]
func UserGet(app *core.App, router fiber.Router) {
	router.Get("/user/info", func(c *fiber.Ctx) error {
		// login status
		user := common.GetUserByReq(app, c)
		isLogin := !user.IsEmpty()

		var p ParamsUserGet
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		if !isLogin {
			user = app.Dao().FindUser(p.Name, p.Email)
		} else {
			if p.Name != "" || p.Email != "" {
				return common.RespError(c, 400, "Not necessary to query user info when logged in")
			}
		}

		if user.IsEmpty() {
			return common.RespData(c, ResponseUserGet{
				User:          nil,
				IsLogin:       isLogin,
				Notifies:      []entity.CookedNotify{},
				NotifiesCount: 0,
			})
		}

		// unread notifies
		unreadNotifies := app.Dao().CookAllNotifies(app.Dao().FindUnreadNotifies(user.ID))
		cockedUser := app.Dao().CookUser(&user)

		return common.RespData(c, ResponseUserGet{
			User:          &cockedUser,
			IsLogin:       isLogin,
			Notifies:      unreadNotifies,
			NotifiesCount: len(unreadNotifies),
		})
	})
}
