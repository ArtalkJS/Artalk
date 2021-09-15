package http

import (
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

func ActionLogin(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")
	if name == "" {
		name = c.QueryParam("name")
	}
	if email == "" {
		email = c.QueryParam("email")
	}
	if password == "" {
		password = c.QueryParam("password")
	}

	if name == "" || email == "" || password == "" {
		return RespError(c, "Incomplete parameters.")
	}

	// record action for limiting action
	RecordAction(c)

	var user model.User
	lib.DB.Where("name = ? AND email = ?", name, email).First(&user) // name = ? OR email = ?
	if user.IsEmpty() || user.Password != password {
		return RespError(c, "验证失败")
	}

	return RespData(c, Map{
		"token": LoginGetUserToken(user),
	})
}
