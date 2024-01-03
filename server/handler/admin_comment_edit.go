package handler

import (
	"fmt"
	"strconv"

	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsCommentEdit struct {
	SiteName string `json:"site_name"` // The site name of your content scope

	Content     string `json:"content"`      // The comment content
	PageKey     string `json:"page_key"`     // The comment page_key
	Nick        string `json:"nick"`         // The comment nick
	Email       string `json:"email"`        // The comment email
	Link        string `json:"link"`         // The comment link
	Rid         string `json:"rid"`          // The comment rid
	UA          string `json:"ua"`           // The comment ua
	IP          string `json:"ip"`           // The comment ip
	IsCollapsed bool   `json:"is_collapsed"` // The comment is_collapsed
	IsPending   bool   `json:"is_pending"`   // The comment is_pending
	IsPinned    bool   `json:"is_pinned"`    // The comment is_pinned
}

type ResponseCommentEdit struct {
	Data entity.CookedComment `json:"data"`
}

// @Summary      Edit Comment
// @Description  Edit a specific comment
// @Tags         Comment
// @Param        id             path  int                true  "The comment ID you want to edit"
// @Param        comment        body  ParamsCommentEdit  true  "The comment data"
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponseCommentEdit
// @Failure      400  {object}  Map{msg=string}
// @Failure      403  {object}  Map{msg=string}
// @Failure      404  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /comments/{id} [put]
func AdminCommentEdit(app *core.App, router fiber.Router) {
	router.Put("/comments/:id", func(c *fiber.Ctx) error {
		var p ParamsCommentEdit
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		id, _ := c.ParamsInt("id")

		// find comment
		comment := app.Dao().FindComment(uint(id))
		if comment.IsEmpty() {
			return common.RespError(c, 404, i18n.T("{{name}} not found", Map{"name": i18n.T("Comment")}))
		}

		if !common.IsAdminHasSiteAccess(app, c, comment.SiteName) {
			return common.RespError(c, 403, i18n.T("Access denied"))
		}

		// check params
		if p.Email != "" && !utils.ValidateEmail(p.Email) {
			return common.RespError(c, 400, i18n.T("Invalid {{name}}", Map{"name": i18n.T("Email")}))
		}
		if p.Link != "" && !utils.ValidateURL(p.Link) {
			return common.RespError(c, 400, i18n.T("Invalid {{name}}", Map{"name": i18n.T("Link")}))
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
		originalUser := app.Dao().FetchUserForComment(&comment)
		if p.Nick == "" {
			p.Nick = originalUser.Name
		}
		if p.Email == "" {
			p.Email = originalUser.Email
		}

		// find or save new user
		user := app.Dao().FindCreateUser(p.Nick, p.Email, p.Link)
		if user.ID != comment.UserID {
			comment.UserID = user.ID
		}

		// user link update
		if p.Link != "" && p.Link != user.Link {
			user.Link = p.Link
			app.Dao().UpdateUser(&user)
		}

		// pageKey
		if p.PageKey != "" && p.PageKey != comment.PageKey {
			app.Dao().FindCreatePage(p.PageKey, "", p.SiteName)
			comment.PageKey = p.PageKey
		}

		if p.UA != "" {
			comment.UA = p.UA
		}
		if p.IP != "" {
			comment.IP = p.IP
		}

		comment.IsCollapsed = p.IsCollapsed
		comment.IsPinned = p.IsPinned

		if p.IsPending != comment.IsPending {
			// 待审状态发生改变
			comment.IsPending = p.IsPending

			// 待审状态被修改为 false，则重新发送邮件通知
			if !comment.IsPending {
				if err := RenotifyWhenPendingModified(app, &comment); err != nil {
					log.Error("[RenotifyWhenPendingModified] error: ", err)
					return common.RespError(c, 500, "Renotify Err: "+err.Error())
				}
			}
		}

		if err := app.Dao().UpdateComment(&comment); err != nil {
			return common.RespError(c, 500, i18n.T("{{name}} save failed", Map{"name": i18n.T("Comment")}))
		}

		return common.RespData(c, ResponseCommentEdit{
			Data: app.Dao().CookComment(&comment),
		})
	})
}

func RenotifyWhenPendingModified(app *core.App, comment *entity.Comment) (err error) {
	if comment.Rid == 0 {
		return // Root 评论不发送通知，因为这个评论已经被管理员看到了
	}

	pComment := app.Dao().FindComment(comment.Rid)
	if app.Dao().FetchUserForComment(&pComment).IsAdmin {
		return // 回复对象是管理员，则不再发送通知，因为已经看到了
	}

	if comment.UserID == pComment.UserID {
		return // 自己回复自己，不通知
	}

	notify := app.Dao().FindCreateNotify(pComment.UserID, comment.ID)
	if notify.IsEmailed {
		return // 邮件已经发送过，则不再重复发送
	}

	// 设置通知为未读状态
	if err := app.Dao().NotifySetInitial(&notify); err != nil {
		return fmt.Errorf("func NotifySetInitial err: %w", err)
	}

	// 邮件通知
	emailService, err := core.AppService[*core.EmailService](app)
	if err != nil {
		return err
	}
	emailService.AsyncSend(&notify)

	return
}
