package http

import (
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsAdminCacheFlush struct {
	FlushAll bool `mapstructure:"flush_all"`
}

// 缓存清理
func (a *action) AdminCacheFlush(c echo.Context) error {
	var p ParamsAdminCacheFlush
	if isOK, resp := ParamsDecode(c, ParamsAdminCacheFlush{}, &p); !isOK {
		return resp
	}

	if p.FlushAll {
		go func() {
			model.CacheFlushAll()
		}()

		return RespData(c, Map{
			"msg": "缓存清理任务已在后台开始执行，稍等片刻完成...",
		})
	}

	return RespError(c, "参数错误")
}
