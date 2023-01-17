package common

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/ArtalkJS/Artalk/server/middleware/cors"
)

var CorsConf = &cors.Config{
	AllowOrigins:     "",
	AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
	AllowCredentials: true, // allow cors with cookies
}

func ReloadCorsAllowOrigins() {
	allowOriginsArr := GetCorsAllowOrigins()

	allowOrigins := strings.Join(allowOriginsArr, ", ")
	{
		if len(allowOriginsArr) == 0 {
			// 无配置的情况全部放行
			// 如程序第一次运行的时候
			allowOrigins = "*"
		}
		if utils.ContainsStr(config.Instance.TrustedDomains, "*") {
			// 通配符关闭 origin 检测
			allowOrigins = "*"
		}
	}
	CorsConf.AllowOrigins = allowOrigins
}

func GetCorsAllowOrigins() []string {
	allowURLs := []string{}
	allowURLs = append(allowURLs, config.Instance.TrustedDomains...) // 导入配置中的可信域名
	for _, site := range query.FindAllSitesCooked() {                // 导入数据库中的站点 urls
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
