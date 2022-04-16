package http

import (
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
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

	user := model.FindUser(p.Name, p.Email) // name = ? OR email = ?

	if user.IsEmpty() {
		return RespError(c, "验证失败")
	}

	// 密码验证
	bcryptPrefix := "(bcrypt)"
	md5Prefix := "(md5)"
	passwordOK := false
	switch {
	case strings.HasPrefix(user.Password, bcryptPrefix):
		err := bcrypt.CompareHashAndPassword(
			[]byte(strings.TrimPrefix(user.Password, bcryptPrefix)),
			[]byte(p.Password),
		)

		if err == nil {
			passwordOK = true
		}
	case strings.HasPrefix(user.Password, md5Prefix):
		if strings.EqualFold(strings.TrimPrefix(user.Password, md5Prefix),
			fmt.Sprintf("%x", md5.Sum([]byte(p.Password)))) {
			passwordOK = true
		}
	default:
		if user.Password == p.Password {
			passwordOK = true
		}
	}

	if !passwordOK {
		return RespError(c, "验证失败")
	}

	return RespData(c, Map{
		"token": LoginGetUserToken(user),
	})
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// 获取当前登录状态
func ActionLoginStatus(c echo.Context) error {
	return RespData(c, Map{
		"is_login": CheckIsAdminReq(c),
	})
}
