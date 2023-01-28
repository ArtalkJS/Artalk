package handler

import (
	"github.com/ArtalkJS/Artalk/internal/cache"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
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

type ResponseAdminUserEdit struct {
	User entity.CookedUserForAdmin `json:"user"`
}

// @Summary      User Edit
// @Description  Edit a specific user
// @Tags         User
// @Param        id             formData  string  true   "the user ID you want to edit"
// @Param        name           formData  string  true   "the user name"
// @Param        email          formData  string  true   "the user email"
// @Param        password       formData  string  false  "the user password"
// @Param        link           formData  string  false  "the user link"
// @Param        is_admin       formData  bool    true   "the user is an admin"
// @Param        site_names     formData  string  false  "the site names associated with the user"
// @Param        receive_email  formData  bool    true   "the user receive email"
// @Param        badge_name     formData  string  false  "the user badge name"
// @Param        badge_color    formData  string  false  "the user badge color (hex format)"
// @Security     ApiKeyAuth
// @Success      200  {object}  common.JSONResult{data=ResponseAdminUserEdit}
// @Router       /admin/user-edit  [post]
func AdminUserEdit(router fiber.Router) {
	router.Post("/user-edit", func(c *fiber.Ctx) error {
		if !common.GetIsSuperAdmin(c) {
			return common.RespError(c, i18n.T("Access denied"))
		}

		var p ParamsAdminUserEdit
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		user := query.FindUserByID(p.ID)
		if user.IsEmpty() {
			return common.RespError(c, i18n.T("{{name}} not found", Map{"name": i18n.T("User")}))
		}

		// 改名名合法性检测
		modifyName := p.Name != user.Name
		modifyEmail := p.Email != user.Email

		if modifyName && modifyEmail && !query.FindUser(p.Name, p.Email).IsEmpty() {
			return common.RespError(c, i18n.T("{{name}} already exists", Map{"name": i18n.T("User")}))
		}

		if !utils.ValidateEmail(p.Email) {
			return common.RespError(c, i18n.T("Invalid {{name}}", Map{"name": i18n.T("Email")}))
		}
		if p.Link != "" && !utils.ValidateURL(p.Link) {
			return common.RespError(c, i18n.T("Invalid {{name}}", Map{"name": i18n.T("Link")}))
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
			return common.RespError(c, i18n.T("{{name}} save failed", Map{"name": i18n.T("User")}))
		}

		return common.RespData(c, ResponseAdminUserEdit{
			User: query.UserToCookedForAdmin(&user),
		})
	})
}
