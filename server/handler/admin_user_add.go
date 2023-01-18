package handler

import (
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ParamsAdminUserAdd struct {
	Name         string `form:"name" validate:"required"`
	Email        string `form:"email" validate:"required"`
	Password     string `form:"password"`
	Link         string `form:"link"`
	IsAdmin      bool   `form:"is_admin" validate:"required"`
	SiteNames    string `form:"site_names"`
	ReceiveEmail bool   `form:"receive_email" validate:"required"`
	BadgeName    string `form:"badge_name"`
	BadgeColor   string `form:"badge_color"`
}

// POST /api/admin/user-add
func AdminUserAdd(router fiber.Router) {
	router.Post("/user-add", func(c *fiber.Ctx) error {
		if !common.GetIsSuperAdmin(c) {
			return common.RespError(c, i18n.T("No access"))
		}

		var p ParamsAdminUserAdd
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		if !query.FindUser(p.Name, p.Email).IsEmpty() {
			return common.RespError(c, i18n.T("User already exists"))
		}

		if !utils.ValidateEmail(p.Email) {
			return common.RespError(c, "Invalid email")
		}
		if p.Link != "" && !utils.ValidateURL(p.Link) {
			return common.RespError(c, "Invalid link")
		}

		user := entity.User{}
		user.Name = p.Name
		user.Email = p.Email
		user.Link = p.Link
		user.IsAdmin = p.IsAdmin
		user.SiteNames = p.SiteNames
		user.ReceiveEmail = p.ReceiveEmail
		user.BadgeName = p.BadgeName
		user.BadgeColor = p.BadgeColor

		if p.Password != "" {
			err := user.SetPasswordEncrypt(p.Password)
			if err != nil {
				logrus.Errorln(err)
				return common.RespError(c, i18n.T("Failed to set password"))
			}
		}

		err := query.CreateUser(&user)
		if err != nil {
			return common.RespError(c, "user "+i18n.T("Failed to create"))
		}

		return common.RespData(c, common.Map{
			"user": query.UserToCookedForAdmin(&user),
		})
	})
}
