package http

import (
	"fmt"

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
	// CORS 配置
	// for Preflight Request
	// 非法 Origin 浏览器拦截继续的请求
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true, // allow cors with cookies
		AllowOriginFunc: func(origin string) (bool, error) {
			if lib.ContainsStr(config.Instance.TrustedDomains, "*") {
				return true, nil // 通配符关闭 origin 检测
			}

			allowURLs := []string{}
			allowURLs = append(allowURLs, config.Instance.TrustedDomains...) // 导入配置中的可信域名
			for _, site := range model.FindAllSitesCooked() {                // 导入数据库中的站点 urls
				allowURLs = append(allowURLs, site.Urls...)
			}

			if len(allowURLs) == 0 {
				// 无配置的情况全部放行
				// 如程序第一次运行的时候
				return true, nil
			}

			if GetIsAllowOrigin(origin, allowURLs) {
				return true, nil
			}

			return false, nil
		},
	}))
}
