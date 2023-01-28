package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ParamsImgUpload struct {
	Name  string `form:"name" validate:"required"`
	Email string `form:"email" validate:"required"`

	PageKey   string `form:"page_key" validate:"required"`
	PageTitle string `form:"page_title"`

	SiteName string

	SiteID  uint
	SiteAll bool
}

type ResponseImgUpload struct {
	ImgFile string `json:"img_file"`
	ImgURL  string `json:"img_url"`
}

// @Summary      Image Upload
// @Description  Upload image from this endpoint
// @Tags         Upload
// @Param        file           formData  file    true   "upload file in preparation for import"
// @Param        name           formData  string  true   "the username"
// @Param        email          formData  string  true   "the user email"
// @Param        page_key       formData  string  true   "the page key"
// @Param        page_title     formData  string  false  "the page title"
// @Param        site_name      formData  string  false  "the site name of your content scope"
// @Security     ApiKeyAuth
// @Success      200  {object}  common.JSONResult{data=ResponseImgUpload}
// @Router       /img-upload  [post]
func ImgUpload(router fiber.Router) {
	router.Post("/img-upload", func(c *fiber.Ctx) error {
		// 功能开关 (管理员始终开启)
		if !config.Instance.ImgUpload.Enabled && !common.CheckIsAdminReq(c) {
			return common.RespError(c, i18n.T("Image upload forbidden"), common.Map{
				"img_upload_enabled": false,
			})
		}

		// 传入参数解析
		var p ParamsImgUpload
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		if !utils.ValidateEmail(p.Email) {
			return common.RespError(c, i18n.T("Invalid {{name}}", Map{"name": i18n.T("Email")}))
		}

		// use site
		common.UseSite(c, &p.SiteName, &p.SiteID, &p.SiteAll)

		// 记录请求次数 (for 请求频率限制)
		common.RecordAction(c)

		// find page
		// page := entity.FindPage(p.PageKey, p.PageTitle)
		// ip := c.RealIP()
		// ua := c.Request().UserAgent()

		// 图片大小限制 (Based on content length)
		if config.Instance.ImgUpload.MaxSize != 0 {
			if int64(c.Request().Header.ContentLength()) > config.Instance.ImgUpload.MaxSize*1024*1024 {
				return common.RespError(c, i18n.T("Image exceeds {{file_size}} limit", Map{
					"file_size": fmt.Sprintf("%dMB", config.Instance.ImgUpload.MaxSize),
				}))
			}
		}

		// 获取 Form
		file, err := c.FormFile("file")
		if err != nil {
			logrus.Error(err)
			return common.RespError(c, "File read failed")
		}

		// 打开文件
		src, err := file.Open()
		if err != nil {
			logrus.Error(err)
			return common.RespError(c, "File open failed")
		}
		defer src.Close()

		// 读取文件
		buf, err := io.ReadAll(src)
		if err != nil {
			logrus.Error(err)
			return common.RespError(c, "File read failed")
		}

		// 大小限制 (Based on content read)
		if config.Instance.ImgUpload.MaxSize != 0 {
			if int64(len(buf)) > config.Instance.ImgUpload.MaxSize*1024*1024 {
				return common.RespError(c, i18n.T("Image exceeds {{file_size}} limit", Map{
					"file_size": fmt.Sprintf("%dMB", config.Instance.ImgUpload.MaxSize),
				}))
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
		if !utils.ContainsStr(allowMines, fileMine) {
			return common.RespError(c, i18n.T("Unsupported formats"))
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
		if err := utils.EnsureDir(config.Instance.ImgUpload.Path); err != nil {
			logrus.Error(err)
			return common.RespError(c, "Folder creation failed")
		}

		fileFullPath := strings.TrimSuffix(config.Instance.ImgUpload.Path, "/") + "/" + filename
		dst, err := os.Create(fileFullPath)
		if err != nil {
			logrus.Error(err)
			return common.RespError(c, "File creation failed")
		}
		defer dst.Close()

		// 写入图片文件
		if _, err = dst.Write(buf); err != nil {
			logrus.Error(err)
			return common.RespError(c, "File write failed")
		}

		// 生成外部可访问链接
		baseURL := config.Instance.ImgUpload.PublicPath
		if baseURL == "" {
			baseURL = config.IMG_UPLOAD_PUBLIC_PATH
		}
		imgURL := path.Join(baseURL, filename)

		// 使用 upgit
		if config.Instance.ImgUpload.Upgit.Enabled {
			upgitURL := execUpgitUpload(fileFullPath)
			if upgitURL == "" || !utils.ValidateURL(upgitURL) {
				// 上传失败，删除源图片文件
				var err = os.Remove(fileFullPath)
				if err != nil {
					logrus.Error(err)
				}

				logrus.Error("[IMG_UPLOAD] [upgit] upgit output: ", upgitURL)
				return common.RespError(c, i18n.T("Upload image via {{method}} failed", Map{"method": "upgit"}))
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
		return common.RespData(c, ResponseImgUpload{
			ImgFile: filename,
			ImgURL:  imgURL,
		})
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

	result, _ := io.ReadAll(stdout)
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
