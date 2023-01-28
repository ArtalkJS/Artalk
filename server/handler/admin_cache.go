package handler

import (
	"github.com/ArtalkJS/Artalk/internal/cache"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsAdminCacheWarm struct {
}

// @Summary      Cache Warming
// @Description  Cache warming helps you hit the cache on the user's first visit
// @Tags         Cache
// @Security     ApiKeyAuth
// @Success      200  {object}  common.JSONResult
// @Router       /admin/cache-warm  [post]
func AdminCacheWarm(router fiber.Router) {
	router.Post("/cache-warm", func(c *fiber.Ctx) error {
		var p ParamsAdminCacheWarm
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		if !common.GetIsSuperAdmin(c) {
			return common.RespError(c, i18n.T("Access denied"))
		}

		go func() {
			cache.CacheWarmUp()
		}()

		return common.RespData(c, common.Map{
			"msg": i18n.T("Task executing in background, please wait..."),
		})
	})
}

type ParamsAdminCacheFlush struct {
	FlushAll bool `form:"flush_all"`
}

// @Summary      Cache Flush
// @Description  Flush Cache when application runs
// @Tags         Cache
// @Param        flush_all      formData  int     false  "flush all cache" example(1)
// @Security     ApiKeyAuth
// @Success      200  {object}  common.JSONResult
// @Router       /admin/cache-flush  [post]
func AdminCacheFlush(router fiber.Router) {
	router.Post("/cache-flush", func(c *fiber.Ctx) error {
		var p ParamsAdminCacheFlush
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		if !common.GetIsSuperAdmin(c) {
			return common.RespError(c, i18n.T("Access denied"))
		}

		if p.FlushAll {
			go func() {
				cache.CacheFlushAll()
			}()

			return common.RespData(c, common.Map{
				"msg": i18n.T("Task executing in background, please wait..."),
			})
		}

		return common.RespError(c, i18n.T("Invalid {{name}}", Map{"name": i18n.T("Parameter")}))
	})
}
