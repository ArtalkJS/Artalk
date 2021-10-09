package http

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

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

	// find site
	if isOK, resp := CheckSite(c, p.SiteName, &p.SiteID); !isOK {
		return resp
	}

	switch p.Type {
	case "artalk_v1":
		var comments []ArtalkV1_Comment
		err := json.Unmarshal([]byte(p.Data), &comments)
		if err != nil {
			return RespError(c, "json unmarshal err:"+err.Error())
		}

		errs := []string{}
		for _, c := range comments {
			page := model.FindCreatePage(c.PageKey, "", "", p.SiteName)
			user := model.FindCreateUser(c.Nick, c.Email)

			user.Link = c.Link
			user.LastIP = c.IP
			user.LastUA = c.UA
			model.UpdateUser(&user)

			nComment := model.Comment{
				Content:     c.Content,
				Rid:         c.Rid,
				UserID:      user.ID,
				PageKey:     page.Key,
				IP:          c.IP,
				UA:          c.UA,
				IsPending:   c.IsPending,
				IsCollapsed: c.IsCollapsed,
				SiteName:    p.SiteName,
			}

			err := lib.DB.Create(&nComment).Error
			if err != nil {
				logrus.Error("Save Comment error: ", err)
				errs = append(errs, fmt.Sprintf("%d: %s", c.ID, err.Error()))
			}
		}

		if len(errs) != 0 {
			return RespError(c, "导入发生错误："+strings.Join(errs, "\n"))
		}
		return RespSuccess(c)
	}

	return RespError(c, "invalid type")
}
