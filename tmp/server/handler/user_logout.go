package handler

import (
	"time"

	"github.com/ArtalkJS/ArtalkGo/internal/config"
	"github.com/ArtalkJS/ArtalkGo/server/common"
	"github.com/gofiber/fiber/v2"
)

// POST /api/logout
func UserLogout(router fiber.Router) {
	router.Post("/logout", func(c *fiber.Ctx) error {

		if !config.Instance.Cookie.Enabled {
			return common.RespError(c, "API 未启用 Cookie")
		}

		if common.GetJwtStrByReqCookie(c) == "" {
			return common.RespError(c, "未登录，无需注销")
		}

		// same as login, remove cookie
		setAuthCookie(c, "", time.Now().AddDate(0, 0, -1))

		return common.RespSuccess(c)
	})
}
