package http

import (
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ParamsAdminUserAdd struct {
	Name         string `mapstructure:"name" param:"required"`
	Email        string `mapstructure:"email" param:"required"`
	Password     string `mapstructure:"password"`
	Link         string `mapstructure:"link"`
	IsAdmin      bool   `mapstructure:"is_admin" param:"required"`
	SiteNames    string `mapstructure:"site_names"`
	ReceiveEmail bool   `mapstructure:"receive_email" param:"required"`
	BadgeName    string `mapstructure:"badge_name"`
	BadgeColor   string `mapstructure:"badge_color"`
}

func (a *action) AdminUserAdd(c echo.Context) error {
	if !GetIsSuperAdmin(c) {
		return RespError(c, "无权操作")
	}

	var p ParamsAdminUserAdd
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	if !model.FindUser(p.Name, p.Email).IsEmpty() {
		return RespError(c, "用户已存在")
	}

	if !lib.ValidateEmail(p.Email) {
		return RespError(c, "Invalid email")
	}
	if p.Link != "" && !lib.ValidateURL(p.Link) {
		return RespError(c, "Invalid link")
	}

	user := model.User{}
	user.Name = p.Name
	user.Email = p.Email
	user.Link = p.Link
	user.IsAdmin = p.IsAdmin
	user.SiteNames = p.SiteNames
	user.ReceiveEmail = p.ReceiveEmail
	user.BadgeName = p.BadgeName
	user.BadgeColor = p.BadgeColor

	if p.Password != "" {
		err := user.SetPasswordEncrypt(p.Password)
		if err != nil {
			logrus.Errorln(err)
			return RespError(c, "密码设置失败")
		}
	}

	err := model.CreateUser(&user)
	if err != nil {
		return RespError(c, "user 创建失败")
	}

	return RespData(c, Map{
		"user": user.ToCookedForAdmin(),
	})
}
