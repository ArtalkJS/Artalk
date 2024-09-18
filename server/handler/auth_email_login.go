package handler

import (
	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/i18n"
	"github.com/artalkjs/artalk/v2/internal/log"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/gofiber/fiber/v2"
)

type RequestAuthEmailLogin struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"optional"`
	Code     string `json:"code" validate:"optional"`
}

// @Id           LoginByEmail
// @Summary      Login by email
// @Description  Login by email with verify code (Need send email verify code first) or password
// @Tags         Auth
// @Param        data  body  RequestAuthEmailLogin  true  "The data to login"
// @Success      200  {object}  ResponseUserLogin
// @Failure      400  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Accept       json
// @Produce      json
// @Router       /auth/email/login  [post]
func AuthEmailLogin(app *core.App, router fiber.Router) {
	router.Post("/auth/email/login", common.LimiterGuard(app, func(c *fiber.Ctx) error {
		if !app.Conf().Auth.Email.Enabled {
			return common.RespError(c, 400, "Email auth is not enabled")
		}

		var p RequestAuthEmailLogin
		if ok, resp := common.ParamsDecode(c, &p); !ok {
			return resp
		}

		findUser := func(email string) (user entity.User) {
			users := app.Dao().FindUsersByEmail(email)
			if len(users) == 0 {
				return entity.User{}
			}
			// Select first user, if there are multiple users with same email
			return users[0]
		}

		var user entity.User
		if p.Code != "" {
			// Login by verify code
			if ok, resp := CheckEmailVerifyCode(app, c, p.Email, p.Code); !ok {
				return resp
			}

			user = findUser(p.Email)
		} else if p.Password != "" {
			// Login by password
			user = findUser(p.Email)

			// Check password
			if !user.CheckPassword(p.Password) {
				return common.RespError(c, 401, i18n.T("Password is incorrect"))
			}
		} else {
			return common.RespError(c, 400, "Password or code is required")
		}

		if user.IsEmpty() {
			return common.RespError(c, 401, "User not found")
		}

		// Get user token
		jwtToken, err := common.LoginGetUserToken(user, app.Conf().AppKey, app.Conf().LoginTimeout)
		if err != nil {
			log.Error("[LoginGetUserToken] ", err)
			return common.RespError(c, 500, i18n.T("Login failed"))
		}

		return common.RespData(c, ResponseUserLogin{
			Token: jwtToken,
			User:  app.Dao().CookUser(&user),
		})
	}))
}
