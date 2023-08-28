package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

// @Summary      Version
// @Description  Get the version of Artalk
// @Tags         System
// @Success      200  {object}  common.ApiVersionData
// @Router       /version  [post]
func Version(app *core.App, router fiber.Router) {
	router.All("/version", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(common.GetApiVersionDataMap())
	})
}
