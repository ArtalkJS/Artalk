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

	if p.Name == "" && p.Email == "" {
		return RespError(c, "Please input name or email.")
	}

	var user model.User // 注：user 查找是 AND
	lib.DB.Where("name = ? AND email = ?", p.Name, p.Email).First(&user)

	isLogin := false
	tUser := GetUserByReqToken(c)
	if !tUser.IsEmpty() && tUser.Name == p.Name && tUser.Email == p.Email {
		isLogin = true
		user = tUser
	}

	if user.IsEmpty() {
		return RespData(c, Map{
			"user":     nil,
			"is_login": false,
		})
	}

	return RespData(c, Map{
		"user":     user.ToCooked(),
		"is_login": isLogin,
	})
}
