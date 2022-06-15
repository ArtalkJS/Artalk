package http

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

// 站点隔离 & Origin 控制
func SiteOriginMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			siteName := c.FormValue("site_name")
			siteID := uint(0)
			var site *model.Site = nil

			siteAll := false

			// 请求站点名 == "__ATK_SITE_ALL" 时取消站点隔离
			if siteName == lib.ATK_SITE_ALL {
				if !CheckIsAdminReq(c) {
					return RespError(c, "仅管理员查询允许取消站点隔离")
				}

				siteAll = true
			} else {
				// 请求站点名为空，使用默认 site
				if siteName == "" {
					siteName = strings.TrimSpace(config.Instance.SiteDefault)
					if siteName != "" {
						model.FindCreateSite(siteName) // 默认站点不存在则创建
					}
				}

				findSite := model.FindSite(siteName)
				if findSite.IsEmpty() {
					return RespError(c, fmt.Sprintf("未找到站点：`%s`，请控制台创建站点", siteName), Map{
						"err_no_site": true,
					})
				}
				site = &findSite
				siteID = findSite.ID
			}

			// 检测 Origin 合法性 (防止 CSRF 攻击)
			if isOK, resp := CheckOrigin(c, site); !isOK {
				return resp
			}

			// 设置上下文
			c.Set(lib.CTX_KEY_ATK_SITE_ID, siteID)
			c.Set(lib.CTX_KEY_ATK_SITE_NAME, siteName)
			c.Set(lib.CTX_KEY_ATK_SITE_ALL, siteAll)

			return next(c)
		}
	}
}

func UseSite(c echo.Context, siteName *string, destID *uint, destSiteAll *bool) {
	if destID != nil {
		*destID = c.Get(lib.CTX_KEY_ATK_SITE_ID).(uint)
	}
	if siteName != nil {
		*siteName = c.Get(lib.CTX_KEY_ATK_SITE_NAME).(string)
	}
	if destSiteAll != nil {
		*destSiteAll = c.Get(lib.CTX_KEY_ATK_SITE_ALL).(bool)
	}
}

// 检测 Origin 合法性
// 防止 CSRF 攻击
// @see https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html
func CheckOrigin(c echo.Context, allowSite *model.Site) (bool, error) {
	// 可信来源 URL
	allowSrcURLs := []string{}

	// 用户配置
	allowSrcURLs = append(allowSrcURLs, config.Instance.TrustedDomains...) // 允许配置文件域名
	if allowSite != nil {
		allowSrcURLs = append(allowSrcURLs, allowSite.ToCooked().Urls...) // 允许数据库站点 URLs 中的域名
	}
	if len(allowSrcURLs) == 0 {
		return true, nil // 若用户配置列表中无数据，则取消控制
	}
	if lib.ContainsStr(allowSrcURLs, "*") {
		return true, nil // 列表中出现通配符关闭控制
	}

	host := c.Request().Host
	realHostUnderProxy := c.Request().Header.Get("X-Forwarded-Host")
	if realHostUnderProxy != "" {
		host = realHostUnderProxy
	}

	// 读取 Origin 数据
	// @note Origin 标头在前端 fetch POST 操作中总是携带的，
	// 		 即使配置 Referrer-Policy: no-referrer
	// @see https://stackoverflow.com/questions/42239643/when-do-browsers-send-the-origin-header-when-do-browsers-set-the-origin-to-null
	origin := c.Request().Header.Get(echo.HeaderOrigin)
	if origin == "" || origin == "null" {
		// 从 Referer 获取 Origin
		referer := c.Request().Referer()
		if referer == "" {
			return false, RespError(c, "无效请求，Origin 无法获取")
		}
		origin = referer
	}

	pOrigin, err := url.Parse(origin)
	if err != nil {
		return false, RespError(c, "Origin 不合法")
	}

	// 系统配置：默认允许来自相同域名的请求
	allowSrcURLs = append(allowSrcURLs, c.Scheme()+"://"+host)

	allowSrcURLs = lib.RemoveDuplicates(allowSrcURLs) // 去重
	for _, a := range allowSrcURLs {
		a = strings.TrimSpace(a)
		if a == "" {
			continue
		}
		pAllow, err := url.Parse(a)
		if err != nil {
			continue
		}

		// 在可信来源列表中匹配 Referer 的 host 部分 (含端口) 则放行
		// @see https://web.dev/referrer-best-practices/
		// Referrer-Policy 不能设为 no-referer，
		// Chrome v85+ 默认为：strict-origin-when-cross-origin。
		// 前端页面 head 不配置 <meta name="referrer" content="no-referer" />，
		// 浏览器默认都会至少携带 Origin 数据 (不带 path，但包含端口)
		if pAllow.Host == pOrigin.Host {
			return true, nil
		}
	}

	return false, RespError(c, "非法请求，请检查可信域名配置")
}
