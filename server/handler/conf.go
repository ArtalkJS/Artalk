package handler

import (
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

// @Summary      Config
// @Description  Get system configurations
// @Tags     System
// @Success  200  {object}  common.JSONResult{data=config.Config}
// @Router   /conf  [get]
func Conf(router fiber.Router) {
	router.All("/conf", func(c *fiber.Ctx) error {
		return common.RespData(c, common.GetApiPublicConfDataMap(c))
	})
}
