package http

import (
	"time"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/labstack/echo/v4"
)

func (a *action) Logout(c echo.Context) error {
	if !config.Instance.Cookie.Enabled {
		return RespError(c, "API 未启用 Cookie")
	}

	if GetJwtStrByReqCookie(c) == "" {
		return RespError(c, "未登录，无需注销")
	}

	// same as login, remove cookie
	setAuthCookie(c, "", time.Now().AddDate(0, 0, -1))

	return RespSuccess(c)
}
