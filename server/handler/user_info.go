package handler

import (
	"errors"

	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/common"
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
// @Tags         Auth
// @Security     ApiKeyAuth
// @Param        user  query  ParamsUserInfo  true  "The user to query"
// @Produce      json
// @Success      200  {object}  ResponseUserInfo
// @Failure      400  {object}  Map{msg=string}
// @Router       /user  [get]
func UserInfo(app *core.App, router fiber.Router) {
	router.Get("/user", func(c *fiber.Ctx) error {
		var p ParamsUserInfo
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// login status
		user, err := common.GetUserByReq(app, c)
		isLogin := !user.IsEmpty()

		// Anonymous
		if errors.Is(err, common.ErrTokenNotProvided) && p.Name != "" && p.Email != "" {
			// If not login, find user by name and email
			user = app.Dao().FindUser(p.Name, p.Email)
		}

		unreadNotifies := common.UnreadNotifications(app, c)

		// If user not found
		if user.IsEmpty() {
			return common.RespData(c, ResponseUserInfo{
				User:          nil,
				IsLogin:       false,
				Notifies:      unreadNotifies,
				NotifiesCount: len(unreadNotifies),
			})
		}

		// Names and email addresses are not proof of identity. Keep only public
		// profile fields for anonymous autocomplete and the admin login prompt.
		cockedUser := entity.CookedUser{
			ID: user.ID, Name: user.Name, Link: user.Link,
			BadgeName: user.BadgeName, BadgeColor: user.BadgeColor, IsAdmin: user.IsAdmin,
		}
		if isLogin {
			cockedUser = app.Dao().CookUser(&user)
		}

		return common.RespData(c, ResponseUserInfo{
			User:          &cockedUser,
			IsLogin:       isLogin,
			Notifies:      unreadNotifies,
			NotifiesCount: len(unreadNotifies),
		})
	})
}
