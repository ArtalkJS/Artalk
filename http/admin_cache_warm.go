package http

import (
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsAdminCacheWarm struct {
}

// 缓存预热
func (a *action) AdminCacheWarm(c echo.Context) error {
	var p ParamsAdminCacheWarm
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	if !GetIsSuperAdmin(c) {
		return RespError(c, "无权访问")
	}

	go func() {
		model.CacheWarmUp()
	}()

	return RespData(c, Map{
		"msg": "缓存预热任务已在后台开始执行，稍等片刻完成...",
	})
}
