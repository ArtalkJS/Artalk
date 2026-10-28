package handler

import (
	"sync"

	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/i18n"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/gofiber/fiber/v2"
)

var allPageFetch pageFetchState

// Keep task state and its lock together. No lock is held during network work.
type pageFetchState struct {
	mu          sync.Mutex
	running     bool
	done, total int
}

func (s *pageFetchState) start() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running {
		return false
	}
	s.running, s.done, s.total = true, 0, 0
	return true
}

func (s *pageFetchState) finish() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.running = false
}

func (s *pageFetchState) progress(done, total int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.done, s.total = done, total
}

func (s *pageFetchState) snapshot() (running bool, done, total int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running, s.done, s.total
}

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
		progress, done, total := allPageFetch.snapshot()
		return common.RespData(c, ResponsePageFetchStatus{
			Msg:        i18n.T("{{done}} of {{total}} done", Map{"done": done, "total": total}),
			IsProgress: progress,
			Done:       done,
			Total:      total,
		})
	}))
}
