package common

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/ArtalkJS/Artalk/server/middleware/cors"
)

// TODO move cors conf to app instance neither global variable
var CorsConf = &cors.Config{
	AllowOrigins:     "",
	AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
	AllowCredentials: true, // allow cors with cookies
}

func ReloadCorsAllowOrigins(app *core.App) {
	allowOriginsArr := GetCorsAllowOrigins(app)

	allowOrigins := strings.Join(allowOriginsArr, ", ")
	{
		if len(allowOriginsArr) == 0 {
			// 无配置的情况全部放行
			// 如程序第一次运行的时候
			allowOrigins = "*"
		}
		if utils.ContainsStr(app.Conf().TrustedDomains, "*") {
			// 通配符关闭 origin 检测
			allowOrigins = "*"
		}
	}

	// TODO prevent use global variable
	CorsConf.AllowOrigins = allowOrigins
}

func GetCorsAllowOrigins(app *core.App) []string {
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
