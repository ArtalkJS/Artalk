package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsUserInfo struct {
	Name  string `query:"name" json:"name" validate:"optional"`   // The username
	Email string `query:"email" json:"email" validate:"optional"` // The user email
}

type ResponseUserInfo struct {
	User          *entity.CookedUser    `json:"user"`
	IsLogin       bool                  `json:"is_login"`
	Notifies      []entity.CookedNotify `json:"notifies"`
	NotifiesCount int                   `json:"notifies_count"`
}

// @Id           GetUser
// @Summary      Get User Info
// @Description  Get user info to prepare for login or check current user status
// @Tags         Account
// @Security     ApiKeyAuth
// @Param        user  query  ParamsUserInfo  true  "The user to query"
// @Produce      json
// @Success      200  {object}  ResponseUserInfo
// @Failure      400  {object}  Map{msg=string}
// @Router       /user  [get]
func UserInfo(app *core.App, router fiber.Router) {
	router.Get("/user", func(c *fiber.Ctx) error {
		// login status
		user := common.GetUserByReq(app, c)
		isLogin := !user.IsEmpty()

		var p ParamsUserInfo
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
			return common.RespData(c, ResponseUserInfo{
				User:          nil,
				IsLogin:       isLogin,
				Notifies:      []entity.CookedNotify{},
				NotifiesCount: 0,
			})
		}

		// unread notifies
		unreadNotifies := app.Dao().CookAllNotifies(app.Dao().FindUnreadNotifies(user.ID))
		cockedUser := app.Dao().CookUser(&user)

		return common.RespData(c, ResponseUserInfo{
			User:          &cockedUser,
			IsLogin:       isLogin,
			Notifies:      unreadNotifies,
			NotifiesCount: len(unreadNotifies),
		})
	})
}
