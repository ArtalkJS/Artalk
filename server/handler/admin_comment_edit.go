package handler

import (
	"strconv"

	"github.com/ArtalkJS/Artalk/internal/email"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsCommentEdit struct {
	// 查询值
	ID       uint `form:"id" validate:"required"`
	SiteName string
	SiteID   uint
	SiteAll  bool

	// 可修改
	Content     string `form:"content"`
	PageKey     string `form:"page_key"`
	Nick        string `form:"nick"`
	Email       string `form:"email"`
	Link        string `form:"link"`
	Rid         string `form:"rid"`
	UA          string `form:"ua"`
	IP          string `form:"ip"`
	IsCollapsed bool   `form:"is_collapsed"`
	IsPending   bool   `form:"is_pending"`
	IsPinned    bool   `form:"is_pinned"`
}

type ResponseCommentEdit struct {
	Comment entity.CookedComment `json:"comment"`
}

// @Summary      Comment Edit
// @Description  Edit a specific comment
// @Tags         Comment
// @Param        id             formData  int     true   "the comment ID you want to edit"
// @Param        site_name      formData  string  false  "the site name of your content scope"
// @Param        content        formData  string  false  "the comment content"
// @Param        page_key       formData  string  false  "the comment page_key"
// @Param        nick           formData  string  false  "the comment nick"
// @Param        email          formData  string  false  "the comment email"
// @Param        link           formData  string  false  "the comment link"
// @Param        rid            formData  string  false  "the comment rid"
// @Param        ua             formData  string  false  "the comment ua"
// @Param        ip             formData  string  false  "the comment ip"
// @Param        is_collapsed   formData  bool    false  "the comment is_collapsed"
// @Param        is_pending     formData  bool    false  "the comment is_pending"
// @Param        is_pinned      formData  bool    false  "the comment is_pinned"
// @Security     ApiKeyAuth
// @Success      200  {object}  common.JSONResult{data=ResponseCommentEdit}
// @Router       /admin/comment-edit [post]
func AdminCommentEdit(router fiber.Router) {
	router.Post("/comment-edit", func(c *fiber.Ctx) error {
		var p ParamsCommentEdit
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// use site
		common.UseSite(c, &p.SiteName, &p.SiteID, &p.SiteAll)

		// find comment
		comment := query.FindComment(p.ID)
		if comment.IsEmpty() {
			return common.RespError(c, i18n.T("{{name}} not found", Map{"name": i18n.T("Comment")}))
		}

		if !common.IsAdminHasSiteAccess(c, comment.SiteName) {
			return common.RespError(c, i18n.T("Access denied"))
		}

		// check params
		if p.Email != "" && !utils.ValidateEmail(p.Email) {
			return common.RespError(c, i18n.T("Invalid {{name}}", Map{"name": i18n.T("Email")}))
		}
		if p.Link != "" && !utils.ValidateURL(p.Link) {
			return common.RespError(c, i18n.T("Invalid {{name}}", Map{"name": i18n.T("Link")}))
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
		originalUser := query.FetchUserForComment(&comment)
		if p.Nick == "" {
			p.Nick = originalUser.Name
		}
		if p.Email == "" {
			p.Email = originalUser.Email
		}

		// find or save new user
		user := query.FindCreateUser(p.Nick, p.Email, p.Link)
		if user.ID != comment.UserID {
			comment.UserID = user.ID
		}

		// user link update
		if p.Link != "" && p.Link != user.Link {
			user.Link = p.Link
			query.UpdateUser(&user)
		}

		// pageKey
		if p.PageKey != "" && p.PageKey != comment.PageKey {
			query.FindCreatePage(p.PageKey, "", p.SiteName)
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

		if err := query.UpdateComment(&comment); err != nil {
			return common.RespError(c, i18n.T("{{name}} save failed", Map{"name": i18n.T("Comment")}))
		}

		return common.RespData(c, ResponseCommentEdit{
			Comment: query.CookComment(&comment),
		})
	})
}

func RenotifyWhenPendingModified(comment *entity.Comment) {
	if comment.Rid == 0 {
		return // Root 评论不发送通知，因为这个评论已经被管理员看到了
	}

	pComment := query.FindComment(comment.Rid)
	if query.FetchUserForComment(&pComment).IsAdmin {
		return // 回复对象是管理员，则不再发送通知，因为已经看到了
	}

	if comment.UserID == pComment.UserID {
		return // 自己回复自己，不通知
	}

	notify := query.FindCreateNotify(pComment.UserID, comment.ID)
	if notify.IsEmailed {
		return // 邮件已经发送过，则不再重复发送
	}

	notify.SetComment(*comment)
	query.NotifySetInitial(&notify)

	// 邮件通知
	email.AsyncSend(&notify)
}
