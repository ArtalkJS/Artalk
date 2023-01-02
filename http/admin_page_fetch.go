package http

import (
	"fmt"

	"github.com/ArtalkJS/ArtalkGo/internal/config"
	"github.com/ArtalkJS/ArtalkGo/internal/entity"
	"github.com/ArtalkJS/ArtalkGo/internal/query"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ParamsAdminPageFetch struct {
	ID       uint `mapstructure:"id"`
	SiteName string

	GetStatus bool `mapstructure:"get_status"`
}

var allPageFetching = false
var allPageFetchDone = 0
var allPageFetchTotal = 0

func (a *action) AdminPageFetch(c echo.Context) error {
	var p ParamsAdminPageFetch
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	UseSite(c, &p.SiteName, nil, nil)

	// 状态获取
	if p.GetStatus {
		if allPageFetching {
			return RespData(c, Map{
				"msg":         fmt.Sprintf("已完成 %d 共 %d 个", allPageFetchDone, allPageFetchTotal),
				"is_progress": true,
			})
		} else {
			return RespData(c, Map{
				"msg":         "",
				"is_progress": false,
			})
		}
	}

	// 更新全部站点
	if p.SiteName != "" {
		if allPageFetching {
			return RespError(c, "任务正在进行中，请稍等片刻")
		}

		// 异步执行
		go func() {
			allPageFetching = true
			allPageFetchDone = 0
			allPageFetchTotal = 0
			var pages []entity.Page
			db := a.db.Model(&entity.Page{})
			if p.SiteName != config.ATK_SITE_ALL {
				db = db.Where(&entity.Page{SiteName: p.SiteName})
			}
			db.Find(&pages)

			allPageFetchTotal = len(pages)
			for _, p := range pages {
				if err := query.FetchPageFromURL(&p); err != nil {
					logrus.Error(c, "[api_admin_page_fetch] page fetch error: "+err.Error())
				} else {
					allPageFetchDone++
				}
			}
			allPageFetching = false
		}()

		return RespSuccess(c)
	}

	page := query.FindPageByID(p.ID)
	if page.IsEmpty() {
		return RespError(c, "page not found")
	}

	if !IsAdminHasSiteAccess(c, page.SiteName) {
		return RespError(c, "无权操作")
	}

	if err := query.FetchPageFromURL(&page); err != nil {
		return RespError(c, "page fetch error: "+err.Error())
	}

	return RespData(c, Map{
		"page": query.CookPage(&page),
	})
}
