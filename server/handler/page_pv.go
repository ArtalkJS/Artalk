package handler

import (
	"hash/maphash"
	"sync"
	"time"

	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/limiter"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsPagePV struct {
	PageKey   string `json:"page_key" validate:"required"`   // The page key
	PageTitle string `json:"page_title" validate:"optional"` // The page title
	SiteName  string `json:"site_name" validate:"optional"`  // The site name of your content scope
}

type ResponsePagePV struct {
	PV int `json:"pv"`
}

// @Id           LogPv
// @Summary      Increase Page Views (PV)
// @Description  Increase and get the number of page views
// @Tags         Page
// @Param        page  body  ParamsPagePV  true  "The page to record pv"
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponsePagePV
// @Router       /pages/pv  [post]
func PagePV(app *core.App, router fiber.Router) {
	// Bound both creation work and lock memory when callers supply arbitrary keys.
	creationLimit := limiter.NewRateLimiter(60, 600, time.Minute)
	var locks [64]sync.Mutex
	seed := maphash.MakeSeed()

	router.Post("/pages/pv", func(c *fiber.Ctx) error {
		var p ParamsPagePV
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		if _, ok, resp := common.CheckSiteExist(app, c, p.SiteName); !ok {
			return resp
		}

		mutex := &locks[maphash.String(seed, p.SiteName+"\x00"+p.PageKey)%uint64(len(locks))]
		mutex.Lock()
		defer mutex.Unlock()

		// Existing pages never consume creation quota, even after a limit is hit.
		page := app.Dao().FindPage(p.PageKey, p.SiteName)
		if page.IsEmpty() {
			if !creationLimit.Allow(c.IP()) {
				return common.RespData(c, ResponsePagePV{PV: 0})
			}
			page = app.Dao().FindCreatePage(p.PageKey, p.PageTitle, p.SiteName)
			if page.IsEmpty() {
				return common.RespError(c, 500, "Failed to create page")
			}
		}

		if err := app.Dao().IncrementPagePV(&page); err != nil {
			return common.RespError(c, 500, "Failed to update page views")
		}

		return common.RespData(c, ResponsePagePV{
			PV: page.PV,
		})
	})
}
