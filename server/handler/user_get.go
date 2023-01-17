package handler

import (
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsUserGet struct {
	Name  string `form:"name"`
	Email string `form:"email"`
}

// POST /api/user-get
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
			return common.RespData(c, common.Map{
				"user":         nil,
				"is_login":     isLogin,
				"unread":       []interface{}{},
				"unread_count": 0,
			})
		}

		// unread notifies
		unreadNotifies := query.CookAllNotifies(query.FindUnreadNotifies(user.ID))

		return common.RespData(c, common.Map{
			"user":         query.CookUser(&user),
			"is_login":     isLogin,
			"unread":       unreadNotifies,
			"unread_count": len(unreadNotifies),
		})
	})
}
