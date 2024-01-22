package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsUserCreate struct {
	Name         string `json:"name" validate:"required"`          // The user name
	Email        string `json:"email" validate:"required"`         // The user email
	Password     string `json:"password" validate:"optional"`      // The user password
	Link         string `json:"link" validate:"optional"`          // The user link
	IsAdmin      bool   `json:"is_admin" validate:"required"`      // The user is an admin
	ReceiveEmail bool   `json:"receive_email" validate:"required"` // The user receive email
	BadgeName    string `json:"badge_name" validate:"optional"`    // The user badge name
	BadgeColor   string `json:"badge_color" validate:"optional"`   // The user badge color (hex format)
}

type ResponseUserCreate struct {
	entity.CookedUserForAdmin
}

// @Id           CreateUser
// @Summary      Create User
// @Description  Create a new user
// @Tags         User
// @Param        user  body  ParamsUserCreate  true  "The user data"
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponseUserCreate
// @Failure      400  {object}  Map{msg=string}
// @Failure      403  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /users  [post]
func UserCreate(app *core.App, router fiber.Router) {
	router.Post("/users", common.AdminGuard(app, func(c *fiber.Ctx) error {
		var p ParamsUserCreate
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		if !app.Dao().FindUser(p.Name, p.Email).IsEmpty() {
			return common.RespError(c, 400, i18n.T("{{name}} already exists", Map{"name": i18n.T("User")}))
		}

		if !utils.ValidateEmail(p.Email) {
			return common.RespError(c, 400, i18n.T("Invalid {{name}}", Map{"name": i18n.T("Email")}))
		}
		if p.Link != "" && !utils.ValidateURL(p.Link) {
			return common.RespError(c, 400, i18n.T("Invalid {{name}}", Map{"name": i18n.T("Link")}))
		}

		user := entity.User{}
		user.Name = p.Name
		user.Email = p.Email
		user.Link = p.Link
		user.IsAdmin = p.IsAdmin
		user.ReceiveEmail = p.ReceiveEmail
		user.BadgeName = p.BadgeName
		user.BadgeColor = p.BadgeColor

		if p.Password != "" {
			err := user.SetPasswordEncrypt(p.Password)
			if err != nil {
				log.Errorln(err)
				return common.RespError(c, 500, i18n.T("Password update failed"))
			}
		}

		err := app.Dao().CreateUser(&user)
		if err != nil {
			return common.RespError(c, 500, i18n.T("{{name}} creation failed", Map{"name": i18n.T("User")}))
		}

		return common.RespData(c, ResponseUserCreate{
			CookedUserForAdmin: app.Dao().UserToCookedForAdmin(&user),
		})
	}))
}
