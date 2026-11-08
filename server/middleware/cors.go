package middleware

import (
	"net/url"
	"strings"

	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CheckOriginTrusted(app *core.App, origin string) bool {
	for _, rawURL := range app.Dao().TrustedURLs(app.Conf().TrustedDomains) {
		if allowed, err := utils.URLOrigin(rawURL); err == nil && allowed == origin {
			return true
		}
		if matchWildcardOrigin(rawURL, origin) {
			return true
		}
	}
	return false
}

// Wildcards grant browser access to one subdomain level. They never authorize
// page-title fetching, which uses the exact origins from Dao.TrustedOrigins.
func matchWildcardOrigin(pattern, origin string) bool {
	p, err := url.Parse(strings.TrimSpace(pattern))
	if err != nil || (p.Scheme != "http" && p.Scheme != "https") || p.User != nil {
		return false
	}
	host := strings.ToLower(p.Hostname())
	if !strings.HasPrefix(host, "*.") || len(host) == 2 || strings.Count(host, "*") != 1 || strings.HasSuffix(p.Host, ":") {
		return false
	}

	u, err := url.Parse(origin)
	if err != nil || u.Scheme != p.Scheme || u.Port() != p.Port() || strings.HasSuffix(u.Host, ":") || origin != u.Scheme+"://"+u.Host {
		return false
	}
	suffix := host[1:]
	subdomain, found := strings.CutSuffix(strings.ToLower(u.Hostname()), suffix)
	if !found || subdomain == "" || len(subdomain) > 63 || strings.HasPrefix(subdomain, "-") || strings.HasSuffix(subdomain, "-") {
		return false
	}
	for _, c := range subdomain {
		if !(c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '-') {
			return false
		}
	}
	return true
}

func CheckURLTrusted(app *core.App, targetUrl string) (trusted bool, origin string, err error) {
	origin, err = utils.URLOrigin(targetUrl)
	if err != nil {
		return false, "", err
	}
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
