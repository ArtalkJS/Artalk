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
	PageKey string `mapstructure:"page_key" param:"required"`
	Token   string `mapstructure:"token"`
}

type ResponseAdd struct {
	Comment model.CookedComment `json:"comment"`
}

func ActionAdd(c echo.Context) error {
	var p ParamsAdd
	if isOK, resp := ParamsDecode(c, ParamsAdd{}, &p); !isOK {
		return resp
	}

	if !ValidateEmail(p.Email) {
		return RespError(c, "Invalid email.")
	}
	if p.Link != "" && !ValidateURL(p.Link) {
		return RespError(c, "Invalid link.")
	}

	ip := c.RealIP()
	ua := c.Request().UserAgent()

	// record action for limiting action
	RecordAction(c)

	// find user
	user := FindUser(p.Name, p.Email)
	if user.IsEmpty() {
		user = NewUser(p.Name, p.Email, p.Link) // save a new user
	}

	// admin user check
	if !user.IsEmpty() && user.IsAdmin() {
		if !CheckIsAdminReq(c) {
			return RespError(c, "需要验证管理员身份", Map{"need_password": true})
		}
	}

	// find page
	page := FindPage(p.PageKey)
	if page.IsEmpty() {
		page = NewPage(p.PageKey)
	}

	// check if the user is allowed to comment
	if isAllowed, resp := CheckIfAllowed(c, user, page); !isAllowed {
		return resp
	}

	if user.ID == 0 || page.Key == "" {
		logrus.Error("Cannot get real user and page.")
		return RespError(c, "评论失败")
	}

	// check reply comment
	var parentComment model.Comment
	if p.Rid != 0 {
		parentComment = FindComment(p.Rid)
		if parentComment.IsEmpty() {
			return RespError(c, "找不到父评论")
		}
		if parentComment.PageKey != p.PageKey {
			return RespError(c, "与父评论的 pageKey 不一致")
		}
	}

	comment := model.Comment{
		Content: p.Content,
		Rid:     p.Rid,
		UserID:  user.ID,
		PageKey: page.Key,
		IP:      ip,
		UA:      ua,
	}

	// default comment type
	if config.Instance.Moderator.PendingDefault {
		comment.Type = model.CommentPending
	}

	// save to database
	err := lib.DB.Create(&comment).Error
	if err != nil {
		logrus.Error("Save Comment error: ", err)
		return RespError(c, "评论失败")
	}

	// update user
	user.Link = p.Link
	user.LastIP = ip
	user.LastUA = ua
	UpdateUser(&user)

	// send email
	if comment.Rid != 0 {
		email.Send(comment.ToCookedForEmail(), parentComment.ToCookedForEmail())
	}
	email.SendToAdmin(comment.ToCookedForEmail())

	return RespData(c, ResponseAdd{
		Comment: comment.ToCooked(),
	})
}
