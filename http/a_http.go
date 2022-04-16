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
	allowOrigins = append(allowOrigins, config.Instance.AllowOrigins...) // 跨域独立配置
	for _, v := range config.Instance.TrustedDomains {                   // 可信域名配置
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
	e.Use(ActionLimitMiddleware(ActionLimitConf))

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

	// api
	api := e.Group("/api")

	api.POST("/add", ActionAdd)
	api.POST("/get", ActionGet)
	api.POST("/user-get", ActionUserGet)
	api.GET("/login", ActionLogin)
	api.POST("/login", ActionLogin)
	api.POST("/mark-read", ActionMarkRead)
	api.POST("/vote", ActionVote)
	api.POST("/pv", ActionPV)

	// api/upload-img
	if config.Instance.ImgUpload.Path == "" {
		config.Instance.ImgUpload.Path = "./data/artalk-img/"
		logrus.Warn("图片上传功能 img_upload.path 未配置，使用默认值：" + config.Instance.ImgUpload.Path)
	}
	api.POST("/img-upload", ActionImgUpload)
	e.Static(ImgUpload_RoutePath, config.Instance.ImgUpload.Path) // 静态可访问图片存放目录

	// api/captcha
	ca := api.Group("/captcha")
	ca.GET("/refresh", ActionCaptchaGet)
	ca.GET("/get", ActionCaptchaGet)
	ca.GET("/check", ActionCaptchaCheck)
	ca.GET("/status", ActionCaptchaStatus)

	// api/admin
	admin := api.Group("/admin", middleware.JWTWithConfig(CommonJwtConfig)) // use jwt

	admin.POST("/comment-edit", ActionAdminCommentEdit)
	admin.POST("/comment-del", ActionAdminCommentDel)

	admin.POST("/page-get", ActionAdminPageGet)
	admin.POST("/page-edit", ActionAdminPageEdit)
	admin.POST("/page-del", ActionAdminPageDel)
	admin.POST("/page-fetch", ActionAdminPageFetch)

	admin.POST("/site-get", ActionAdminSiteGet)
	admin.POST("/site-add", ActionAdminSiteAdd)
	admin.POST("/site-edit", ActionAdminSiteEdit)
	admin.POST("/site-del", ActionAdminSiteDel)
	admin.POST("/setting-get", ActionAdminSettingGet)
	admin.POST("/setting-save", ActionAdminSettingSave)
	admin.POST("/import", ActionAdminImport)
	admin.POST("/export", ActionAdminExport)
	// admin.POST("/vote-sync", ActionAdminVoteSync) // 数据导入功能未关注 vote 部分，暂时注释

	admin.POST("/send-mail", ActionAdminSendMail)

	// conf
	api.Any("/conf", func(c echo.Context) error {
		return c.JSON(200, GetApiPublicConfDataMap(c))
	})

	// version
	api.Any("/version", func(c echo.Context) error {
		return c.JSON(200, GetApiVersionDataMap())
	})
}
