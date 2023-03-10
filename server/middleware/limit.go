package middleware

import (
	"path"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ActionLimitConf struct {
	ProtectPaths []string
}

// 操作限制 中间件
func ActionLimitMiddleware(conf ActionLimitConf) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 关闭验证码功能，直接 Skip
		if !config.Instance.Captcha.Enabled {
			return c.Next()
		}

		// 路径是否启用操作限制
		pathInList := false
		for _, p := range conf.ProtectPaths {
			if path.Clean(c.Path()) == path.Clean(p) {
				pathInList = true
				break
			}
		}
		if !pathInList {
			// 不启用的 path 直接放行
			return c.Next()
		}

		// 管理员直接放行
		if common.CheckIsAdminReq(c) {
			return c.Next()
		}

		userIP := c.IP()
		isNeedCheck := common.IsReqNeedCaptchaCheck(c)

		// 总是需要验证码模式
		if config.Instance.Captcha.Always {
			if common.GetAlwaysCaptchaMode_Pass(userIP) {
				common.SetAlwaysCaptchaMode_Pass(userIP, false) // 总是需要验证码，放行一次后再次需要验证码
				isNeedCheck = false
			} else {
				isNeedCheck = true
			}
		} else {
			// 超时模式：重置计数
			if config.Instance.Captcha.ActionReset != -1 {
				if !common.IsActionInTimeFrame(c) { // 超时
					common.ResetActionRecord(c) // 重置计数
					isNeedCheck = false         // 放行
				}
			}
		}

		// 是否需要验证
		if isNeedCheck {
			respData := common.Map{
				"need_captcha": true,
			}

			if config.Instance.Captcha.CaptchaType == config.TypeImage {
				// 图片验证码
				respData["img_data"] = common.GetNewImageCaptchaBase64(c.IP())
			} else {
				// iframe 验证模式
				respData["iframe"] = true
				// 前端新版不会再用到 img_data，向旧版响应 ArtalkFrontend out-of-date 图片
				respData["img_data"] = "data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 160 40'%3E%3Cdefs%3E%3Cstyle%3E.a%7Bfill:%23328ce6%3B%7D.b%7Bfont-size:12px%3Bfill:%23fff%3Bfont-family:sans-serif%3B%7D%3C/style%3E%3C/defs%3E%3Crect class='a' width='160' height='40'/%3E%3Ctext class='b' transform='translate(18.37 16.67)'%3EArtalk Frontend%3Ctspan x='0' y='14.4'%3EOut-Of-Date.%3C/tspan%3E%3C/text%3E%3C/svg%3E"
			}

			return common.RespError(c, i18n.T("Captcha required"), respData)
		}

		// 放行
		return c.Next()
	}
}
