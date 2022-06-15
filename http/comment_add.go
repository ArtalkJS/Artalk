package http

import (
	"strings"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/lib/anti_spam"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/ArtalkJS/ArtalkGo/model/notify_launcher"
	"github.com/labstack/echo/v4"

	"github.com/sirupsen/logrus"
)

type ParamsAdd struct {
	Name    string `mapstructure:"name"`
	Email   string `mapstructure:"email"`
	Link    string `mapstructure:"link"`
	Content string `mapstructure:"content" param:"required"`
	Rid     uint   `mapstructure:"rid"`
	UA      string `mapstructure:"ua"`

	PageKey   string `mapstructure:"page_key" param:"required"`
	PageTitle string `mapstructure:"page_title"`

	Token    string `mapstructure:"token"`
	SiteName string
	SiteID   uint
}

type ResponseAdd struct {
	Comment model.CookedComment `json:"comment"`
}

func (a *action) Add(c echo.Context) error {
	var p ParamsAdd
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	if strings.TrimSpace(p.Name) == "" {
		return RespError(c, "昵称不能为空")
	}
	if strings.TrimSpace(p.Email) == "" {
		return RespError(c, "邮箱不能为空")
	}

	if !lib.ValidateEmail(p.Email) {
		return RespError(c, "Invalid email")
	}
	if p.Link != "" && !lib.ValidateURL(p.Link) {
		return RespError(c, "Invalid link")
	}

	ip := c.RealIP()
	ua := c.Request().UserAgent()

	// 允许传入修正后的 UA
	if p.UA != "" {
		ua = p.UA
	}

	// record action for limiting action
	RecordAction(c)

	// use site
	UseSite(c, &p.SiteName, &p.SiteID, nil)

	// find page
	page := model.FindCreatePage(p.PageKey, p.PageTitle, p.SiteName)

	// check if the user is allowed to comment
	if isAllowed, resp := CheckIsAllowed(c, p.Name, p.Email, page, p.SiteName); !isAllowed {
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
	}

	// default comment type
	if !CheckIsAdminReq(c) && config.Instance.Moderator.PendingDefault {
		// 不是管理员评论 && 配置开启评论默认待审
		comment.IsPending = true
	}

	// save to database
	err := model.CreateComment(&comment)
	if err != nil {
		logrus.Error("Save Comment error: ", err)
		return RespError(c, "评论失败")
	}

	// 异步执行
	go func() {
		// Page Update
		if page.ToCooked().URL != "" && page.Title == "" {
			page.FetchURL()
		}

		// 垃圾检测
		if !CheckIsAdminReq(c) { // 忽略检查管理员
			anti_spam.SyncSpamCheck(&comment, c) // 同步执行
		}

		// 通知发送
		notify_launcher.SendNotify(&comment, &parentComment)
	}()

	return RespData(c, ResponseAdd{
		Comment: comment.ToCooked(),
	})
}
