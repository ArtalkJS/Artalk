package middleware

import (
	"github.com/ArtalkJS/ArtalkGo/server/common"
	"github.com/ArtalkJS/ArtalkGo/server/middleware/cors"
	"github.com/gofiber/fiber/v2"
)

func CorsMiddleware() func(*fiber.Ctx) error {
	common.ReloadCorsAllowOrigins()
	return cors.New(common.CorsConf)
}
