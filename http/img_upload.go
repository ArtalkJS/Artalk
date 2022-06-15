package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
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

	SiteName string

	SiteID  uint
	SiteAll bool
}

func (a *action) ImgUpload(c echo.Context) error {
	// 功能开关 (管理员始终开启)
	if !config.Instance.ImgUpload.Enabled && !CheckIsAdminReq(c) {
		return RespError(c, "禁止图片上传", Map{
			"img_upload_enabled": false,
		})
	}

	// 传入参数解析
	var p ParamsImgUpload
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	if !lib.ValidateEmail(p.Email) {
		return RespError(c, "Invalid email")
	}

	// use site
	UseSite(c, &p.SiteName, &p.SiteID, &p.SiteAll)

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
		"image/jpeg", "image/png", "image/gif", "image/webp", "image/bmp",
		// "image/svg+xml",
	}
	if !lib.ContainsStr(allowMines, fileMine) {
		return RespError(c, "不支持的格式")
	}

	// 图片文件名
	mineToExts := map[string]string{
		"image/jpeg": ".jpg",
		"image/png":  ".png",
		"image/gif":  ".gif",
		"image/webp": ".webp",
		"image/bmp":  ".bmp",
		// "image/svg+xml": ".svg",
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
		logrus.Error(err)
		return RespError(c, "图片文件创建失败")
	}
	defer dst.Close()

	// 写入图片文件
	if _, err = dst.Write(buf); err != nil {
		logrus.Error(err)
		return RespError(c, "图片文件写入失败")
	}

	// 生成外部可访问链接
	baseURL := config.Instance.ImgUpload.PublicPath
	if baseURL == "" {
		baseURL = ImgUpload_RoutePath
	}
	imgURL := path.Join(baseURL, filename)

	// 使用 upgit
	if config.Instance.ImgUpload.Upgit.Enabled {
		upgitURL := execUpgitUpload(fileFullPath)
		if upgitURL == "" || !lib.ValidateURL(upgitURL) {
			// 上传失败，删除源图片文件
			var err = os.Remove(fileFullPath)
			if err != nil {
				logrus.Error(err)
			}

			logrus.Error("[IMG_UPLOAD] [upgit] upgit output: ", upgitURL)
			return RespError(c, "图片通过 upgit 上传失败")
		}

		// 上传成功，删除本地文件
		if config.Instance.ImgUpload.Upgit.DelLocal {
			var err = os.Remove(fileFullPath)
			if err != nil {
				logrus.Error(err)
			}
		}

		// 使用从 upgit 获取的图片 URL
		imgURL = upgitURL
	}

	// 响应数据
	return RespData(c, Map{
		"img_file": filename,
		"img_url":  imgURL,
	})
}

// 调用 upgit 上传图片获得 URL
func execUpgitUpload(filename string) string {
	LogTag := "[IMG_UPLOAD] [upgit] "

	// 处理参数
	cmdStrSplitted := strings.Split(config.Instance.ImgUpload.Upgit.Exec, " ")
	execApp := cmdStrSplitted[0]
	execArgs := []string{}
	for i, arg := range cmdStrSplitted {
		if i > 0 {
			execArgs = append(execArgs, arg)
		}
	}
	execArgs = append(execArgs, filename)

	// 执行命令
	cmd := exec.Command(execApp, execArgs...)
	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		logrus.Error(LogTag, "cmd.Start: ", err)
		return ""
	}

	result, _ := ioutil.ReadAll(stdout)
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				logrus.Error(LogTag, "Exit Status: ", status.ExitStatus())
			}
		} else {
			logrus.Error(LogTag, "cmd.Wait: ", err)
		}

		return ""
	}

	return strings.TrimSpace(string(result))
}
