package http

import (
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsUserGet struct {
	Name  string `mapstructure:"name" param:"required"`
	Email string `mapstructure:"email" param:"required"`

	Site   string `mapstructure:"site"`
	SiteID uint
}

type ResponseUserGet struct {
	IsAdmin bool `json:"is_admin"`
}

func ActionUserGet(c echo.Context) error {
	var p ParamsUserGet
	if isOK, resp := ParamsDecode(c, ParamsUserGet{}, &p); !isOK {
		return resp
	}

	// find site
	p.SiteID = HandleSiteParam(p.Site)
	if isOK, resp := CheckSite(c, p.SiteID); !isOK {
		return resp
	}

	var user model.User // 注：user 查找是 AND
	lib.DB.Where(&model.User{
		Name:   p.Name,
		Email:  p.Email,
		SiteID: p.SiteID,
	}).First(&user)

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
