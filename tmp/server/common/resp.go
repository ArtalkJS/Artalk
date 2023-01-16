package common

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// JSONResult JSON 响应数据结构
type JSONResult struct {
	Success bool        `json:"success"`         // 是否成功
	Msg     string      `json:"msg,omitempty"`   // 消息
	Data    interface{} `json:"data,omitempty"`  // 数据
	Extra   interface{} `json:"extra,omitempty"` // 数据
}

// RespJSON is normal json result
func RespJSON(c *fiber.Ctx, msg string, data interface{}, success bool) error {
	return c.Status(http.StatusOK).JSON(&JSONResult{
		Success: success,
		Msg:     msg,
		Data:    data,
	})
}

// RespData is just response data
func RespData(c *fiber.Ctx, data interface{}) error {
	return c.Status(http.StatusOK).JSON(&JSONResult{
		Success: true,
		Data:    data,
	})
}

// RespSuccess is just response success
func RespSuccess(c *fiber.Ctx, msg ...string) error {
	respData := &JSONResult{
		Success: true,
	}

	// 可选参数 msg
	if len(msg) > 0 {
		respData.Msg = msg[0]
	}

	return c.Status(http.StatusOK).JSON(respData)
}

// RespError is just response error
func RespError(c *fiber.Ctx, msg string, data ...Map) error {
	// log
	path := c.Path()
	if path == "" {
		path = "/"
	}
	LogWithHttpInfo(c).Errorf("[响应] %s %s ==> %s", c.Method(), path, strconv.Quote(msg))

	respData := Map{}
	if len(data) > 0 {
		respData = data[0]
	}

	return c.Status(http.StatusOK).JSON(&JSONResult{
		Success: false,
		Msg:     msg,
		Data:    respData,
	})
}

func LogWithHttpInfo(c *fiber.Ctx) *logrus.Entry {
	fields := logrus.Fields{}

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

	return logrus.WithFields(fields)
}
