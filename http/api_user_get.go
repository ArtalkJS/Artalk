package http

import (
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsUserGet struct {
	Name  string `mapstructure:"name"`
	Email string `mapstructure:"email"`
}

type ResponseUserGet struct {
	IsAdmin bool `json:"is_admin"`
}

func ActionUserGet(c echo.Context) error {
	var p ParamsUserGet
	if isOK, resp := ParamsDecode(c, ParamsUserGet{}, &p); !isOK {
		return resp
	}

	var user model.User
	lib.DB.Where("name = ? AND email = ?", p.Name, p.Email).First(&user)

	if user.IsEmpty() || user.Name != p.Name || user.Email != p.Email {
		return RespData(c, Map{
			"user":     nil,
			"is_login": false,
		})
	}

	// loginned user check
	isLogin := false
	tUser := GetUserByReqToken(c)
	if !tUser.IsEmpty() && tUser.Name == p.Name && tUser.Email == p.Email {
		isLogin = true
		user = tUser
	}

	return RespData(c, Map{
		"user":     user.ToCooked(),
		"is_login": isLogin,
	})
}
