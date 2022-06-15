package http

import (
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsMarkRead struct {
	NotifyKey string `mapstructure:"notify_key"`

	Name    string `mapstructure:"name"`
	Email   string `mapstructure:"email"`
	AllRead bool   `mapstructure:"all_read"`

	SiteName string
	SiteID   uint
	SiteAll  bool
}

func (a *action) MarkRead(c echo.Context) error {
	var p ParamsMarkRead
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	// use site
	UseSite(c, &p.SiteName, &p.SiteID, &p.SiteAll)

	// all read
	if p.AllRead {
		if p.Name == "" || p.Email == "" {
			return RespError(c, "need username and email")
		}

		user := model.FindUser(p.Name, p.Email)
		err := model.UserNotifyMarkAllAsRead(user.ID)
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
