package http

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/araddon/dateparse"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ParamsAdminImporter struct {
	Type     string `mapstructure:"type" param:"required"`
	Data     string `mapstructure:"data" param:"required"`
	SiteName string `mapstructure:"site_name"`
	SiteID   uint
}

func ActionAdminImporter(c echo.Context) error {
	if isOK, resp := AdminOnly(c); !isOK {
		return resp
	}

	var p ParamsAdminImporter
	if isOK, resp := ParamsDecode(c, ParamsAdminImporter{}, &p); !isOK {
		return resp
	}

	// check create site
	if model.FindSite(p.SiteName).IsEmpty() {
		site := model.Site{}
		site.Name = p.SiteName
		err := lib.DB.Create(&site).Error
		if err != nil {
			return RespError(c, "站点创建失败")
		}
	}

	// find site
	if isOK, resp := CheckSite(c, &p.SiteName, &p.SiteID, nil); !isOK {
		return resp
	}

	// 数据导入注意事项：
	// 1. 记得恢复创建/更新时间
	// 2. 因为是数据追加模式，ID 会发生变动，需更新 rid

	errs := []string{}
	dbResultLog := func(name string, err error) {
		if err != nil {
			logrus.Error("Save Comment error: ", err)
			errs = append(errs, fmt.Sprintf("%s: %s", name, err.Error()))
		}
	}

	idChanges := map[uint]uint{} // comment_id: old => new

	switch p.Type {
	// Artalk v1 数据源
	case "artalk_v1":
		// 解析数据源
		var v1Comments []ArtalkV1_Comment
		err := json.Unmarshal([]byte(p.Data), &v1Comments)
		if err != nil {
			return RespError(c, "json unmarshal err:"+err.Error())
		}

		for _, oc := range v1Comments {
			if strings.TrimSpace(oc.Content) == "" {
				continue
			}

			// 创建 page 和 user
			page := model.FindCreatePage(oc.PageKey, "", p.SiteName)
			user := model.FindCreateUser(oc.Nick, oc.Email)

			user.Link = oc.Link
			user.LastIP = oc.IP
			user.LastUA = oc.UA
			model.UpdateUser(&user)

			// 创建新 comment 实例
			nComment := model.Comment{
				Content:     oc.Content,
				Rid:         oc.Rid,
				UserID:      user.ID,
				PageKey:     page.Key,
				IP:          oc.IP,
				UA:          oc.UA,
				IsPending:   oc.IsPending,
				IsCollapsed: oc.IsCollapsed,
				SiteName:    p.SiteName,
			}

			// 日期恢复
			denverLoc, _ := time.LoadLocation(config.Instance.TimeZone) // 时区
			time.Local = denverLoc
			t, _ := dateparse.ParseIn(oc.Date, denverLoc)
			nComment.CreatedAt = t
			nComment.UpdatedAt = t

			// 保存到数据库
			err := lib.DB.Create(&nComment).Error
			dbResultLog(fmt.Sprintf("保存 old_id:%d new_id:%d", oc.ID, nComment.ID), err)

			idChanges[oc.ID] = nComment.ID
		}

		// reply id 重建
		for _, oc := range v1Comments {
			if oc.Rid != 0 {
				if _, isExist := idChanges[oc.ID]; !isExist {
					continue
				}

				nComment := model.FindComment(idChanges[oc.ID], p.SiteName)
				nComment.Rid = idChanges[oc.Rid]
				err := lib.DB.Save(&nComment).Error
				dbResultLog(fmt.Sprintf("rid 更新 new_id:%d new_rid:%d", nComment.ID, idChanges[oc.Rid]), err)
			}
		}
	default:
		return RespError(c, "invalid type")
	}

	if len(errs) != 0 {
		return RespError(c, "导入发生错误："+strings.Join(errs, "\n"))
	}
	return RespSuccess(c)
}

type ArtalkV1_Comment struct {
	ID          uint   `json:"id"`
	Content     string `json:"content"`
	Nick        string `json:"nick"`
	Email       string `json:"email"`
	Link        string `json:"link"`
	UA          string `json:"ua"`
	PageKey     string `json:"page_key"`
	Rid         uint   `json:"rid"`
	IP          string `json:"ip"`
	Date        string `json:"date"`
	IsPending   bool   `json:"is_pending"`
	IsCollapsed bool   `json:"is_collapsed"`
}
