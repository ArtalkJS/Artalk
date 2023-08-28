package middleware

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/ArtalkJS/Artalk/server/middleware/cors"
	"github.com/gofiber/fiber/v2"
)

func CorsMiddleware(app *core.App) func(*fiber.Ctx) error {
	common.ReloadCorsAllowOrigins(app)
	return cors.New(common.CorsConf)
}
