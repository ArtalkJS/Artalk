package handler

import (
	"fmt"

	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/i18n"
	"github.com/artalkjs/artalk/v2/internal/log"
	"github.com/artalkjs/artalk/v2/internal/utils"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsCommentUpdate struct {
	SiteName string `json:"site_name" validate:"required"` // The site name of your content scope

	Content     string `json:"content" validate:"required"`      // The comment content
	PageKey     string `json:"page_key" validate:"required"`     // The comment page_key
	Nick        string `json:"nick" validate:"optional"`         // The comment nick
	Email       string `json:"email" validate:"optional"`        // The comment email
	Link        string `json:"link" validate:"optional"`         // The comment link
	Rid         uint   `json:"rid" validate:"required"`          // The comment rid
	UA          string `json:"ua" validate:"optional"`           // The comment ua
	IP          string `json:"ip" validate:"optional"`           // The comment ip
	IsCollapsed bool   `json:"is_collapsed" validate:"required"` // The comment is_collapsed
	IsPending   bool   `json:"is_pending" validate:"required"`   // The comment is_pending
	IsPinned    bool   `json:"is_pinned" validate:"required"`    // The comment is_pinned
}

type ResponseCommentUpdate struct {
	entity.CommentDTO
}

// @Id           UpdateComment
// @Summary      Update Comment
// @Description  Update a specific comment
// @Tags         Comment
// @Param        id             path  int                true  "The comment ID you want to update"
// @Param        comment        body  ParamsCommentUpdate  true  "The comment data"
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponseCommentUpdate
// @Failure      400  {object}  Map{msg=string}
// @Failure      403  {object}  Map{msg=string}
// @Failure      404  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /comments/{id} [put]
func CommentUpdate(app *core.App, router fiber.Router) {
	router.Put("/comments/:id", common.AdminGuard(app, func(c *fiber.Ctx) error {
		var p ParamsCommentUpdate
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		id, _ := c.ParamsInt("id")

		// check site exist
		if _, ok, resp := common.CheckSiteExist(app, c, p.SiteName); !ok {
			return resp
		}

		// find comment
		comment := app.Dao().FindComment(uint(id))
		if comment.IsEmpty() {
			return common.RespError(c, 404, i18n.T("{{name}} not found", Map{"name": i18n.T("Comment")}))
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

		// merge user
		originalUser := app.Dao().FetchUserForComment(&comment)
		if p.Nick == "" {
			p.Nick = originalUser.Name
		}
		if p.Email == "" {
			p.Email = originalUser.Email
		}

		// Keep existing authors, including legacy profiles that no longer pass validation.
		user := originalUser
		if p.Nick != originalUser.Name || p.Email != originalUser.Email {
			var err error
			user, err = app.Dao().FindCreateUser(p.Nick, p.Email, p.Link)
			if err != nil {
				return common.RespError(c, 500, i18n.T("{{name}} save failed", Map{"name": i18n.T("User")}))
			}
			comment.UserID = user.ID
		}

		// user link update
		if p.Link != "" && p.Link != user.Link {
			user.Link = p.Link
			if err := app.Dao().UpdateUser(&user); err != nil {
				return common.RespError(c, 500, i18n.T("{{name}} save failed", Map{"name": i18n.T("User")}))
			}
		}

		// pageKey
		if p.PageKey != "" && p.PageKey != comment.PageKey {
			page := app.Dao().FindCreatePage(p.PageKey, "", p.SiteName)
			if page.IsEmpty() {
				return common.RespError(c, 500, i18n.T("{{name}} save failed", Map{"name": i18n.T("Page")}))
			}
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

		approved := comment.IsPending && !p.IsPending
		comment.IsPending = p.IsPending
		var notify *entity.Notify
		var emailService *core.EmailService
		var saveErr error
		err := app.Dao().Transaction(func(tx *dao.Dao) error {
			if saveErr = tx.UpdateComment(&comment); saveErr != nil {
				return saveErr
			}
			if !approved {
				return nil
			}
			var err error
			notify, err = prepareApprovalNotification(tx, &comment)
			if err != nil || notify == nil {
				return err
			}
			emailService, err = core.AppService[*core.EmailService](app)
			return err
		})
		if err != nil {
			if saveErr != nil {
				return common.RespError(c, 500, i18n.T("{{name}} save failed", Map{"name": i18n.T("Comment")}))
			}
			log.Error("[RenotifyWhenPendingModified] error: ", err)
			return common.RespError(c, 500, "Renotify Err: "+err.Error())
		}
		app.Dao().CacheAction(func(cache *dao.DaoCache) {
			cache.CommentCacheSave(&comment)
		})
		if notify != nil {
			emailService.AsyncSend(notify)
		}

		commentDTO := app.Dao().BuildCommentDTO(&comment)
		commentDTO = fetchIPRegionForComment(app, commentDTO)

		return common.RespData(c, ResponseCommentUpdate{
			CommentDTO: commentDTO,
		})
	}))
}

func prepareApprovalNotification(d *dao.Dao, comment *entity.Comment) (*entity.Notify, error) {
	if comment.Rid == 0 {
		return nil, nil // Root 评论不发送通知，因为这个评论已经被管理员看到了
	}

	pComment := d.FindComment(comment.Rid)
	if d.FetchUserForComment(&pComment).IsAdmin {
		return nil, nil // 回复对象是管理员，则不再发送通知，因为已经看到了
	}

	if comment.UserID == pComment.UserID {
		return nil, nil // 自己回复自己，不通知
	}

	notify := d.FindCreateNotify(pComment.UserID, comment.ID)
	if notify.IsEmailed {
		return nil, nil // 邮件已经发送过，则不再重复发送
	}

	// 设置通知为未读状态
	if err := d.NotifySetInitial(&notify); err != nil {
		return nil, fmt.Errorf("func NotifySetInitial err: %w", err)
	}
	return &notify, nil
}
