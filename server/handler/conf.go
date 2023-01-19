package handler

import (
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

// GET /conf
func Conf(router fiber.Router) {
	router.All("/conf", func(c *fiber.Ctx) error {
		return common.RespData(c, common.GetApiPublicConfDataMap(c))
	})
}
