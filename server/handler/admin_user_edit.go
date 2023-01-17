package handler

import (
	"github.com/ArtalkJS/Artalk/internal/cache"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsAdminUserEdit struct {
	// 查询值
	ID uint `form:"id" validate:"required"`

	// 修改值
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

// POST /api/admin/user-edit
func AdminUserEdit(router fiber.Router) {
	router.Post("/user-edit", func(c *fiber.Ctx) error {
		if !common.GetIsSuperAdmin(c) {
			return common.RespError(c, "无权操作")
		}

		var p ParamsAdminUserEdit
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		user := query.FindUserByID(p.ID)
		if user.IsEmpty() {
			return common.RespError(c, "user 不存在")
		}

		// 改名名合法性检测
		modifyName := p.Name != user.Name
		modifyEmail := p.Email != user.Email

		if modifyName && modifyEmail && !query.FindUser(p.Name, p.Email).IsEmpty() {
			return common.RespError(c, "user 已存在，请更换用户名和邮箱")
		}

		if !utils.ValidateEmail(p.Email) {
			return common.RespError(c, "Invalid email")
		}
		if p.Link != "" && !utils.ValidateURL(p.Link) {
			return common.RespError(c, "Invalid link")
		}

		// 删除原有缓存
		cache.UserCacheDel(&user)

		// 修改 user
		user.Name = p.Name
		user.Email = p.Email
		if p.Password != "" {
			user.SetPasswordEncrypt(p.Password)
		}
		user.Link = p.Link
		user.IsAdmin = p.IsAdmin
		user.SiteNames = p.SiteNames
		user.ReceiveEmail = p.ReceiveEmail
		user.BadgeName = p.BadgeName
		user.BadgeColor = p.BadgeColor

		err := query.UpdateUser(&user)
		if err != nil {
			return common.RespError(c, "user 保存失败")
		}

		return common.RespData(c, common.Map{
			"user": query.UserToCookedForAdmin(&user),
		})
	})
}
