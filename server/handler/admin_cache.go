package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

// @Summary      Warm-Up Cache
// @Description  Cache warming helps you to pre-load the cache to improve the performance of the first request
// @Tags         Cache
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  common.JSONResult
// @Router       /cache/warmup  [post]
func AdminCacheWarm(app *core.App, router fiber.Router) {
	router.Post("/cache/warmup", func(c *fiber.Ctx) error {
		if !common.GetIsSuperAdmin(app, c) {
			return common.RespError(c, i18n.T("Access denied"))
		}

		if !app.Conf().Cache.Enabled {
			return common.RespError(c, "cache disabled")
		}

		go func() {
			app.Dao().CacheWarmUp()
		}()

		return common.RespData(c, common.Map{
			"msg": i18n.T("Task executing in background, please wait..."),
		})
	})
}

// @Summary      Flush Cache
// @Description  Flush all cache on the server
// @Tags         Cache
// @Security     ApiKeyAuth
// @Success      200  {object}  common.JSONResult
// @Produce      json
// @Router       /cache/flush  [post]
func AdminCacheFlush(app *core.App, router fiber.Router) {
	router.Post("/cache/flush", func(c *fiber.Ctx) error {
		if !common.GetIsSuperAdmin(app, c) {
			return common.RespError(c, i18n.T("Access denied"))
		}

		if !app.Conf().Cache.Enabled {
			return common.RespError(c, "cache disabled")
		}

		go func() {
			app.Dao().CacheFlushAll()
		}()

		return common.RespData(c, common.Map{
			"msg": i18n.T("Task executing in background, please wait..."),
		})
	})
}
