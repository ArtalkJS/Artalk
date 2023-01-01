package http

import (
	"github.com/ArtalkJS/ArtalkGo/internal/entity"
	"github.com/ArtalkJS/ArtalkGo/internal/query"
	"github.com/labstack/echo/v4"
)

type ParamsAdminUserGet struct {
	Limit  int    `mapstructure:"limit"`
	Offset int    `mapstructure:"offset"`
	Type   string `mapstructure:"type"`
}

type ResponseAdminUserGet struct {
	Total int64                       `json:"total"`
	Users []entity.CookedUserForAdmin `json:"users"`
}

func (a *action) AdminUserGet(c echo.Context) error {
	if !GetIsSuperAdmin(c) {
		return RespError(c, "无权操作")
	}

	var p ParamsAdminUserGet
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	// 准备 query
	db := a.db.Model(&entity.User{}).Order("created_at DESC")

	// 总共条数
	var total int64
	db.Count(&total)

	// 类型筛选
	if p.Type == "" {
		p.Type = "all" // 默认类型
	}

	if p.Type == "admin" {
		db = db.Where("is_admin = ?", true)
	} else if p.Type == "in_conf" {
		db = db.Where("is_in_conf = ?", true)
	}

	// 数据分页
	db = db.Scopes(Paginate(p.Offset, p.Limit))

	// 查找
	var users []entity.User
	db.Find(&users)

	var cookedUsers []entity.CookedUserForAdmin
	for _, u := range users {
		cookedUsers = append(cookedUsers, query.UserToCookedForAdmin(&u))
	}

	return RespData(c, ResponseAdminUserGet{
		Users: cookedUsers,
		Total: total,
	})
}
