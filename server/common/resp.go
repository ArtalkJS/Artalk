package common

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/gofiber/fiber/v2"
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
	// log
	path := c.Path()
	if path == "" {
		path = "/"
	}
	LogWithHttpInfo(c).Errorf("[Response] %s %s ==> %s", c.Method(), path, strconv.Quote(msg))

	respData := Map{}
	if len(data) > 0 {
		respData = data[0]
	}

	respData["msg"] = msg

	return c.Status(code).JSON(respData)
}

func LogWithHttpInfo(c *fiber.Ctx) *log.Entry {
	fields := log.Fields{}

	req := c.Request()
	res := c.Response()

	path := string(req.URI().Path())
	if path == "" {
		path = "/"
	}

	id := c.Get(fiber.HeaderXRequestID)
	if id == "" {
		id = c.GetRespHeader(fiber.HeaderXRequestID)
	}
	fields["id"] = id
	fields["ip"] = strings.Join(c.IPs(), ", ")
	fields["host"] = string(req.Host())
	fields["referer"] = string(req.Header.Referer())
	fields["user_agent"] = string(req.Header.UserAgent())
	fields["status"] = res.StatusCode()
	//fields["headers"] = req.Header

	return log.WithFields(fields)
}
