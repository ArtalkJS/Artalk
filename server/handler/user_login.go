package handler

import (
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type ParamsLogin struct {
	Name     string `json:"name"`                         // The username
	Email    string `json:"email" validate:"required"`    // The user email
	Password string `json:"password" validate:"required"` // The user password
}

type ResponseLogin struct {
	Token string            `json:"token"`
	User  entity.CookedUser `json:"user"`
}

// @Summary      Get Access Token
// @Description  Login user by name or email
// @Tags         Account
// @Param        user  body  ParamsLogin  true  "The user login data"
// @Accept       json
// @Produce      json
// @Success      200  {object}  common.JSONResult{data=ResponseLogin}
// @Failure      400  {object}  common.JSONResult{data=object{need_name_select=[]string}}  "Multiple users with the same email address are matched"
// @Router       /user/access_token  [post]
func UserLogin(app *core.App, router fiber.Router) {
	router.Post("/user/access_token", func(c *fiber.Ctx) error {
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
			users := app.Dao().FindUsersByEmail(p.Email)
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
			user = app.Dao().FindUser(p.Name, p.Email) // name = ? AND email = ?
		}

		if user.IsEmpty() {
			return common.RespError(c, i18n.T("Login failed"))
		}

		// 密码验证
		const bcryptPrefix = "(bcrypt)"
		const md5Prefix = "(md5)"
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

		jwtToken := common.LoginGetUserToken(user, app.Conf().AppKey, app.Conf().LoginTimeout)

		return common.RespData(c, ResponseLogin{
			Token: jwtToken,
			User:  app.Dao().CookUser(&user),
		})
	})
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

type ParamsLoginStatus struct {
	Name  string `query:"name" json:"name"`   // The username
	Email string `query:"email" json:"email"` // The user email
}

type ResponseLoginStatus struct {
	IsAdmin bool `json:"is_admin"`
	IsLogin bool `json:"is_login"`
}

// @Summary      Get Login Status
// @Description  Get user login status by header Authorization
// @Tags         Account
// @Security     ApiKeyAuth
// @Param        user  query  ParamsLoginStatus  true  "The user to query"
// @Produce      json
// @Success      200  {object}  common.JSONResult{data=ResponseLoginStatus}
// @Router       /user/status  [get]
func UserLoginStatus(app *core.App, router fiber.Router) {
	router.Get("/user/status", func(c *fiber.Ctx) error {
		var p ParamsLoginStatus
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		isAdmin := false
		if p.Email != "" && p.Name != "" {
			isAdmin = app.Dao().IsAdminUserByNameEmail(p.Name, p.Email)
		}

		return common.RespData(c, ResponseLoginStatus{
			IsAdmin: isAdmin,
			IsLogin: common.CheckIsAdminReq(app, c),
		})
	})
}
