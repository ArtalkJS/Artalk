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
	Name     string `mapstructure:"name"`
	Email    string `mapstructure:"email" param:"required"`
	Password string `mapstructure:"password" param:"required"`
}

func (a *action) Login(c echo.Context) error {
	var p ParamsLogin
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	// 账户读取
	var user model.User
	if p.Name == "" {
		// 仅 Email 的查询
		if !lib.ValidateEmail(p.Email) {
			return RespError(c, "请输入正确的邮箱")
		}
		users := model.FindUsersByEmail(p.Email)
		if len(users) == 1 {
			// 仅有一个 email 匹配的用户
			user = users[0]
		} else if len(users) > 1 {
			// 存在多个 email 匹配的用户
			userNames := []string{}
			for _, u := range users {
				userNames = append(userNames, u.Name)
			}
			return RespError(c, "需要选择用户名", Map{
				// 前端需做处理让用户选择用户名，
				// 之后再发起带 name 参数的请求
				"need_name_select": userNames,
			})
		}
	} else {
		// Name + Email 的精准查询
		user = model.FindUser(p.Name, p.Email) // name = ? AND email = ?
	}

	// record action for limiting action
	RecordAction(c)

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
		"user":  user.ToCooked(),
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
