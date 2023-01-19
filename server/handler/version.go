package handler

import (
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

// GET /version
func Version(router fiber.Router) {
	router.All("/version", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(common.GetApiVersionDataMap())
	})
}
