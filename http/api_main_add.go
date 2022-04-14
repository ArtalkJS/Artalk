package http

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/lib/email"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
	"github.com/nikoksr/notify"
	"github.com/sirupsen/logrus"
)

type ParamsAdd struct {
	Name    string `mapstructure:"name" param:"required"`
	Email   string `mapstructure:"email" param:"required"`
	Link    string `mapstructure:"link"`
	Content string `mapstructure:"content" param:"required"`
	Rid     uint   `mapstructure:"rid"`
	UA      string `mapstructure:"ua"`

	PageKey   string `mapstructure:"page_key" param:"required"`
	PageTitle string `mapstructure:"page_title"`

	Token    string `mapstructure:"token"`
	SiteName string `mapstructure:"site_name"`
	SiteID   uint
}

type ResponseAdd struct {
	Comment model.CookedComment `json:"comment"`
}

func ActionAdd(c echo.Context) error {
	var p ParamsAdd
	if isOK, resp := ParamsDecode(c, ParamsAdd{}, &p); !isOK {
		return resp
	}

	if !lib.ValidateEmail(p.Email) {
		return RespError(c, "Invalid email")
	}
	if p.Link != "" && !lib.ValidateURL(p.Link) {
		return RespError(c, "Invalid link")
	}

	ip := c.RealIP()
	ua := c.Request().UserAgent()

	// 仅允许针对 Win11 的 UA 修正
	if p.UA != "" {
		if matchWin11, _ := regexp.MatchString(`Windows\W+NT\W+11.0`, p.UA); matchWin11 {
			ua = p.UA
		}
	}

	// record action for limiting action
	RecordAction(c)

	// find site
	if isOK, resp := CheckSite(c, &p.SiteName, &p.SiteID, nil); !isOK {
		return resp
	}

	// find page
	page := model.FindCreatePage(p.PageKey, p.PageTitle, p.SiteName)

	// check if the user is allowed to comment
	if isAllowed, resp := CheckIfAllowed(c, p.Name, p.Email, page, p.SiteName); !isAllowed {
		return resp
	}

	// check reply comment
	var parentComment model.Comment
	if p.Rid != 0 {
		parentComment = model.FindComment(p.Rid)
		if parentComment.IsEmpty() {
			return RespError(c, "找不到父评论")
		}
		if parentComment.PageKey != p.PageKey {
			return RespError(c, "与父评论的 pageKey 不一致")
		}
		if !parentComment.IsAllowReply() {
			return RespError(c, "不允许回复该评论")
		}
	}

	// find user
	user := model.FindCreateUser(p.Name, p.Email, p.Link)
	if user.ID == 0 || page.Key == "" {
		logrus.Error("Cannot get user or page")
		return RespError(c, "评论失败")
	}

	// update user
	user.Link = p.Link
	user.LastIP = ip
	user.LastUA = ua
	user.Name = p.Name // for 若用户修改用户名大小写
	user.Email = p.Email
	model.UpdateUser(&user)

	comment := model.Comment{
		Content:  p.Content,
		PageKey:  page.Key,
		SiteName: p.SiteName,

		UserID: user.ID,
		IP:     ip,
		UA:     ua,

		Rid: p.Rid,

		IsPending:   false,
		IsCollapsed: false,
		IsPinned:    false,

		User: user,
		Page: page,
	}

	// default comment type
	if config.Instance.Moderator.PendingDefault {
		comment.IsPending = true
	}

	// save to database
	err := lib.DB.Create(&comment).Error
	if err != nil {
		logrus.Error("Save Comment error: ", err)
		return RespError(c, "评论失败")
	}

	// 异步执行
	go func() {
		// update page
		if page.ToCooked().URL != "" && page.Title == "" {
			page.FetchURL()
		}

		// 垃圾检测
		if !CheckIsAdminReq(c) {
			comment.SpamCheck(c)
		}

		// 邮件通知发送
		EmailSend(&comment, &parentComment)

		// 其他通知
		UseNotify(&comment, &parentComment)
	}()

	return RespData(c, ResponseAdd{
		Comment: comment.ToCooked(),
	})
}

// 邮件发送 (from comment to parentComment)
func EmailSend(comment *model.Comment, parentComment *model.Comment) {
	if !config.Instance.Email.Enabled {
		return
	}

	// 自己回复自己，不提醒
	if comment.UserID == parentComment.UserID {
		return
	}

	// 邮件回复对方
	if !parentComment.IsEmpty() && !comment.IsPending {
		notify := model.FindCreateNotify(parentComment.UserID, comment.ID)
		notify.Comment = *comment
		email.AsyncSend(&notify)
	}

	// 邮件通知管理员
	admins := model.GetAllAdmins()
	userIsAdmin := func(userID uint) bool {
		for _, admin := range admins {
			if admin.ID == userID {
				return true
			}
		}
		return false
	}

	if parentComment.IsEmpty() && len(admins) > 0 {
		// TODO: 增加用户的站点隔离，指定管理员分配网站
		for _, admin := range admins {
			// 管理员自己回复自己，不提醒
			if comment.UserID == admin.ID {
				continue
			}

			// 管理员评论不回复给其他管理员
			if userIsAdmin(comment.UserID) {
				continue
			}

			notify := model.FindCreateNotify(admin.ID, comment.ID)
			notify.Comment = *comment
			email.AsyncSendToAdmin(&notify, &admin) // 发送邮件给管理员
		}
	}
}

var NotifyCtx = context.Background()

// 其他通知方式
func UseNotify(comment *model.Comment, parentComment *model.Comment) {
	// 忽略管理员回复
	coUser := comment.FetchUser()
	if coUser.IsAdmin {
		return
	}

	// 如果不是 root 评论，并且回复的不是管理员，直接忽略
	if !parentComment.IsEmpty() && !parentComment.FetchUser().IsAdmin {
		return
	}

	// 评论内容
	coContent := comment.Content
	if len(coContent) > 280 {
		coContent = string([]rune(coContent)[0:280]) + "..." // 截取文字
	}

	// 消息内容
	title := fmt.Sprintf("来自 @%s 的回复", coUser.Name)
	msg := fmt.Sprintf("%s\n\n%s", coContent, comment.GetLinkToReply())

	// 使用 notify 发送
	go func() {
		_ = notify.Send(NotifyCtx, title, msg)
	}()

	// 飞书
	go func() {
		sendLark(title + "\n\n" + msg)
	}()

	// Bark
	go func() {
		sendBark(title, msg)
	}()
}

func sendLark(msg string) {
	larkConf := config.Instance.Notify.Lark
	if !larkConf.Enabled {
		return
	}

	sendData := fmt.Sprintf(`{"msg_type":"text","content":{"text":%s}}`, strconv.Quote(msg))
	result, err := http.Post(larkConf.WebhookURL, "application/json", strings.NewReader(sendData))
	if err != nil {
		logrus.Error("[飞书]", " 消息发送失败：", err)
		return
	}

	defer result.Body.Close()
}

func sendBark(title string, msg string) {
	barkConf := config.Instance.Notify.Bark
	if !barkConf.Enabled {
		return
	}

	result, err := http.Get(fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(barkConf.Server, "/"), url.QueryEscape(title), url.QueryEscape(msg)))
	if err != nil {
		logrus.Error("[Bark]", " 消息发送失败：", err)
		return
	}

	defer result.Body.Close()
}
