package http

import (
	"fmt"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
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
	allowOrigins := []string{}
	for _, v := range config.Instance.TrustedDomains { // 可信域名配置
		if !lib.ContainsStr(allowOrigins, v) {
			allowOrigins = append(allowOrigins, v)
		}
	}
	if lib.ContainsStr(allowOrigins, "*") { // 通配符关闭跨域控制
		allowOrigins = []string{"*"}
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowOrigins,
	}))

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

	// Jwt Config
	CommonJwtConfig = middleware.JWTConfig{
		Claims:        &jwtCustomClaims{},
		ContextKey:    "user",
		SigningMethod: "HS256",
		TokenLookup:   "header:Authorization,query:token,form:token",
		AuthScheme:    "Bearer",
		SigningKey:    []byte(config.Instance.AppKey),
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
