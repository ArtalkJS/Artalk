package http

import (
	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/lib/email"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ParamsAdd struct {
	Name    string `mapstructure:"name" param:"required"`
	Email   string `mapstructure:"email" param:"required"`
	Link    string `mapstructure:"link"`
	Content string `mapstructure:"content" param:"required"`
	Rid     uint   `mapstructure:"rid"`

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
		parentComment = model.FindComment(p.Rid, p.SiteName)
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
	user := model.FindCreateUser(p.Name, p.Email)
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

	// update page
	if page.ToCooked().URL != "" && page.Title == "" {
		go func() {
			page.FetchURL()
		}()
	}

	// 异步执行
	go func() {
		// 垃圾检测
		if !CheckIsAdminReq(c) {
			comment.SpamCheck(c)
		}

		// 邮件通知发送
		AsyncSendEmail(&comment, &parentComment)
	}()

	return RespData(c, ResponseAdd{
		Comment: comment.ToCooked(),
	})
}

// 异步邮件发送
func AsyncSendEmail(comment *model.Comment, parentComment *model.Comment) {
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
	var admins []model.User
	lib.DB.Where("is_admin = 1").Find(&admins)

	if parentComment.IsEmpty() && len(admins) > 0 {
		// TODO: 增加用户的站点隔离，指定管理员分配网站
		for _, admin := range admins {
			if comment.UserID == admin.ID { // 管理员自己回复自己，不提醒
				continue
			}

			notify := model.FindCreateNotify(admin.ID, comment.ID)
			notify.Comment = *comment
			email.AsyncSendToAdmin(&notify, &admin) // 发送邮件给管理员
		}
	}
}
