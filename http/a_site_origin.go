package http

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

// 不启用 Origin 控制的 API paths
var SiteOriginSkips = []string{
	"/api/user-get",
	"/api/login",
}

// 站点隔离 & Origin 控制
func SiteOriginMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 忽略白名单
			for _, p := range SiteOriginSkips {
				if path.Clean(c.Request().URL.Path) == path.Clean(p) {
					return next(c)
				}
			}

			siteName := c.FormValue("site_name")
			siteID := uint(0)
			var site *model.Site = nil

			siteAll := false
			isSuperAdmin := GetIsSuperAdmin(c)

			// 请求站点名 == "__ATK_SITE_ALL" 时取消站点隔离
			if siteName == lib.ATK_SITE_ALL {
				if !isSuperAdmin {
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
					return RespError(c, fmt.Sprintf("未找到站点：`%s`，请在控制台创建站点", siteName), Map{
						"err_no_site": true,
					})
				}
				site = &findSite
				siteID = findSite.ID
			}

			// 检测 Origin 合法性 (防止跨域的 CSRF 攻击)
			if !isSuperAdmin { // 管理员忽略 Origin 检测
				if isOK, resp := CheckOrigin(c, site); !isOK {
					return resp
				}
			}

			// 设置 Context Values
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
// 防止跨域的 CSRF 攻击
// @see https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html
func CheckOrigin(c echo.Context, allowSite *model.Site) (bool, error) {
	// 可信来源 URL
	allowURLs := []string{}

	// 用户配置
	allowURLs = append(allowURLs, config.Instance.TrustedDomains...) // 允许配置文件域名
	if allowSite != nil {
		allowURLs = append(allowURLs, allowSite.ToCooked().Urls...) // 允许数据库站点 URLs 中的域名
	}
	if lib.ContainsStr(allowURLs, "*") {
		return true, nil // 列表中出现通配符关闭控制
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

	// 允许同源请求
	host := c.Request().Host
	realHostUnderProxy := c.Request().Header.Get("X-Forwarded-Host")
	if realHostUnderProxy != "" {
		host = realHostUnderProxy
	}
	allowURLs = append(allowURLs, c.Scheme()+"://"+host)

	// 判断 Origin 是否被允许
	if GetIsAllowOrigin(origin, allowURLs) {
		return true, nil
	}

	return false, RespError(c, "非法请求，请检查可信域名配置")
}

// 判断 Origin 是否被允许
// origin is 'schema://hostname:port',
// allowURLs is a collection of url strings
func GetIsAllowOrigin(origin string, allowURLs []string) bool {
	// Origin 合法性检测
	originP, err := url.Parse(origin)
	if err != nil || originP.Scheme == "" || originP.Host == "" {
		return false
	}

	// 提取 URLs 检测 Origin 是否匹配
	for _, u := range allowURLs {
		u = strings.TrimSpace(u)
		if u == "" {
			continue
		}

		urlP, err := url.Parse(u)
		if err != nil || urlP.Scheme == "" || urlP.Host == "" {
			continue
		}

		// 在可信来源列表中匹配 Referer 的 host 部分 (含端口) 则放行
		// @see https://web.dev/referrer-best-practices/
		// Referrer-Policy 不能设为 no-referer，
		// Chrome v85+ 默认为：strict-origin-when-cross-origin。
		// 前端页面 head 不配置 <meta name="referrer" content="no-referer" />，
		// 浏览器默认都会至少携带 Origin 数据 (不带 path，但包含端口)
		if urlP.Scheme == originP.Scheme && urlP.Host == originP.Host {
			return true
		}
	}

	return false
}
