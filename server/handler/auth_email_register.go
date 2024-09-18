package handler

import (
	"strings"

	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/utils"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/gofiber/fiber/v2"
)

type RequestAuthEmailRegister struct {
	Code     string `json:"code" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Name     string `json:"name" validate:"optional"`
	Link     string `json:"link" validate:"optional"`
	Password string `json:"password" validate:"required"`
}

// @Id           RegisterByEmail
// @Summary      Register by email
// @Description  Register by email and verify code (if user exists, will update user, like forget or change password. Need send email verify code first)
// @Tags         Auth
// @Param        data  body  RequestAuthEmailRegister  true  "The data to register"
// @Success      200  {object}  ResponseUserLogin
// @Failure      400  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Accept       json
// @Produce      json
// @Router       /auth/email/register  [post]
func AuthEmailRegister(app *core.App, router fiber.Router) {
	router.Post("/auth/email/register", common.LimiterGuard(app, func(c *fiber.Ctx) error {
		if !app.Conf().Auth.Email.Enabled {
			return common.RespError(c, 400, "Email auth is not enabled")
		}

		var p RequestAuthEmailRegister
		if ok, resp := common.ParamsDecode(c, &p); !ok {
			return resp
		}

		// Trim form
		p.Name = strings.TrimSpace(p.Name)
		p.Email = strings.TrimSpace(p.Email)
		p.Link = strings.TrimSpace(p.Link)
		p.Password = strings.TrimSpace(p.Password)

		// Check email
		if !utils.ValidateEmail(p.Email) {
			return common.RespError(c, 400, "Invalid email")
		}

		// Check link
		if p.Link != "" && !utils.ValidateURL(p.Link) {
			return common.RespError(c, 400, "Invalid link")
		}

		// Check password
		if len(p.Password) < 6 {
			return common.RespError(c, 400, "Password must be at least 6 characters")
		}

		// Select first user, if there are multiple users with same email
		// If name is empty, it is the mode of updating user (like forget password)
		if p.Name == "" {
			users := app.Dao().FindUsersByEmail(p.Email)
			if len(users) == 0 {
				return common.RespError(c, 400, "User not found")
			}
			p.Name = users[0].Name
		}

		// Check email verify code
		if ok, resp := CheckEmailVerifyCode(app, c, p.Email, p.Code); !ok {
			return resp
		}

		// Create user
		user, err := app.Dao().FindCreateUser(p.Name, p.Email, p.Link)
		if err != nil {
			return common.RespError(c, 500, "Failed to create user")
		}

		// Update user
		if err := user.SetPasswordEncrypt(p.Password); err != nil {
			return common.RespError(c, 500, "Failed to encrypt password")
		}
		if p.Link != "" {
			user.Link = p.Link
		}
		if err := app.Dao().UpdateUser(&user); err != nil {
			return common.RespError(c, 500, "Failed to update user")
		}

		// Login
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
