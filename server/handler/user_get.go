package handler

import (
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsUserGet struct {
	Name  string `form:"name"`
	Email string `form:"email"`
}

type ResponseUserGet struct {
	User        *entity.CookedUser    `json:"user"`
	IsLogin     bool                  `json:"is_login"`
	Unread      []entity.CookedNotify `json:"unread"`
	UnreadCount int                   `json:"unread_count"`
}

// @Summary      User Info Get
// @Description  Get user info to prepare for login or check current user status
// @Tags         User
// @Param        name           formData  string  false  "the username"
// @Param        email          formData  string  false  "the user email"
// @Security     ApiKeyAuth
// @Success      200  {object}  common.JSONResult{data=ResponseUserGet}
// @Router       /user-get  [post]
func UserGet(router fiber.Router) {
	router.Post("/user-get", func(c *fiber.Ctx) error {
		var p ParamsUserGet
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// login status
		isLogin := !common.GetUserByReq(c).IsEmpty()

		user := query.FindUser(p.Name, p.Email)
		if user.IsEmpty() {
			return common.RespData(c, ResponseUserGet{
				User:        nil,
				IsLogin:     isLogin,
				Unread:      []entity.CookedNotify{},
				UnreadCount: 0,
			})
		}

		// unread notifies
		unreadNotifies := query.CookAllNotifies(query.FindUnreadNotifies(user.ID))

		cockedUser := query.CookUser(&user)

		return common.RespData(c, ResponseUserGet{
			User:        &cockedUser,
			IsLogin:     isLogin,
			Unread:      unreadNotifies,
			UnreadCount: len(unreadNotifies),
		})
	})
}
