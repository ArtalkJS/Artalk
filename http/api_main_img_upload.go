package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// 图片目录路由重写路径
var ImgUpload_RoutePath = "/static/images/"

type ParamsImgUpload struct {
	Name  string `mapstructure:"name" param:"required"`
	Email string `mapstructure:"email" param:"required"`

	PageKey   string `mapstructure:"page_key" param:"required"`
	PageTitle string `mapstructure:"page_title"`

	SiteName string `mapstructure:"site_name"`

	SiteID  uint
	SiteAll bool
}

func ActionImgUpload(c echo.Context) error {
	// 功能开关
	if !config.Instance.ImgUpload.Enabled {
		return RespError(c, "图片上传功能已禁用", Map{
			"img_upload_enabled": false,
		})
	}

	// 传入参数解析
	var p ParamsImgUpload
	if isOK, resp := ParamsDecode(c, ParamsImgUpload{}, &p); !isOK {
		return resp
	}

	if !lib.ValidateEmail(p.Email) {
		return RespError(c, "Invalid email")
	}

	// find site
	if isOK, resp := CheckSite(c, &p.SiteName, &p.SiteID, &p.SiteAll); !isOK {
		return resp
	}

	// 记录请求次数 (for 请求频率限制)
	RecordAction(c)

	// find page
	// page := model.FindPage(p.PageKey, p.PageTitle)
	// ip := c.RealIP()
	// ua := c.Request().UserAgent()

	// 图片大小限制 (Based on content length)
	if config.Instance.ImgUpload.MaxSize != 0 {
		if c.Request().ContentLength > config.Instance.ImgUpload.MaxSize*1024*1024 {
			return RespError(c, fmt.Sprintf("图片大小超过限制 %dMB", config.Instance.ImgUpload.MaxSize))
		}
	}

	// 获取 Form
	file, err := c.FormFile("file")
	if err != nil {
		logrus.Error(err)
		return RespError(c, "文件获取失败")
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		logrus.Error(err)
		return RespError(c, "文件打开失败")
	}
	defer src.Close()

	// 读取文件
	buf, err := ioutil.ReadAll(src)
	if err != nil {
		logrus.Error(err)
		return RespError(c, "文件读取失败")
	}

	// 大小限制 (Based on content read)
	if config.Instance.ImgUpload.MaxSize != 0 {
		if int64(len(buf)) > config.Instance.ImgUpload.MaxSize*1024*1024 {
			return RespError(c, fmt.Sprintf("图片大小超过限制 %dMB", config.Instance.ImgUpload.MaxSize))
		}
	}

	// 文件格式判断
	// @link https://mimesniff.spec.whatwg.org/
	// @link https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Basics_of_HTTP/MIME_types
	fileMine := http.DetectContentType(buf)
	allowMines := []string{
		"image/jpeg", "image/png", "image/gif",
		"image/svg+xml", "image/webp", "image/bmp",
	}
	if !lib.ContainsStr(allowMines, fileMine) {
		return RespError(c, "不支持的格式")
	}

	// 图片文件名
	mineToExts := map[string]string{
		"image/jpeg":    ".jpg",
		"image/png":     ".png",
		"image/gif":     ".gif",
		"image/svg+xml": ".svg",
		"image/webp":    ".webp",
		"image/bmp":     ".bmp",
	}

	t := time.Now()
	filename := t.Format("20060102-150405.000") + mineToExts[fileMine]

	// 创建图片目标文件
	if err := lib.EnsureDir(config.Instance.ImgUpload.Path); err != nil {
		logrus.Error(err)
		return RespError(c, "创建图片存放文件夹失败")
	}

	fileFullPath := strings.TrimSuffix(config.Instance.ImgUpload.Path, "/") + "/" + filename
	dst, err := os.Create(fileFullPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// 写入图片文件
	if _, err = dst.Write(buf); err != nil {
		return err
	}

	// 生成外部可访问链接
	baseURL := config.Instance.ImgUpload.PublicPath
	if baseURL == "" {
		baseURL = ImgUpload_RoutePath
	}

	// 响应数据
	return RespData(c, Map{
		"img_file": filename,
		"img_url":  path.Join(baseURL, filename),
	})
}
