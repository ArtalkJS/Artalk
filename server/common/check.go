package common

import (
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func CheckIsAllowed(c *fiber.Ctx, name string, email string, page entity.Page, siteName string) (bool, error) {
	isAdminUser := query.IsAdminUserByNameEmail(name, email)

	// 如果用户是管理员，或者当前页只能管理员评论
	if isAdminUser || page.AdminOnly {
		if !CheckIsAdminReq(c) {
			return false, RespError(c, i18n.T("Admin access required"), Map{"need_login": true})
		}
	}

	return true, nil
}

func CheckIsAdminReq(c *fiber.Ctx) bool {
	jwt := GetJwtInstanceByReq(c)
	if jwt == nil {
		return false
	}

	user := GetUserByJwt(jwt)
	return user.IsAdmin
}

func GetIsSuperAdmin(c *fiber.Ctx) bool {
	user := GetUserByReq(c)
	return user.IsAdmin && user.SiteNames == ""
}

func IsAdminHasSiteAccess(c *fiber.Ctx, siteName string) bool {
	user := GetUserByReq(c)
	cookedUser := query.CookUser(&user)

	if !user.IsAdmin {
		return false
	}

	if !GetIsSuperAdmin(c) && !utils.ContainsStr(cookedUser.SiteNames, siteName) {
		// 如果账户分配了站点，并且待操作的站点并非处于分配的站点列表
		return false
	}

	return true
}
