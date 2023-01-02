package http

import (
	"github.com/ArtalkJS/ArtalkGo/internal/query"
	"github.com/labstack/echo/v4"
)

type ParamsUserGet struct {
	Name  string `mapstructure:"name"`
	Email string `mapstructure:"email"`
}

func (a *action) UserGet(c echo.Context) error {
	var p ParamsUserGet
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	// login status
	isLogin := !GetUserByReq(c).IsEmpty()

	user := query.FindUser(p.Name, p.Email)
	if user.IsEmpty() {
		return RespData(c, Map{
			"user":         nil,
			"is_login":     isLogin,
			"unread":       []interface{}{},
			"unread_count": 0,
		})
	}

	// unread notifies
	unreadNotifies := query.CookAllNotifies(query.FindUnreadNotifies(user.ID))

	return RespData(c, Map{
		"user":         query.CookUser(&user),
		"is_login":     isLogin,
		"unread":       unreadNotifies,
		"unread_count": len(unreadNotifies),
	})
}
