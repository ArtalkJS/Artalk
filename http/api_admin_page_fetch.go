package http

import (
	"fmt"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ParamsAdminPageFetch struct {
	ID       uint   `mapstructure:"id"`
	SiteName string `mapstructure:"site_name"`
}

var allPageFetching = false
var allPageFetchDone = 0
var allPageFetchTotal = 0

func ActionAdminPageFetch(c echo.Context) error {
	if isOK, resp := AdminOnly(c); !isOK {
		return resp
	}

	var p ParamsAdminPageFetch
	if isOK, resp := ParamsDecode(c, ParamsAdminPageFetch{}, &p); !isOK {
		return resp
	}

	// 更新全部站点
	if p.SiteName != "" {
		if allPageFetching {
			return RespError(c, fmt.Sprintf("页面更新正在执行中，已完成 %d 共 %d 个", allPageFetchDone, allPageFetchTotal))
		}

		// 异步执行
		go func() {
			allPageFetching = true
			allPageFetchDone = 0
			allPageFetchTotal = 0
			var pages []model.Page
			query := lib.DB
			if p.SiteName != lib.ATK_SITE_ALL {
				query = query.Where(&model.Page{SiteName: p.SiteName})
			}
			query.Find(&pages)

			allPageFetchTotal = len(pages)
			for _, p := range pages {
				if err := p.FetchURL(); err != nil {
					logrus.Error(c, "[api_admin_page_fetch] page fetch error: "+err.Error())
				} else {
					allPageFetchDone++
				}
			}
			allPageFetching = false
		}()

		return RespSuccess(c)
	}

	page := model.FindPageByID(p.ID)
	if page.IsEmpty() {
		return RespError(c, "page not found")
	}

	if err := page.FetchURL(); err != nil {
		return RespError(c, "page fetch error: "+err.Error())
	}

	return RespData(c, Map{
		"page": page.ToCooked(),
	})
}
