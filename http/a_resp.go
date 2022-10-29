package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/labstack/echo/v4"
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
func RespSuccess(c echo.Context, msg ...string) error {
	respData := &JSONResult{
		Success: true,
	}

	// 可选参数 msg
	if len(msg) > 0 {
		respData.Msg = msg[0]
	}

	return c.JSON(http.StatusOK, respData)
}

// RespError is just response error
func RespError(c echo.Context, msg string, data ...Map) error {
	// log
	req := c.Request()
	path := req.URL.Path
	if path == "" {
		path = "/"
	}
	LogWithHttpInfo(c).Errorf("[响应] %s %s ==> %s", req.Method, path, strconv.Quote(msg))

	respData := Map{}
	if len(data) > 0 {
		respData = data[0]
	}

	return c.JSON(http.StatusOK, &JSONResult{
		Success: false,
		Msg:     msg,
		Data:    respData,
	})
}

func LogWithHttpInfo(c echo.Context) *logrus.Entry {
	fields := logrus.Fields{}

	req := c.Request()
	res := c.Response()

	path := req.URL.Path
	if path == "" {
		path = "/"
	}

	id := req.Header.Get(echo.HeaderXRequestID)
	if id == "" {
		id = res.Header().Get(echo.HeaderXRequestID)
	}
	fields["id"] = id
	fields["ip"] = c.RealIP()
	fields["host"] = req.Host
	fields["referer"] = req.Referer()
	fields["user_agent"] = req.UserAgent()
	fields["status"] = res.Status
	//fields["headers"] = req.Header

	return logrus.WithFields(fields)
}

func GetApiVersionDataMap() Map {
	return Map{
		"app":            "artalk-go",
		"version":        strings.TrimPrefix(lib.Version, "v"),
		"commit_hash":    lib.CommitHash,
		"fe_min_version": strings.TrimPrefix(lib.FeMinVersion, "v"),
	}
}

func GetApiPublicConfDataMap(c echo.Context) Map {
	isAdmin := CheckIsAdminReq(c)
	imgUpload := config.Instance.ImgUpload.Enabled
	if isAdmin {
		imgUpload = true // 管理员始终允许上传图片
	}

	frontendConf := config.Instance.Frontend
	if frontendConf == nil {
		frontendConf = make(map[string]interface{})
	}
	frontendConf["imgUpload"] = &imgUpload

	return Map{
		"img_upload":    imgUpload,
		"frontend_conf": frontendConf,
	}
}
