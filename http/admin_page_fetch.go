package http

import (
	"fmt"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
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
			var pages []model.Page
			query := a.db.Model(&model.Page{})
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

	if !IsAdminHasSiteAccess(c, page.SiteName) {
		return RespError(c, "无权操作")
	}

	if err := page.FetchURL(); err != nil {
		return RespError(c, "page fetch error: "+err.Error())
	}

	return RespData(c, Map{
		"page": page.ToCooked(),
	})
}
