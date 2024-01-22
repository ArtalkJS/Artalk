package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

// @Id           Conf
// @Summary      Get System Configs
// @Description  Get System Configs for UI
// @Tags         System
// @Produce      json
// @Success      200  {object}  common.ConfData
// @Router       /conf  [get]
func Conf(app *core.App, router fiber.Router) {
	router.All("/conf", func(c *fiber.Ctx) error {
		return common.RespData(c, common.GetApiPublicConfDataMap(app, c))
	})
}
