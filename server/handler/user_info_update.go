package handler

import (
	"strings"

	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/i18n"
	"github.com/artalkjs/artalk/v2/internal/utils"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/gofiber/fiber/v2"
)

type RequestUserInfoUpdate struct {
	Email string `json:"email" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Link  string `json:"link" validate:"optional"`
	Code  string `json:"code" validate:"optional"`
}

type ResponseUserInfoUpdate struct {
	User entity.CookedUser `json:"user"`
}

// @Id           UpdateProfile
// @Summary      Update user profile
// @Description  Update user profile when user is logged in
// @Tags         Auth
// @Security     ApiKeyAuth
// @Param        data  body  RequestUserInfoUpdate  true  "The profile data to update"
// @Success      200  {object}  ResponseUserInfoUpdate
// @Failure      400  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Accept       json
// @Produce      json
// @Router       /user  [post]
func UserInfoUpdate(app *core.App, router fiber.Router) {
	router.Post("/user", common.LoginGuard(app, func(c *fiber.Ctx, user entity.User) error {
		var p RequestUserInfoUpdate
		if ok, resp := common.ParamsDecode(c, &p); !ok {
			return resp
		}

		// Trim form
		p.Name = strings.TrimSpace(p.Name)
		p.Email = strings.TrimSpace(p.Email)
		p.Link = strings.TrimSpace(p.Link)

		// Modify name
		if p.Name != user.Name {
			// Check name exists
			// (Allows the same name but different email)
			findUserByName := app.Dao().FindUser(p.Name, p.Email) // Find by `name` AND `email`
			if !findUserByName.IsEmpty() && findUserByName.ID != user.ID {
				// If user name exists but not current user
				// Allow current user to change the same name but case not same
				return common.RespError(c, 400, i18n.T("{{name}} already exists", map[string]interface{}{"name": i18n.T("Username")}))
			}

			user.Name = p.Name
		}

		// Modify email
		// (Verify sent email code first)
		if p.Email != user.Email {
			// Check email format
			if !utils.ValidateEmail(p.Email) {
				return common.RespError(c, 400, i18n.T("Invalid {{name}}", map[string]interface{}{"name": i18n.T("Email")}))
			}

			// Check email exists
			var findUserByEmail entity.User
			app.Dao().DB().Where("LOWER(email) = LOWER(?)", p.Email).First(&findUserByEmail)
			if !findUserByEmail.IsEmpty() {
				return common.RespError(c, 400, i18n.T("{{name}} already exists", map[string]interface{}{"name": i18n.T("Email")}))
			}

			// Check email verify code
			if ok, resp := CheckEmailVerifyCode(app, c, p.Email, p.Code); !ok {
				return resp
			}

			user.Email = p.Email
		}

		// Modify link
		if p.Link != "" {
			// Check link format
			if !utils.ValidateURL(p.Link) {
				user.Link = "https://" + p.Link
			}

			user.Link = p.Link
		}

		if err := app.Dao().UpdateUser(&user); err != nil {
			return common.RespError(c, 500, "Failed to update user")
		}

		return common.RespData(c, ResponseUserInfoUpdate{
			User: app.Dao().CookUser(&user),
		})
	}))
}
