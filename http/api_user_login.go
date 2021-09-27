package http

import (
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsLogin struct {
	Name     string `mapstructure:"name" param:"required"`
	Email    string `mapstructure:"email" param:"required"`
	Password string `mapstructure:"password" param:"required"`
}

func ActionLogin(c echo.Context) error {
	var p ParamsLogin
	if isOK, resp := ParamsDecode(c, ParamsLogin{}, &p); !isOK {
		return resp
	}

	// record action for limiting action
	RecordAction(c)

	var user model.User
	lib.DB.Where(&model.User{
		Name:  p.Name,
		Email: p.Email,
	}).First(&user) // name = ? OR email = ?
	if user.IsEmpty() || user.Password != p.Password {
		return RespError(c, "验证失败")
	}

	return RespData(c, Map{
		"token": LoginGetUserToken(user),
	})
}
