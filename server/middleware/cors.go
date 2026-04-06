package middleware

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func getCorsAllowOrigins(app *core.App) []string {
	allowURLs := []string{}
	allowURLs = append(allowURLs, app.Conf().TrustedDomains...) // 导入配置中的可信域名
	for _, site := range app.Dao().FindAllSitesCooked() {       // 导入数据库中的站点 urls
		allowURLs = append(allowURLs, site.Urls...)
	}

	allowOrigins := []string{}
	for _, u := range allowURLs {
		u = strings.TrimSpace(u)
		if u == "" {
			continue
		}

		urlP, err := url.Parse(u)
		if err != nil || urlP.Scheme == "" || urlP.Host == "" {
			continue
		}

		allowOrigins = append(allowOrigins, fmt.Sprintf("%s://%s", urlP.Scheme, urlP.Host))
	}

	return allowOrigins
}

func CheckOriginTrusted(app *core.App, origin string) bool {
	for _, allowed := range getCorsAllowOrigins(app) {
		if allowed == origin {
			return true
		}
		// Support wildcard subdomain matching (e.g. "https://*.example.com")
		if matchWildcardOrigin(allowed, origin) {
			return true
		}
	}
	return false
}

// matchWildcardOrigin checks if origin matches a wildcard pattern like "https://*.example.com".
// The wildcard only matches a single subdomain level in the host part.
func matchWildcardOrigin(pattern, origin string) bool {
	// Pattern must contain "*." in the host part
	schemeEnd := strings.Index(pattern, "://")
	if schemeEnd == -1 {
		return false
	}
	patternScheme := pattern[:schemeEnd]
	patternHost := pattern[schemeEnd+3:]

	if !strings.HasPrefix(patternHost, "*.") {
		return false
	}

	// Parse the origin
	originSchemeEnd := strings.Index(origin, "://")
	if originSchemeEnd == -1 {
		return false
	}
	originScheme := origin[:originSchemeEnd]
	originHost := origin[originSchemeEnd+3:]

	// Schemes must match
	if patternScheme != originScheme {
		return false
	}

	// The origin host must end with the wildcard's base domain
	// e.g. pattern "*.example.com" should match "sub.example.com" but not "example.com" itself
	baseDomain := patternHost[1:] // ".example.com"
	if !strings.HasSuffix(originHost, baseDomain) {
		return false
	}

	// The part before the base domain must not contain a dot (single level only)
	subdomain := originHost[:len(originHost)-len(baseDomain)]
	if subdomain == "" || strings.Contains(subdomain, ".") {
		return false
	}

	return true
}

func CheckURLTrusted(app *core.App, targetUrl string) (trusted bool, origin string, err error) {
	u, err := url.Parse(targetUrl)
	if err != nil {
		return false, "", fmt.Errorf("invalid URL")
	}
	origin = u.Scheme + "://" + u.Host
	trusted = CheckOriginTrusted(app, origin)
	return trusted, origin, nil
}

func CorsMiddleware(app *core.App) func(*fiber.Ctx) error {
	return cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return CheckOriginTrusted(app, origin)
		},
	})
}
