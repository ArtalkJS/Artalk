package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsAdminUserGet struct {
	Limit  int `query:"limit" json:"limit"`   // The limit for pagination
	Offset int `query:"offset" json:"offset"` // The offset for pagination
}

type ResponseAdminUserGet struct {
	Total int64                       `json:"total"`
	Data  []entity.CookedUserForAdmin `json:"data"`
}

// @Summary      Get User List
// @Description  Get a list of users by some conditions
// @Tags         User
// @Param        type     path   string              false  "The type of users"  Enums(all, admin, in_conf)
// @Param        options  query  ParamsAdminUserGet  true   "The options"
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponseAdminUserGet
// @Failure      403  {object}  Map{msg=string}
// @Router       /users/{type}  [get]
func AdminUserGet(app *core.App, router fiber.Router) {
	router.Get("/users/:type?", func(c *fiber.Ctx) error {
		if !common.GetIsSuperAdmin(app, c) {
			return common.RespError(c, 403, i18n.T("Access denied"))
		}

		listType := c.Params("type", "all") // 默认类型

		var p ParamsAdminUserGet
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// 准备 query
		q := app.Dao().DB().Model(&entity.User{}).Order("created_at DESC")

		// 总共条数
		var total int64
		q.Count(&total)

		// 类型筛选
		if listType == "admin" {
			q = q.Where("is_admin = ?", true)
		} else if listType == "in_conf" {
			q = q.Where("is_in_conf = ?", true)
		}

		// 数据分页
		q = q.Scopes(Paginate(p.Offset, p.Limit))

		// 查找
		var users []entity.User
		q.Find(&users)

		var cookedUsers []entity.CookedUserForAdmin
		for _, u := range users {
			cookedUsers = append(cookedUsers, app.Dao().UserToCookedForAdmin(&u))
		}

		return common.RespData(c, ResponseAdminUserGet{
			Data:  cookedUsers,
			Total: total,
		})
	})
}
