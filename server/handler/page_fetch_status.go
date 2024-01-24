package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

var (
	allPageFetching   = false
	allPageFetchDone  = 0
	allPageFetchTotal = 0
)

type ResponsePageFetchStatus struct {
	Msg        string `json:"msg"`         // The message of the task status
	IsProgress bool   `json:"is_progress"` // If the task is in progress
	Done       int    `json:"done"`        // The number of pages that have been fetched
	Total      int    `json:"total"`       // The total number of pages
}

// @Id           GetPageFetchStatus
// @Summary      Get Pages Fetch Status
// @Description  Get the status of the task of fetching all pages
// @Tags         Page
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  ResponsePageFetchStatus
// @Router       /pages/fetch/status  [get]
func PageFetchStatus(app *core.App, router fiber.Router) {
	router.Get("/pages/fetch/status", common.AdminGuard(app, func(c *fiber.Ctx) error {
		return common.RespData(c, ResponsePageFetchStatus{
			Msg:        i18n.T("{{done}} of {{total}} done", Map{"done": allPageFetchDone, "total": allPageFetchTotal}),
			IsProgress: allPageFetching,
			Done:       allPageFetchDone,
			Total:      allPageFetchTotal,
		})
	}))
}
