package http

import (
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsAdminUserGet struct {
	Limit  int    `mapstructure:"limit"`
	Offset int    `mapstructure:"offset"`
	Type   string `mapstructure:"type"`
}

type ResponseAdminUserGet struct {
	Total int64                      `json:"total"`
	Users []model.CookedUserForAdmin `json:"users"`
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
	query := a.db.Model(&model.User{}).Order("created_at DESC")

	// 总共条数
	var total int64
	query.Count(&total)

	// 类型筛选
	if p.Type == "" {
		p.Type = "all" // 默认类型
	}

	if p.Type == "admin" {
		query = query.Where("is_admin = ?", true)
	} else if p.Type == "in_conf" {
		query = query.Where("is_in_conf = ?", true)
	}

	// 数据分页
	query = query.Scopes(Paginate(p.Offset, p.Limit))

	// 查找
	var users []model.User
	query.Find(&users)

	var cookedUsers []model.CookedUserForAdmin
	for _, u := range users {
		cookedUsers = append(cookedUsers, u.ToCookedForAdmin())
	}

	return RespData(c, ResponseAdminUserGet{
		Users: cookedUsers,
		Total: total,
	})
}
