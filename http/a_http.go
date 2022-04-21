package http

import (
	"fmt"
	"net/http"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/pkged"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echolog "github.com/onrik/logrus/echo"
	"github.com/sirupsen/logrus"
)

type Map = map[string]interface{}

func Run() {
	e := echo.New()
	e.HideBanner = true

	// 跨域控制
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

	// Action Limit
	ActionLimitConf := ActionLimitConf{
		// 启用操作限制路径白名单
		ProtectPaths: []string{
			"/api/add",
			"/api/login",
			"/api/vote",
			"/api/img-upload",
		},
	}

	if config.Instance.Captcha.Enabled {
		e.Use(ActionLimitMiddleware(ActionLimitConf))
	}

	CommonJwtConfig = middleware.JWTConfig{
		Claims:        &jwtCustomClaims{},
		ContextKey:    "user",
		SigningMethod: "HS256",
		TokenLookup:   "header:Authorization,query:token,form:token",
		AuthScheme:    "Bearer",
		SigningKey:    []byte(config.Instance.AppKey),
	}

	// Route
	InitRoute(e)

	listenAddr := fmt.Sprintf("%s:%d", config.Instance.Host, config.Instance.Port)

	if config.Instance.SSL.Enabled {
		e.Logger.Fatal(e.StartTLS(listenAddr, config.Instance.SSL.CertPath, config.Instance.SSL.KeyPath))
	} else {
		e.Logger.Fatal(e.Start(listenAddr))
	}
}

func InitRoute(e *echo.Echo) {
	f, err := pkged.Open("/frontend")
	if err != nil {
		logrus.Fatal(err)
		return
	}

	fileServer := http.FileServer(f)
	e.Any("/*", echo.WrapHandler(fileServer))

	// action
	action := &action{
		db: lib.DB,
	}

	// api
	api := e.Group("/api")

	api.POST("/add", action.Add)
	api.POST("/get", action.Get)
	api.POST("/user-get", action.UserGet)
	api.GET("/login", action.Login)
	api.POST("/login", action.Login)
	api.GET("/login-status", action.LoginStatus)
	api.POST("/login-status", action.LoginStatus)
	api.POST("/mark-read", action.MarkRead)
	api.POST("/vote", action.Vote)
	api.POST("/pv", action.PV)

	// api/upload-img
	if config.Instance.ImgUpload.Path == "" {
		config.Instance.ImgUpload.Path = "./data/artalk-img/"
		logrus.Warn("图片上传功能 img_upload.path 未配置，使用默认值：" + config.Instance.ImgUpload.Path)
	}
	api.POST("/img-upload", action.ImgUpload)
	e.Static(ImgUpload_RoutePath, config.Instance.ImgUpload.Path) // 静态可访问图片存放目录

	// api/captcha
	if config.Instance.Captcha.Enabled {
		ca := api.Group("/captcha")
		ca.GET("/refresh", action.CaptchaGet)
		ca.POST("/refresh", action.CaptchaGet)
		ca.GET("/get", action.CaptchaGet)
		ca.POST("/get", action.CaptchaGet)
		ca.GET("/check", action.CaptchaCheck)
		ca.POST("/check", action.CaptchaCheck)
		ca.GET("/status", action.CaptchaStatus)
		ca.POST("/status", action.CaptchaStatus)
	}

	// api/admin
	admin := api.Group("/admin", middleware.JWTWithConfig(CommonJwtConfig)) // use jwt

	admin.POST("/comment-edit", action.AdminCommentEdit)
	admin.POST("/comment-del", action.AdminCommentDel)

	admin.POST("/page-get", action.AdminPageGet)
	admin.POST("/page-edit", action.AdminPageEdit)
	admin.POST("/page-del", action.AdminPageDel)
	admin.POST("/page-fetch", action.AdminPageFetch)

	admin.POST("/site-get", action.AdminSiteGet)
	admin.POST("/site-add", action.AdminSiteAdd)
	admin.POST("/site-edit", action.AdminSiteEdit)
	admin.POST("/site-del", action.AdminSiteDel)
	admin.POST("/setting-get", action.AdminSettingGet)
	admin.POST("/setting-save", action.AdminSettingSave)
	admin.POST("/import", action.AdminImport)
	admin.POST("/import-upload", action.AdminImportUpload)
	admin.POST("/export", action.AdminExport)
	// admin.POST("/vote-sync", action.AdminVoteSync) // 数据导入功能未关注 vote 部分，暂时注释

	admin.POST("/send-mail", action.AdminSendMail)

	// conf
	api.Any("/conf", func(c echo.Context) error {
		return c.JSON(200, GetApiPublicConfDataMap(c))
	})

	// version
	api.Any("/version", func(c echo.Context) error {
		return c.JSON(200, GetApiVersionDataMap())
	})
}
