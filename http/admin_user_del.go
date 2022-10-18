package http

import (
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsAdminUserDel struct {
	ID uint `mapstructure:"id" param:"required"`
}

func (a *action) AdminUserDel(c echo.Context) error {
	if !GetIsSuperAdmin(c) {
		return RespError(c, "无权操作")
	}

	var p ParamsAdminUserDel
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	user := model.FindUserByID(p.ID)
	if user.IsEmpty() {
		return RespError(c, "user 不存在")
	}

	err := model.DelUser(&user)
	if err != nil {
		return RespError(c, "user 删除失败")
	}
	return RespSuccess(c)
}
