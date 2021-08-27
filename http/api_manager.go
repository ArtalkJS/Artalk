package http

import (
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

func ActionLogin(c echo.Context) error {
	username := c.QueryParam("user")
	password := c.QueryParam("password")
	if username == "" || password == "" {
		return RespError(c, "Incomplete parameters.")
	}

	// record action for limiting action
	RecordAction(c)

	var user model.User
	lib.DB.Where("name = ? OR email = ?", username, username).First(&user)
	if user.IsEmpty() || user.Password != password {
		return RespError(c, "验证失败")
	}

	return RespData(c, Map{
		"token": LoginGetUserToken(user),
	})
}

func ActionManagerEdit(c echo.Context) error {
	// user := c.Get("user").(*jwt.Token)
	// claims := user.Claims.(*jwtCustomClaims)
	// name := claims.Name
	return nil
}

func ActionManagerDel(c echo.Context) error {
	return nil
}

func ActionManagerSendMail(c echo.Context) error {
	return RespSuccess(c)
}
