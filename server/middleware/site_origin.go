package middleware

import (
	"net/url"
	"path"
	"strings"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

// 不启用 Origin 控制的 API paths
var SiteOriginSkips = []string{
	"/api/user-get",
	"/api/login",
}

// 站点隔离 & Origin 控制
func SiteOriginMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 忽略白名单
		for _, p := range SiteOriginSkips {
			if path.Clean(c.Path()) == path.Clean(p) {
				return c.Next()
			}
		}

		siteName := c.FormValue("site_name")
		siteID := uint(0)
		var site *entity.Site = nil

		siteAll := false
		isSuperAdmin := common.GetIsSuperAdmin(c)

		// 请求站点名 == "__ATK_SITE_ALL" 时取消站点隔离
		if siteName == config.ATK_SITE_ALL {
			if !isSuperAdmin {
				return common.RespError(c, "Only admin can query sites with disable isolation")
			}

			siteAll = true
		} else {
			// 请求站点名为空，使用默认 site
			if siteName == "" {
				siteName = strings.TrimSpace(config.Instance.SiteDefault)
				if siteName != "" {
					query.FindCreateSite(siteName) // 默认站点不存在则创建
				}
			}

			findSite := query.FindSite(siteName)
			if findSite.IsEmpty() {
				return common.RespError(c,
					i18n.T("Site `{{name}}` not found. Please create it in control center.", map[string]interface{}{"name": siteName}),
					common.Map{
						"err_no_site": true,
					},
				)
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
		c.Locals(config.CTX_KEY_ATK_SITE_ID, siteID)
		c.Locals(config.CTX_KEY_ATK_SITE_NAME, siteName)
		c.Locals(config.CTX_KEY_ATK_SITE_ALL, siteAll)

		return c.Next()
	}
}

// 检测 Origin 合法性
// 防止跨域的 CSRF 攻击
// @see https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html
func CheckOrigin(c *fiber.Ctx, allowSite *entity.Site) (bool, error) {
	// 可信来源 URL
	allowURLs := []string{}

	// 用户配置
	allowURLs = append(allowURLs, config.Instance.TrustedDomains...) // 允许配置文件域名
	if allowSite != nil {
		allowURLs = append(allowURLs, query.CookSite(allowSite).Urls...) // 允许数据库站点 URLs 中的域名
	}
	if utils.ContainsStr(allowURLs, "*") {
		return true, nil // 列表中出现通配符关闭控制
	}

	// 读取 Origin 数据
	// @note Origin 标头在前端 fetch POST 操作中总是携带的，
	// 		 即使配置 Referrer-Policy: no-referrer
	// @see https://stackoverflow.com/questions/42239643/when-do-browsers-send-the-origin-header-when-do-browsers-set-the-origin-to-null
	origin := c.Get(fiber.HeaderOrigin)
	if origin == "" || origin == "null" {
		// 从 Referer 获取 Origin
		referer := string(c.Request().Header.Referer())
		if referer == "" {
			return false, common.RespError(c, i18n.T("Invalid request")+", "+i18n.T("Unable to get `{{name}}`", map[string]interface{}{"name": "origin"}))
		}
		origin = referer
	}

	// 允许同源请求
	// c.Request().Host()、c.Request().URI().Scheme() 不能直接获取，
	// 如果套一层代理将不能获得正确值。
	// 使用 c.Protocol() 将读取 X-Forwarded-，应注意 header spoofing
	allowURLs = append(allowURLs, c.Protocol()+"://"+c.Hostname())

	// 判断 Origin 是否被允许
	if GetIsAllowOrigin(origin, allowURLs) {
		return true, nil
	}

	return false, common.RespError(c, i18n.T("Invalid request. Please check your `trusted_domains` config."))
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
