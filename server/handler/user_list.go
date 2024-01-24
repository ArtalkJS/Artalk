package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsUserList struct {
	Limit  int `query:"limit" json:"limit" validate:"optional"`   // The limit for pagination
	Offset int `query:"offset" json:"offset" validate:"optional"` // The offset for pagination
}

type ResponseAdminUserList struct {
	Total int64                       `json:"count"`
	Users []entity.CookedUserForAdmin `json:"users"`
}

// @Id           GetUsers
// @Summary      Get User List
// @Description  Get a list of users by some conditions
// @Tags         User
// @Param        type     path   string              false  "The type of users"  Enums(all, admin, in_conf)
// @Param        options  query  ParamsUserList  true   "The options"
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponseAdminUserList
// @Failure      403  {object}  Map{msg=string}
// @Router       /users/{type}  [get]
func UserList(app *core.App, router fiber.Router) {
	router.Get("/users/:type?", common.AdminGuard(app, func(c *fiber.Ctx) error {
		listType := c.Params("type", "all") // 默认类型

		var p ParamsUserList
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

		return common.RespData(c, ResponseAdminUserList{
			Users: cookedUsers,
			Total: total,
		})
	}))
}
