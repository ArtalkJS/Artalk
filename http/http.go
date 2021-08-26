package http

import (
	"net/http"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	e.Logger = echolog.NewLogger(logrus.StandardLogger(), "")
	e.Use(echolog.Middleware(echolog.DefaultConfig))

	InitRoute(e)

	e.Logger.Fatal(e.Start(config.Instance.HttpAddr))
}

// JSONResult JSON 响应数据结构
type JSONResult struct {
	Success bool        `json:"success"`         // 是否成功
	Msg     string      `json:"msg,omitempty"`   // 消息
	Data    interface{} `json:"data,omitempty"`  // 数据
	Extra   interface{} `json:"extra,omitempty"` // 数据
}

// RespJSON is normal json result
func RespJSON(c echo.Context, msg string, data interface{}, success bool) error {
	return c.JSON(http.StatusOK, &JSONResult{
		Success: success,
		Msg:     msg,
		Data:    data,
	})
}

// RespData is just response data
func RespData(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, &JSONResult{
		Success: true,
		Data:    data,
	})
}

// RespSuccess is just response success
func RespSuccess(c echo.Context) error {
	return c.JSON(http.StatusOK, &JSONResult{
		Success: true,
	})
}

// RespError is just response error
func RespError(c echo.Context, msg string, details ...string) error {
	extraMap := Map{}
	if details != nil {
		extraMap["errDetails"] = details
	}

	return c.JSON(http.StatusOK, &JSONResult{
		Success: false,
		Msg:     msg,
		Extra:   extraMap,
	})
}
