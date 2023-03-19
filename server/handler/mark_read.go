package handler

import (
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsMarkRead struct {
	CommentID uint   `form:"comment_id"`
	NotifyKey string `form:"notify_key"`

	Name    string `form:"name"`
	Email   string `form:"email"`
	AllRead bool   `form:"all_read"`

	SiteName string
	SiteID   uint
	SiteAll  bool
}

// @Summary      Notify Mark Read
// @Description  Mark specific notification as read for user
// @Tags         Notify
// @Param        comment_id  formData  int     true   "the comment id of the notify you want to mark as read"
// @Param        notify_key  formData  string  true   "the key of the notify"
// @Param        name        formData  string  false  "the username"
// @Param        email       formData  string  false  "the user email"
// @Param        all_read    formData  bool    false  "the option if mark all user's notify as read"
// @Param        site_name   formData  string  false  "the site name of your content scope"
// @Success      200  {object}  common.JSONResult
// @Router       /mark-read  [post]
func MarkRead(router fiber.Router) {
	router.Post("/mark-read", func(c *fiber.Ctx) error {
		var p ParamsMarkRead
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// use site
		common.UseSite(c, &p.SiteName, &p.SiteID, &p.SiteAll)

		// all read
		if p.AllRead {
			if p.Name == "" || p.Email == "" {
				return common.RespError(c, "username or email cannot be empty")
			}

			user := query.FindUser(p.Name, p.Email)
			err := query.UserNotifyMarkAllAsRead(user.ID)
			if err != nil {
				return common.RespError(c, err.Error())
			}

			return common.RespSuccess(c)
		}

		// find notify
		notify := query.FindNotifyForComment(p.CommentID, p.NotifyKey)
		if notify.IsEmpty() {
			return common.RespError(c, i18n.T("{{name}} not found", Map{"name": i18n.T("Notify")}))
		}

		if notify.IsRead {
			return common.RespSuccess(c)
		}

		// update notify
		err := query.NotifySetRead(&notify)
		if err != nil {
			return common.RespError(c, i18n.T("{{name}} save failed", Map{"name": i18n.T("Notify")}))
		}

		return common.RespSuccess(c)
	})
}
