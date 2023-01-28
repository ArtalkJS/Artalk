package handler

import (
	"crypto/md5"
	"fmt"
	"strings"
	"time"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type ParamsLogin struct {
	Name     string `form:"name"`
	Email    string `form:"email" validate:"required"`
	Password string `form:"password" validate:"required"`
}

type ResponseLogin struct {
	Token string            `json:"token"`
	User  entity.CookedUser `json:"user"`
}

// @Summary      User Login
// @Description  Login user by name or email
// @Tags         User
// @Param        name      formData  string  false  "the username"
// @Param        email     formData  string  true   "the user email"
// @Param        password  formData  string  true   "the user password"
// @Success      200  {object}  common.JSONResult{data=ResponseLogin}
// @Failure      400  {object}  common.JSONResult{data=object{need_name_select=[]string}}  "Multiple users with the same email address are matched"
// @Router       /login  [post]
func UserLogin(router fiber.Router) {
	router.Post("/login", func(c *fiber.Ctx) error {
		var p ParamsLogin
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// 账户读取
		var user entity.User
		if p.Name == "" {
			// 仅 Email 的查询
			if !utils.ValidateEmail(p.Email) {
				return common.RespError(c, i18n.T("Invalid {{name}}", Map{"name": i18n.T("Email")}))
			}
			users := query.FindUsersByEmail(p.Email)
			if len(users) == 1 {
				// 仅有一个 email 匹配的用户
				user = users[0]
			} else if len(users) > 1 {
				// 存在多个 email 匹配的用户
				userNames := []string{}
				for _, u := range users {
					userNames = append(userNames, u.Name)
				}
				return common.RespError(c, "Need to select username", common.Map{
					// 前端需做处理让用户选择用户名，
					// 之后再发起带 name 参数的请求
					"need_name_select": userNames,
				})
			}
		} else {
			// Name + Email 的精准查询
			user = query.FindUser(p.Name, p.Email) // name = ? AND email = ?
		}

		// record action for limiting action
		common.RecordAction(c)

		if user.IsEmpty() {
			return common.RespError(c, i18n.T("Login failed"))
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
			return common.RespError(c, i18n.T("Login failed"))
		}

		jwtToken := common.LoginGetUserToken(user)
		setAuthCookie(c, jwtToken, time.Now().Add(time.Second*time.Duration(config.Instance.LoginTimeout)))

		return common.RespData(c, ResponseLogin{
			Token: jwtToken,
			User:  query.CookUser(&user),
		})
	})
}

func setAuthCookie(c *fiber.Ctx, jwtToken string, expires time.Time) {
	if !config.Instance.Cookie.Enabled {
		return
	}

	// save jwt to cookie
	cookie := new(fiber.Cookie)
	cookie.Name = config.COOKIE_KEY_ATK_AUTH
	cookie.Value = jwtToken
	cookie.Expires = expires

	// @see https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Cookies
	// @see https://owasp.org/www-project-web-security-testing-guide/v41/4-Web_Application_Security_Testing/06-Session_Management_Testing/02-Testing_for_Cookies_Attributes
	cookie.Path = "/"
	cookie.HTTPOnly = true // prevent XSS
	cookie.Secure = true   // https only
	cookie.SameSite = ""   // for cors-request

	// @note cookie secure is not working on localhost
	// @see https://bugs.chromium.org/p/chromium/issues/detail?id=1177877#c7

	c.Cookie(cookie)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

type ParamsLoginStatus struct {
	Name  string `form:"name"`
	Email string `form:"email"`
}

type ResponseLoginStatus struct {
	IsAdmin bool `json:"is_admin"`
	IsLogin bool `json:"is_login"`
}

// @Summary      User Login Status
// @Description  Get user login status by header Authorization
// @Tags         User
// @Param        name           formData  string  false  "the username"
// @Param        email          formData  string  false  "the user email"
// @Param        password       formData  string  true   "the user password"
// @Security     ApiKeyAuth
// @Success      200  {object}  common.JSONResult{data=ResponseLoginStatus}
// @Router       /login-status  [post]
func UserLoginStatus(router fiber.Router) {
	router.Post("/login-status", func(c *fiber.Ctx) error {
		var p ParamsLoginStatus
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		isAdmin := false
		if p.Email != "" && p.Name != "" {
			isAdmin = query.IsAdminUserByNameEmail(p.Name, p.Email)
		}

		return common.RespData(c, ResponseLoginStatus{
			IsAdmin: isAdmin,
			IsLogin: common.CheckIsAdminReq(c),
		})
	})
}
