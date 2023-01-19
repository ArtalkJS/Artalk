package handler

import (
	"strings"

	"github.com/ArtalkJS/Artalk/internal/anti_spam"
	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/notify_launcher"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ParamsAdd struct {
	Name    string `form:"name"`
	Email   string `form:"email"`
	Link    string `form:"link"`
	Content string `form:"content" validate:"required"`
	Rid     uint   `form:"rid"`
	UA      string `form:"ua"`

	PageKey   string `form:"page_key" validate:"required"`
	PageTitle string `form:"page_title"`

	Token    string `form:"token"`
	SiteName string
	SiteID   uint
}

type ResponseAdd struct {
	Comment entity.CookedComment `json:"comment"`
}

// GET /api/add
func CommentAdd(router fiber.Router) {
	router.Post("/add", func(c *fiber.Ctx) error {
		var p ParamsAdd
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		if strings.TrimSpace(p.Name) == "" {
			return common.RespError(c, i18n.T("{{name}} cannot be empty", Map{"name": i18n.T("Nickname")}))
		}
		if strings.TrimSpace(p.Email) == "" {
			return common.RespError(c, i18n.T("{{name}} cannot be empty", Map{"name": i18n.T("Email")}))
		}

		if !utils.ValidateEmail(p.Email) {
			return common.RespError(c, i18n.T("Invalid {{name}}", Map{"name": i18n.T("Email")}))
		}
		if p.Link != "" && !utils.ValidateURL(p.Link) {
			return common.RespError(c, i18n.T("Invalid {{name}}", Map{"name": i18n.T("Link")}))
		}

		ip := c.IP()
		ua := string(c.Request().Header.UserAgent())

		// 允许传入修正后的 UA
		if p.UA != "" {
			ua = p.UA
		}

		// record action for limiting action
		common.RecordAction(c)

		// use site
		common.UseSite(c, &p.SiteName, &p.SiteID, nil)

		// find page
		page := query.FindCreatePage(p.PageKey, p.PageTitle, p.SiteName)

		// check if the user is allowed to comment
		if isAllowed, resp := common.CheckIsAllowed(c, p.Name, p.Email, page, p.SiteName); !isAllowed {
			return resp
		}

		// check reply comment
		var parentComment entity.Comment
		if p.Rid != 0 {
			parentComment = query.FindComment(p.Rid)
			if parentComment.IsEmpty() {
				return common.RespError(c, i18n.T("{{name}} not found", Map{"name": i18n.T("Parent comment")}))
			}
			if parentComment.PageKey != p.PageKey {
				return common.RespError(c, "Inconsistent with the page_key of the parent comment")
			}
			if !parentComment.IsAllowReply() {
				return common.RespError(c, i18n.T("Cannot reply to this comment"))
			}
		}

		// find user
		user := query.FindCreateUser(p.Name, p.Email, p.Link)
		if user.ID == 0 || page.Key == "" {
			logrus.Error("Cannot get user or page")
			return common.RespError(c, i18n.T("Comment failed"))
		}

		// update user
		user.Link = p.Link
		user.LastIP = ip
		user.LastUA = ua
		user.Name = p.Name // for 若用户修改用户名大小写
		user.Email = p.Email
		query.UpdateUser(&user)

		comment := entity.Comment{
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
		if !common.CheckIsAdminReq(c) && config.Instance.Moderator.PendingDefault {
			// 不是管理员评论 && 配置开启评论默认待审
			comment.IsPending = true
		}

		// save to database
		err := query.CreateComment(&comment)
		if err != nil {
			logrus.Error("Save Comment error: ", err)
			return common.RespError(c, i18n.T("Comment failed"))
		}

		// 异步执行
		go func() {
			// Page Update
			if query.CookPage(&page).URL != "" && page.Title == "" {
				query.FetchPageFromURL(&page)
			}

			// 垃圾检测
			if !common.CheckIsAdminReq(c) { // 忽略检查管理员
				anti_spam.SyncSpamCheck(&comment, c) // 同步执行
			}

			// 通知发送
			notify_launcher.SendNotify(&comment, &parentComment)
		}()

		return common.RespData(c, ResponseAdd{
			Comment: query.CookComment(&comment),
		})
	})
}
