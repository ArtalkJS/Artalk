package middleware

import (
	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"golang.org/x/exp/slices"
)

func CheckOriginTrusted(app *core.App, origin string) bool {
	return slices.Contains(app.Dao().TrustedOrigins(app.Conf().TrustedDomains), origin)
}

func CheckURLTrusted(app *core.App, targetUrl string) (trusted bool, origin string, err error) {
	origin, err = utils.URLOrigin(targetUrl)
	if err != nil {
		return false, "", err
	}
	trusted = CheckOriginTrusted(app, origin)
	return trusted, origin, nil
}

func CorsMiddleware(app *core.App) func(*fiber.Ctx) error {
	return cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return CheckOriginTrusted(app, origin)
		},
	})
}
