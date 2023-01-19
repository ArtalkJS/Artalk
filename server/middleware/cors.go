package middleware

import (
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/ArtalkJS/Artalk/server/middleware/cors"
	"github.com/gofiber/fiber/v2"
)

func CorsMiddleware() func(*fiber.Ctx) error {
	common.ReloadCorsAllowOrigins()
	return cors.New(common.CorsConf)
}
