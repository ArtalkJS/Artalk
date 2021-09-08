package http

import (
	"net/http"
	"strings"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/markbates/pkger"
	echolog "github.com/onrik/logrus/echo"
	"github.com/sirupsen/logrus"
)

type Map = map[string]interface{}

func Run() {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: config.Instance.AllowOrigin,
	}))

	// Logger
	e.Logger = echolog.NewLogger(logrus.StandardLogger(), "")
	logConf := echolog.DefaultConfig
	logConf.Fields = []string{"ip", "latency", "status", "referer", "user_agent"} // "headers"
	e.Use(echolog.Middleware(logConf))

	// Action Limit
	ActionPermissionConf := ActionPermissionConf{
		Skipper: func(c echo.Context) bool {
			// 不启用操作限制的 path
			skipPath := []string{
				"/api/captcha/",
			}

			for _, p := range skipPath {
				if strings.HasPrefix(c.Request().URL.Path, p) {
					return true
				}
			}

			return false
		},
	}
	e.Use(ActionPermission(ActionPermissionConf))

	CommonJwtConfig = middleware.JWTConfig{
		Claims:        &jwtCustomClaims{},
		ContextKey:    "user",
		SigningMethod: "HS256",
		TokenLookup:   "header:X-Auth-Token,query:token",
		SigningKey:    []byte(config.Instance.AppKey),
	}

	// Route
	InitRoute(e)

	e.Logger.Fatal(e.Start(config.Instance.HttpAddr))
}

func InitRoute(e *echo.Echo) {
	f, err := pkger.Open("/frontend")
	if err != nil {
		logrus.Fatal(err)
		return
	}

	fileServer := http.FileServer(f)
	e.GET("/*", echo.WrapHandler(fileServer))

	// api
	api := e.Group("/api")

	api.POST("/add", ActionAdd)
	api.GET("/get", ActionGet)
	api.GET("/user", ActionUser)
	api.GET("/login", ActionLogin)
	api.POST("/login", ActionLogin)

	// api/captcha
	ca := api.Group("/captcha")
	ca.GET("/get", ActionCaptchaGet)
	ca.GET("/check", ActionCaptchaCheck)

	// api/manager
	manager := api.Group("/manager", middleware.JWTWithConfig(CommonJwtConfig)) // use jwt
	manager.GET("/edit", ActionManagerEdit)
	manager.GET("/del", ActionManagerDel)
	manager.GET("/send-mail", ActionManagerSendMail)
}
