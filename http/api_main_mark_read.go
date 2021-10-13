package http

import (
	"errors"
	"time"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsMarkRead struct {
	NotifyKey string `mapstructure:"notify_key"`

	Name    string `mapstructure:"name"`
	Email   string `mapstructure:"email"`
	AllRead bool   `mapstructure:"all_read"`

	SiteName string `mapstructure:"site_name"`
	SiteID   uint
	SiteAll  bool
}

func ActionMarkRead(c echo.Context) error {
	var p ParamsMarkRead
	if isOK, resp := ParamsDecode(c, ParamsMarkRead{}, &p); !isOK {
		return resp
	}

	// find site
	if isOK, resp := CheckSite(c, &p.SiteName, &p.SiteID, &p.SiteAll); !isOK {
		return resp
	}

	// all read
	if p.AllRead {
		if p.Name == "" || p.Email == "" {
			return RespError(c, "need username and email")
		}

		user := model.FindUser(p.Name, p.Email)
		err := NotifyMarkAllAsRead(user.ID)
		if err != nil {
			return RespError(c, err.Error())
		}

		return RespSuccess(c)
	}

	// find notify
	notify := model.FindNotifyByKey(p.NotifyKey)
	if notify.IsEmpty() {
		return RespError(c, "notify key is wrong")
	}

	if notify.IsRead {
		return RespSuccess(c)
	}

	// update notify
	err := notify.SetRead()
	if err != nil {
		return RespError(c, "notify save error")
	}

	return RespSuccess(c)
}

func NotifyMarkAllAsRead(userID uint) error {
	if userID == 0 {
		return errors.New("user not found")
	}

	lib.DB.Model(&model.Notify{}).Where("user_id = ?", userID).Updates(&model.Notify{
		IsRead: true,
		ReadAt: time.Now(),
	})

	return nil
}
