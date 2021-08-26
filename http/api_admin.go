package http

import (
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

func ActionLogin(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")
	if username == "" || password == "" {
		return RespError(c, "Incomplete parameters.")
	}

	var user model.User
	lib.DB.Where("name = ? OR email = ?", username, username).First(user)
	if user.IsEmpty() || user.Password != password {
		return RespError(c, "验证失败")
	}

	return RespData(c, Map{
		"token": LoginGetUserToken(user),
	})
}

func ActionAdminEdit(c echo.Context) error {
	// user := c.Get("user").(*jwt.Token)
	// claims := user.Claims.(*jwtCustomClaims)
	// name := claims.Name
	return nil
}

func ActionAdminDel(c echo.Context) error {
	return nil
}
