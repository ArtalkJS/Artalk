package http

import (
	"github.com/labstack/echo/v4"
)

func ActionAdminLogin(c echo.Context) error {
	return nil
}

func ActionAdminEdit(c echo.Context) error {
	// user := c.Get("user").(*jwt.Token)
	// claims := user.Claims.(*jwtCustomClaims)
	// name := claims.Name
	return nil
}

func ActionAdminDel(c echo.Context) error {
	return nil
}
