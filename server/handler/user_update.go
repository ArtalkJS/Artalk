package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/dao"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsUserUpdate struct {
	Name         string `json:"name" validate:"required"`          // The user name
	Email        string `json:"email" validate:"required"`         // The user email
	Password     string `json:"password" validate:"optional"`      // The user password
	Link         string `json:"link" validate:"optional"`          // The user link
	IsAdmin      bool   `json:"is_admin" validate:"required"`      // The user is an admin
	ReceiveEmail bool   `json:"receive_email" validate:"required"` // The user receive email
	BadgeName    string `json:"badge_name" validate:"optional"`    // The user badge name
	BadgeColor   string `json:"badge_color" validate:"optional"`   // The user badge color (hex format)
}

type ResponseUserUpdate struct {
	entity.CookedUserForAdmin
}

// @Id           UpdateUser
// @Summary      Update User
// @Description  Update a specific user
// @Tags         User
// @Param        id    path  int                  true  "The user ID you want to update"
// @Param        user  body  ParamsUserUpdate  true  "The user data"
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponseUserUpdate
// @Failure      400  {object}  Map{msg=string}
// @Failure      403  {object}  Map{msg=string}
// @Failure      404  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /users/{id}  [put]
func UserUpdate(app *core.App, router fiber.Router) {
	router.Put("/users/:id", common.AdminGuard(app, func(c *fiber.Ctx) error {
		id, _ := c.ParamsInt("id")

		var p ParamsUserUpdate
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		user := app.Dao().FindUserByID(uint(id))
		if user.IsEmpty() {
			return common.RespError(c, 404, i18n.T("{{name}} not found", Map{"name": i18n.T("User")}))
		}

		// 改名名合法性检测
		modifyName := p.Name != user.Name
		modifyEmail := p.Email != user.Email

		if modifyName && modifyEmail && !app.Dao().FindUser(p.Name, p.Email).IsEmpty() {
			return common.RespError(c, 400, i18n.T("{{name}} already exists", Map{"name": i18n.T("User")}))
		}

		if !utils.ValidateEmail(p.Email) {
			return common.RespError(c, 400, i18n.T("Invalid {{name}}", Map{"name": i18n.T("Email")}))
		}
		if p.Link != "" && !utils.ValidateURL(p.Link) {
			return common.RespError(c, 400, i18n.T("Invalid {{name}}", Map{"name": i18n.T("Link")}))
		}

		// 删除原有缓存
		app.Dao().CacheAction(func(cache *dao.DaoCache) {
			cache.UserCacheDel(&user)
		})

		// 修改 user
		user.Name = p.Name
		user.Email = p.Email
		if p.Password != "" {
			user.SetPasswordEncrypt(p.Password)
		}
		user.Link = p.Link
		user.IsAdmin = p.IsAdmin
		user.ReceiveEmail = p.ReceiveEmail
		user.BadgeName = p.BadgeName
		user.BadgeColor = p.BadgeColor

		err := app.Dao().UpdateUser(&user)
		if err != nil {
			return common.RespError(c, 500, i18n.T("{{name}} save failed", Map{"name": i18n.T("User")}))
		}

		return common.RespData(c, ResponseUserUpdate{
			CookedUserForAdmin: app.Dao().UserToCookedForAdmin(&user),
		})
	}))
}
