package http

import (
	"strconv"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/lib/email"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsCommentEdit struct {
	// 查询值
	ID       uint `mapstructure:"id" param:"required"`
	SiteName string
	SiteID   uint
	SiteAll  bool

	// 可修改
	Content     string `mapstructure:"content"`
	PageKey     string `mapstructure:"page_key"`
	Nick        string `mapstructure:"nick"`
	Email       string `mapstructure:"email"`
	Link        string `mapstructure:"link"`
	Rid         string `mapstructure:"rid"`
	UA          string `mapstructure:"ua"`
	IP          string `mapstructure:"ip"`
	IsCollapsed bool   `mapstructure:"is_collapsed"`
	IsPending   bool   `mapstructure:"is_pending"`
	IsPinned    bool   `mapstructure:"is_pinned"`
}

func (a *action) AdminCommentEdit(c echo.Context) error {
	var p ParamsCommentEdit
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	// use site
	UseSite(c, &p.SiteName, &p.SiteID, &p.SiteAll)

	// find comment
	comment := model.FindComment(p.ID)
	if comment.IsEmpty() {
		return RespError(c, "comment not found")
	}

	if !IsAdminHasSiteAccess(c, comment.SiteName) {
		return RespError(c, "无权操作")
	}

	// check params
	if p.Email != "" && !lib.ValidateEmail(p.Email) {
		return RespError(c, "Invalid email")
	}
	if p.Link != "" && !lib.ValidateURL(p.Link) {
		return RespError(c, "Invalid link")
	}

	// content
	if p.Content != "" {
		comment.Content = p.Content
	}

	// rid
	if p.Rid != "" {
		if rid, err := strconv.Atoi(p.Rid); err == nil {
			comment.Rid = uint(rid)
		}
	}

	// merge user
	originalUser := comment.FetchUser()
	if p.Nick == "" {
		p.Nick = originalUser.Name
	}
	if p.Email == "" {
		p.Email = originalUser.Email
	}

	// find or save new user
	user := model.FindCreateUser(p.Nick, p.Email, p.Link)
	if user.ID != comment.UserID {
		comment.UserID = user.ID
	}

	// user link update
	if p.Link != "" && p.Link != user.Link {
		user.Link = p.Link
		model.UpdateUser(&user)
	}

	// pageKey
	if p.PageKey != "" && p.PageKey != comment.PageKey {
		model.FindCreatePage(p.PageKey, "", p.SiteName)
		comment.PageKey = p.PageKey
	}

	comment.UA = p.UA
	comment.IP = p.IP
	comment.IsCollapsed = p.IsCollapsed
	comment.IsPinned = p.IsPinned

	if p.IsPending != comment.IsPending {
		// 待审状态发生改变
		comment.IsPending = p.IsPending

		// 待审状态被修改为 false，则重新发送邮件通知
		if !comment.IsPending {
			RenotifyWhenPendingModified(&comment)
		}
	}

	if err := model.UpdateComment(&comment); err != nil {
		return RespError(c, "comment save error")
	}

	return RespData(c, Map{
		"comment": comment.ToCooked(),
	})
}

func RenotifyWhenPendingModified(comment *model.Comment) {
	if comment.Rid == 0 {
		return // Root 评论不发送通知，因为这个评论已经被管理员看到了
	}

	pComment := model.FindComment(comment.Rid)
	if pComment.FetchUser().IsAdmin {
		return // 回复对象是管理员，则不再发送通知，因为已经看到了
	}

	if comment.UserID == pComment.UserID {
		return // 自己回复自己，不通知
	}

	notify := model.FindCreateNotify(pComment.UserID, comment.ID)
	if notify.IsEmailed {
		return // 邮件已经发送过，则不再重复发送
	}

	notify.SetComment(*comment)
	notify.SetInitial()

	// 邮件通知
	email.AsyncSend(&notify)
}
