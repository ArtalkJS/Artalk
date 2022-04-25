package http

import (
	"net/http"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/pkged"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func InitRouter(e *echo.Echo) {
	// Init File Server
	f, err := pkged.Open("/frontend")
	if err != nil {
		logrus.Fatal(err)
		return
	}

	fileServer := http.FileServer(f)
	e.Any("/*", echo.WrapHandler(fileServer))

	// All Actions
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
	admin.POST("/vote-sync", action.AdminVoteSync)
	admin.POST("/send-mail", action.AdminSendMail)

	admin.POST("/cache-flush", action.AdminCacheFlush)
	admin.POST("/cache-warm", action.AdminCacheWarm)

	// conf
	api.Any("/conf", func(c echo.Context) error {
		return RespData(c, GetApiPublicConfDataMap(c))
	})

	// version
	api.Any("/version", func(c echo.Context) error {
		return c.JSON(200, GetApiVersionDataMap())
	})
}
