package http

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type ParamsLogin struct {
	Name     string `mapstructure:"name" param:"required"`
	Email    string `mapstructure:"email" param:"required"`
	Password string `mapstructure:"password" param:"required"`
}

func (a *action) Login(c echo.Context) error {
	var p ParamsLogin
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	// record action for limiting action
	RecordAction(c)

	user := model.FindUser(p.Name, p.Email) // name = ? AND email = ?
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

	jwtToken := LoginGetUserToken(user)
	setAuthCookie(c, jwtToken, time.Now().Add(time.Second*time.Duration(config.Instance.LoginTimeout)))

	return RespData(c, Map{
		"token": jwtToken,
	})
}

func setAuthCookie(c echo.Context, jwtToken string, expires time.Time) {
	if !config.Instance.Cookie.Enabled {
		return
	}

	// save jwt token to cookie
	cookie := new(http.Cookie)
	cookie.Name = lib.COOKIE_KEY_ATK_AUTH
	cookie.Value = jwtToken
	cookie.Expires = expires

	// @see https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Cookies
	// @see https://owasp.org/www-project-web-security-testing-guide/v41/4-Web_Application_Security_Testing/06-Session_Management_Testing/02-Testing_for_Cookies_Attributes
	cookie.Path = "/"
	cookie.HttpOnly = true                     // prevent XSS
	cookie.Secure = true                       // https only
	cookie.SameSite = http.SameSiteDefaultMode // for cors-request

	// @note cookie secure is not working on localhost
	// @see https://bugs.chromium.org/p/chromium/issues/detail?id=1177877#c7

	c.SetCookie(cookie)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

type ParamsLoginStatus struct {
	Name  string `mapstructure:"name"`
	Email string `mapstructure:"email"`
}

// 获取当前登录状态
func (a *action) LoginStatus(c echo.Context) error {
	var p ParamsLoginStatus
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	isAdmin := false
	if p.Email != "" && p.Name != "" {
		isAdmin = model.IsAdminUserByNameEmail(p.Name, p.Email)
	}

	return RespData(c, Map{
		"is_admin": isAdmin,
		"is_login": CheckIsAdminReq(c),
	})
}
