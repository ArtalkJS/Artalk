package handler

import (
	"github.com/ArtalkJS/Artalk/internal/cache"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsAdminCacheWarm struct {
}

// 缓存预热
//
// POST /api/admin/cache-warm
func AdminCacheWarm(router fiber.Router) {
	router.Post("/cache-warm", func(c *fiber.Ctx) error {
		var p ParamsAdminCacheWarm
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		if !common.GetIsSuperAdmin(c) {
			return common.RespError(c, i18n.T("No access"))
		}

		go func() {
			cache.CacheWarmUp()
		}()

		return common.RespData(c, common.Map{
			"msg": "缓存预热任务已在后台开始执行，稍等片刻完成...",
		})
	})
}

type ParamsAdminCacheFlush struct {
	FlushAll bool `form:"flush_all"`
}

// 缓存清理
//
// POST /api/admin/cache-flush
func AdminCacheFlush(router fiber.Router) {
	router.Post("/cache-flush", func(c *fiber.Ctx) error {
		var p ParamsAdminCacheFlush
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		if !common.GetIsSuperAdmin(c) {
			return common.RespError(c, i18n.T("No access"))
		}

		if p.FlushAll {
			go func() {
				cache.CacheFlushAll()
			}()

			return common.RespData(c, common.Map{
				"msg": "缓存清理任务已在后台开始执行，稍等片刻完成...",
			})
		}

		return common.RespError(c, i18n.T("Invalid parameter"))
	})
}
