package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsUserLogin struct {
	Name     string `json:"name" validate:"optional"`     // The username
	Email    string `json:"email" validate:"required"`    // The user email
	Password string `json:"password" validate:"required"` // The user password
}

type ResponseUserLogin struct {
	Token string            `json:"token"`
	User  entity.CookedUser `json:"user"`
}

// @Id           Login
// @Summary      Get Access Token
// @Description  Login user by name or email
// @Tags         Auth
// @Param        user  body  ParamsUserLogin  true  "The user login data"
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponseUserLogin
// @Failure      400  {object}  Map{msg=string,data=object{need_name_select=[]string}}  "Multiple users with the same email address are matched"
// @Failure      401  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /user/access_token  [post]
func UserLogin(app *core.App, router fiber.Router) {
	router.Post("/user/access_token", common.LimiterGuard(app, func(c *fiber.Ctx) error {
		var p ParamsUserLogin
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// 账户读取
		var user entity.User
		if p.Name == "" {
			// 仅 Email 的查询
			if !utils.ValidateEmail(p.Email) {
				return common.RespError(c, 400, i18n.T("Invalid {{name}}", Map{"name": i18n.T("Email")}))
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
				return common.RespError(c, 400, "Need to select username", common.Map{
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
			return common.RespError(c, 401, i18n.T("Login failed"))
		}

		// 密码验证
		if !user.CheckPassword(p.Password) {
			return common.RespError(c, 401, i18n.T("Login failed"))
		}

		jwtToken, err := common.LoginGetUserToken(user, app.Conf().AppKey, app.Conf().LoginTimeout)
		if err != nil {
			return common.RespError(c, 500, err.Error())
		}

		return common.RespData(c, ResponseUserLogin{
			Token: jwtToken,
			User:  app.Dao().CookUser(&user),
		})
	}))
}
