package http

import (
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ParamsAdminUserAdd struct {
	Name     string `mapstructure:"name" param:"required"`
	Email    string `mapstructure:"email" param:"required"`
	Link     string `mapstructure:"link"`
	Password string `mapstructure:"password"`
	IsAdmin  bool   `mapstructure:"is_admin"`
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

	user := model.User{}
	user.Name = p.Name
	user.Email = p.Email
	user.Link = p.Link
	user.IsAdmin = p.IsAdmin

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
