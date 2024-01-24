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
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsUpload struct {
}

type ResponseUpload struct {
	FileType  string `json:"file_type" enum:"image"`
	FileName  string `json:"file_name"`
	PublicURL string `json:"public_url"`
}

// @Id           Upload
// @Summary      Upload
// @Description  Upload file from this endpoint
// @Tags         Upload
// @Param        file           formData  file    true   "Upload file"
// @Security     ApiKeyAuth
// @Accept       mpfd
// @Produce      json
// @Success      200  {object}  ResponseUpload
// @Failure      400  {object}  Map{msg=string}
// @Failure      403  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /upload  [post]
func Upload(app *core.App, router fiber.Router) {
	router.Post("/upload", common.LimiterGuard(app, func(c *fiber.Ctx) error {
		// 功能开关 (管理员始终开启)
		if !app.Conf().ImgUpload.Enabled && !common.CheckIsAdminReq(app, c) {
			return common.RespError(c, 403, i18n.T("Image upload forbidden"), common.Map{
				"img_upload_enabled": false,
			})
		}

		// 传入参数解析
		var p ParamsUpload
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// find page
		// page := entity.FindPage(p.PageKey, p.PageTitle)
		// ip := c.RealIP()
		// ua := c.Request().UserAgent()

		// 图片大小限制 (Based on content length)
		if app.Conf().ImgUpload.MaxSize != 0 {
			if int64(c.Request().Header.ContentLength()) > app.Conf().ImgUpload.MaxSize*1024*1024 {
				return common.RespError(c, 400, i18n.T("Image exceeds {{file_size}} limit", Map{
					"file_size": fmt.Sprintf("%dMB", app.Conf().ImgUpload.MaxSize),
				}))
			}
		}

		// 获取 Form
		file, err := c.FormFile("file")
		if err != nil {
			log.Error(err)
			return common.RespError(c, 500, "File read failed")
		}

		// 打开文件
		src, err := file.Open()
		if err != nil {
			log.Error(err)
			return common.RespError(c, 500, "File open failed")
		}
		defer src.Close()

		// 读取文件
		buf, err := io.ReadAll(src)
		if err != nil {
			log.Error(err)
			return common.RespError(c, 500, "File read failed")
		}

		// 大小限制 (Based on content read)
		if app.Conf().ImgUpload.MaxSize != 0 {
			if int64(len(buf)) > app.Conf().ImgUpload.MaxSize*1024*1024 {
				return common.RespError(c, 400, i18n.T("Image exceeds {{file_size}} limit", Map{
					"file_size": fmt.Sprintf("%dMB", app.Conf().ImgUpload.MaxSize),
				}))
			}
		}

		// 文件格式判断
		// The http.DetectContentType function reads the first 512 bytes of a file
		// and uses these bytes (aka the magic number) to determine the file's content type.
		// @link https://mimesniff.spec.whatwg.org/
		// @link https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Basics_of_HTTP/MIME_types
		fileMine := http.DetectContentType(buf)
		allowMines := []string{
			"image/jpeg", "image/png", "image/gif", "image/webp", "image/bmp",
			// "image/svg+xml",
		}
		if !utils.ContainsStr(allowMines, fileMine) {
			return common.RespError(c, 400, i18n.T("Unsupported formats"))
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
		if err := utils.EnsureDir(app.Conf().ImgUpload.Path); err != nil {
			log.Error(err)
			return common.RespError(c, 500, "Folder creation failed")
		}

		fileFullPath := strings.TrimSuffix(app.Conf().ImgUpload.Path, "/") + "/" + filename
		dst, err := os.Create(fileFullPath)
		if err != nil {
			log.Error(err)
			return common.RespError(c, 500, "File creation failed")
		}
		defer dst.Close()

		// 写入图片文件
		if _, err = dst.Write(buf); err != nil {
			log.Error(err)
			return common.RespError(c, 500, "File write failed")
		}

		// 生成外部可访问链接
		baseURL := app.Conf().ImgUpload.PublicPath
		if baseURL == "" {
			baseURL = config.IMG_UPLOAD_PUBLIC_PATH
		}

		var imgURL string
		if utils.ValidateURL(baseURL) {
			// full url
			imgURL = strings.TrimSuffix(baseURL, "/") + "/" + filename
		} else {
			// relative path
			imgURL = path.Join(baseURL, filename)
		}

		// 使用 upgit
		if app.Conf().ImgUpload.Upgit.Enabled {
			upgitURL := execUpgitUpload(app.Conf().ImgUpload.Upgit.Exec, fileFullPath)
			if upgitURL == "" || !utils.ValidateURL(upgitURL) {
				// 上传失败，删除源图片文件
				var err = os.Remove(fileFullPath)
				if err != nil {
					log.Error(err)
				}

				log.Error("[IMG_UPLOAD] [upgit] upgit output: ", upgitURL)
				return common.RespError(c, 500, i18n.T("Upload image via {{method}} failed", Map{"method": "upgit"}))
			}

			// 上传成功，删除本地文件
			if app.Conf().ImgUpload.Upgit.DelLocal {
				var err = os.Remove(fileFullPath)
				if err != nil {
					log.Error(err)
				}
			}

			// 使用从 upgit 获取的图片 URL
			imgURL = upgitURL
		}

		// 响应数据
		return common.RespData(c, ResponseUpload{
			FileType:  "image",
			FileName:  filename,
			PublicURL: imgURL,
		})
	}))
}

// 调用 upgit 上传图片获得 URL
func execUpgitUpload(execCommand string, filename string) string {
	LogTag := "[IMG_UPLOAD] [upgit] "

	// 处理参数
	cmdStrSplitted := strings.Split(execCommand, " ")
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
		log.Error(LogTag, "cmd.Start: ", err)
		return ""
	}

	result, _ := io.ReadAll(stdout)
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				log.Error(LogTag, "Exit Status: ", status.ExitStatus())
			}
		} else {
			log.Error(LogTag, "cmd.Wait: ", err)
		}

		return ""
	}

	return strings.TrimSpace(string(result))
}
