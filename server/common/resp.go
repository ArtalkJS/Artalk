package common

import (
	"net/http"

	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// JSONResult JSON Response data structure
type JSONResult struct {
	Msg  string      `json:"msg,omitempty"`  // Message
	Data interface{} `json:"data,omitempty"` // Data
}

// RespData is just response data
func RespData(c *fiber.Ctx, data interface{}) error {
	return c.Status(http.StatusOK).JSON(data)
}

// RespSuccess is just response success
func RespSuccess(c *fiber.Ctx, msg ...string) error {
	respData := &JSONResult{
		Msg: "Success",
	}

	// Optional msg parameter
	if len(msg) > 0 {
		respData.Msg = msg[0]
	}

	return c.Status(http.StatusOK).JSON(respData)
}

// RespError is just response error
func RespError(c *fiber.Ctx, code int, msg string, data ...Map) error {
	respData := Map{}
	if len(data) > 0 {
		respData = data[0]
	}

	respData["msg"] = msg

	c.Status(code)

	LogWithHttpInfo(c, func(l *zap.SugaredLogger) {
		l.Error(msg)
	})

	return c.JSON(respData)
}

func LogWithHttpInfo(c *fiber.Ctx, logFn func(l *zap.SugaredLogger)) {
	path := string(c.Request().URI().Path())
	if path == "" {
		path = "/"
	}

	id := c.Get(fiber.HeaderXRequestID)
	if id == "" {
		id = c.GetRespHeader(fiber.HeaderXRequestID)
	}

	logger := log.StandardLogger().Sugar().With(
		"id", id,
		"path", path,
		"method", c.Method(),
		"ip", c.IP(),
		"remote_addr", c.Context().RemoteAddr().String(),
		"host", string(c.Request().Host()),
		"referer", string(c.Request().Header.Referer()),
		"user_agent", string(c.Request().Header.UserAgent()),
		"status", c.Response().StatusCode(),
		// "headers", req.Header,
	)

	var skipper = func(l *zap.SugaredLogger) *zap.SugaredLogger {
		return l.WithOptions(zap.AddCallerSkip(2))
	}

	logFn(skipper(logger))
}
