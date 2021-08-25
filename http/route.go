package http

import (
	"net/http"

	"github.com/ArtalkJS/Artalk-API-Go/config"
	"github.com/ArtalkJS/Artalk-API-Go/model"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/markbates/pkger"
	"github.com/sirupsen/logrus"
)

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

	api.GET("/add", ActionAdd)
	api.GET("/get", ActionAdd)
	api.GET("/user", ActionUser)
	api.GET("/login-admin", ActionAdminLogin)

	// api/captcha
	ca := api.Group("/captcha")
	ca.GET("/refresh", ActionCaptchaRefresh)
	ca.GET("/check", ActionCaptchaCheck)

	// api/admin
	admin := api.Group("/admin")
	admin.GET("/edit", ActionAdminEdit)
	admin.GET("/del", ActionAdminDel)

	admin.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte(config.Instance.AppKey),
	}))
}

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type jwtCustomClaims struct {
	UserName  string         `json:"name"`
	UserEmail string         `json:"email"`
	UserType  model.UserType `json:"type"`
	jwt.StandardClaims
}
