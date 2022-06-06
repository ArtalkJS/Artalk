package http

import (
	"fmt"
	libURL "net/url"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echolog "github.com/onrik/logrus/echo"
	"github.com/sirupsen/logrus"
)

func Run() {
	// Create echo instance
	e := echo.New()
	e.HideBanner = true

	// Cors Config
	InitCorsControl(e)

	// Logger
	e.Logger = echolog.NewLogger(logrus.StandardLogger(), "")
	logConf := echolog.DefaultConfig
	logConf.Fields = []string{"ip", "latency", "status", "referer", "user_agent"} // "headers"
	e.Use(echolog.Middleware(logConf))

	// Action Limit Middleware
	if config.Instance.Captcha.Enabled {
		ActionLimitConf := ActionLimitConf{
			// 启用操作限制路径白名单
			ProtectPaths: []string{
				"/api/add",
				"/api/login",
				"/api/vote",
				"/api/img-upload",
			},
		}

		e.Use(ActionLimitMiddleware(ActionLimitConf))
	}

	// Router
	InitRouter(e)

	// Start Listener
	listenAddr := fmt.Sprintf("%s:%d", config.Instance.Host, config.Instance.Port)
	if config.Instance.SSL.Enabled {
		e.Logger.Fatal(e.StartTLS(listenAddr, config.Instance.SSL.CertPath, config.Instance.SSL.KeyPath))
	} else {
		e.Logger.Fatal(e.Start(listenAddr))
	}
}

func InitCorsControl(e *echo.Echo) {
	allowOrigins := []string{}
	pushAllowOrigin := func(u string) {
		if !lib.ContainsStr(allowOrigins, u) {
			allowOrigins = append(allowOrigins, extractURLForCorsConf(u))
		}
	}

	// 导入配置中的可信域名
	for _, v := range config.Instance.TrustedDomains {
		pushAllowOrigin(v)
	}

	// 导入数据库中的站点 urls
	for _, site := range model.FindAllSitesCooked() {
		for _, url := range site.Urls {
			pushAllowOrigin(url)
		}
	}

	// 通配符关闭跨域控制
	if lib.ContainsStr(allowOrigins, "*") {
		allowOrigins = []string{"*"}
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowOrigins,
	}))
}

// 从完整 URL 中提取出 Scheme + Host 部分 (保留端口号)
func extractURLForCorsConf(u string) string {
	pu, err := libURL.Parse(u)
	if err != nil {
		return u
	}

	return pu.Scheme + "://" + pu.Host
}
