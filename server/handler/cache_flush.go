package handler

import (
	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/i18n"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/gofiber/fiber/v2"
)

// @Id           FlushCache
// @Summary      Flush Cache
// @Description  Flush all cache on the server
// @Tags         Cache
// @Security     ApiKeyAuth
// @Success      200  {object}  Map{msg=string}
// @Failure      403  {object}  Map{msg=string}
// @Failure      400  {object}  Map{msg=string}
// @Produce      json
// @Router       /cache/flush  [post]
func CacheFlush(app *core.App, router fiber.Router) {
	router.Post("/cache/flush", common.AdminGuard(app, func(c *fiber.Ctx) error {
		if !app.Conf().Cache.Enabled {
			return common.RespError(c, 400, "cache disabled")
		}

		go func() {
			app.Dao().CacheFlushAll()
		}()

		return common.RespData(c, Map{
			"msg": i18n.T("Task executing in background, please wait..."),
		})
	}))
}
