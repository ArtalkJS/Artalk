package handler

import (
	"io"
	"os"

	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ResponseImportUpload struct {
	// The uploaded file name which can be used to import
	Filename string `json:"filename"`
}

// @Summary      Upload Artrans
// @Description  Upload a file to prepare to import
// @Tags         Transfer
// @Security     ApiKeyAuth
// @Param        file  formData  file  true  "Upload file in preparation for import task"
// @Accept       mpfd
// @Produce      json
// @Success      200  {object}  ResponseImportUpload{filename=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /transfer/upload  [post]
func transferUpload(app *core.App) func(c *fiber.Ctx) error {
	return common.AdminGuard(app, func(c *fiber.Ctx) error {
		// Get file from FormData
		file, err := c.FormFile("file")
		if err != nil {
			log.Error(err)
			return common.RespError(c, 500, "File read failed")
		}

		// Open file
		src, err := file.Open()
		if err != nil {
			log.Error(err)
			return common.RespError(c, 500, "File open failed")
		}
		defer src.Close()

		// Read file to buffer
		buf, err := io.ReadAll(src)
		if err != nil {
			log.Error(err)
			return common.RespError(c, 500, "File read failed")
		}

		// Create temp file
		tmpFile, err := os.CreateTemp("", "artalk-import-file-")
		if err != nil {
			log.Error(err)
			return common.RespError(c, 500, "tmp file creation failed")
		}

		// Write buffer to temp file
		tmpFile.Write(buf)

		// Return filename
		return common.RespData(c, ResponseImportUpload{
			Filename: tmpFile.Name(),
		})
	})
}
