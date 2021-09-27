package http

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func AdminOnly(c echo.Context) (isOK bool, resp error) {
	if !CheckIsAdminInJwtMiddleware(c) {
		return false, RespError(c, "需要验证管理员身份", Map{"need_login": true})
	}

	return true, nil
}

// 中间件会创建一个 user context，通过中间件获取到的 jwt 判断
func CheckIsAdminInJwtMiddleware(c echo.Context) bool {
	jwt := c.Get("user").(*jwt.Token)

	return CheckIsAdminByJwt(jwt)
}
