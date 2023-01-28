package handler

import (
	"github.com/ArtalkJS/Artalk/internal/db"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsAdminUserGet struct {
	Limit  int    `form:"limit"`
	Offset int    `form:"offset"`
	Type   string `form:"type"`
}

type ResponseAdminUserGet struct {
	Total int64                       `json:"total"`
	Users []entity.CookedUserForAdmin `json:"users"`
}

// @Summary      User List
// @Description  Get a list of users by some conditions
// @Tags         User
// @Param        limit          formData  int     false  "the limit for pagination"
// @Param        offset         formData  int     false  "the offset for pagination"
// @Security     ApiKeyAuth
// @Success      200  {object}  common.JSONResult{data=ResponseAdminUserGet}
// @Router       /admin/user-get  [post]
func AdminUserGet(router fiber.Router) {
	router.Post("/user-get", func(c *fiber.Ctx) error {
		if !common.GetIsSuperAdmin(c) {
			return common.RespError(c, i18n.T("Access denied"))
		}

		var p ParamsAdminUserGet
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// 准备 query
		q := db.DB().Model(&entity.User{}).Order("created_at DESC")

		// 总共条数
		var total int64
		q.Count(&total)

		// 类型筛选
		if p.Type == "" {
			p.Type = "all" // 默认类型
		}

		if p.Type == "admin" {
			q = q.Where("is_admin = ?", true)
		} else if p.Type == "in_conf" {
			q = q.Where("is_in_conf = ?", true)
		}

		// 数据分页
		q = q.Scopes(Paginate(p.Offset, p.Limit))

		// 查找
		var users []entity.User
		q.Find(&users)

		var cookedUsers []entity.CookedUserForAdmin
		for _, u := range users {
			cookedUsers = append(cookedUsers, query.UserToCookedForAdmin(&u))
		}

		return common.RespData(c, ResponseAdminUserGet{
			Users: cookedUsers,
			Total: total,
		})
	})
}
