package handler

import (
	"cmp"
	"errors"
	"fmt"

	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/i18n"
	"github.com/artalkjs/artalk/v2/internal/log"
	"github.com/artalkjs/artalk/v2/internal/utils"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsCommentCreate struct {
	Name    string `json:"name" validate:"required"`    // The comment name
	Email   string `json:"email" validate:"required"`   // The comment email
	Link    string `json:"link" validate:"optional"`    // The comment link
	Content string `json:"content" validate:"required"` // The comment content
	Rid     uint   `json:"rid" validate:"optional"`     // The comment rid
	UA      string `json:"ua" validate:"optional"`      // The comment ua

	PageKey   string `json:"page_key" validate:"required"`   // The comment page_key
	PageTitle string `json:"page_title" validate:"optional"` // The comment page_title

	SiteName string `json:"site_name" validate:"required"` // The site name of your content scope
}

type ResponseCommentCreate struct {
	entity.CookedComment
}

// @Id           CreateComment
// @Summary      Create Comment
// @Description  Create a new comment
// @Tags         Comment
// @Param        comment  body  ParamsCommentCreate  true  "The comment data"
// @Security     ApiKeyAuth
// @Success      200  {object}  ResponseCommentCreate
// @Failure      400  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Accept       json
// @Produce      json
// @Router       /comments  [post]
func CommentCreate(app *core.App, router fiber.Router) {
	router.Post("/comments", common.LimiterGuard(app, func(c *fiber.Ctx) error {
		var p ParamsCommentCreate
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		if !utils.ValidateEmail(p.Email) {
			return common.RespError(c, 400, i18n.T("Invalid {{name}}", Map{"name": i18n.T("Email")}))
		}
		if p.Link != "" && !utils.ValidateURL(p.Link) {
			return common.RespError(c, 400, i18n.T("Invalid {{name}}", Map{"name": i18n.T("Link")}))
		}

		if _, ok, resp := common.CheckSiteExist(app, c, p.SiteName); !ok {
			return resp
		}

		// Prepare the arguments for creating comment
		var (
			ip         = c.IP()
			ua         = cmp.Or(p.UA, string(c.Request().Header.UserAgent())) // allows the patched UA from the post data
			referer    = cmp.Or(c.Get("Referer"), c.Get("Origin"))
			isAdmin    = common.CheckIsAdminReq(app, c)
			isVerified = true // for display the verified badge
		)

		// Find or create page
		page := app.Dao().FindCreatePage(p.PageKey, p.PageTitle, p.SiteName)
		if page.Key == "" {
			log.Error("[CommentCreate] FindCreatePage error")
			return common.RespError(c, 500, i18n.T("Comment failed"))
		}

		// Check the page and the user is allowed to comment (admin only check)
		if isAllowed, resp := isAllowComment(app, c, p.Name, p.Email, page.AdminOnly); !isAllowed {
			return resp
		}

		// Check parent comment (reply a comment)
		var parentComment entity.Comment
		if p.Rid != 0 {
			parentComment = app.Dao().FindComment(p.Rid)
			if parentComment.IsEmpty() {
				return common.RespError(c, 404, i18n.T("{{name}} not found", Map{"name": i18n.T("Parent comment")}))
			}
			if parentComment.PageKey != p.PageKey {
				return common.RespError(c, 400, "Inconsistent with the page_key of the parent comment")
			}
			if !parentComment.IsAllowReply() {
				return common.RespError(c, 400, i18n.T("Cannot reply to this comment"))
			}
		}

		// Get the user data
		user, err := common.GetUserByReq(app, c) // if token is provided and a login user
		if errors.Is(err, common.ErrTokenNotProvided) {
			// Anonymous user
			isVerified = false
			if user, err = getUpdateAnonymousUser(app, p.Name, p.Email, p.Link, ip, ua); err != nil {
				return common.RespError(c, 500, err.Error())
			}
		} else if err != nil {
			// Login user error
			log.Error("[CommentCreate] Get user error: ", err)
			return common.RespError(c, 500, i18n.T("Comment failed"))
		}

		// Create new comment entity
		comment := entity.Comment{
			Content:  p.Content,
			PageKey:  page.Key,
			SiteName: p.SiteName,

			UserID: user.ID,
			IP:     ip,
			UA:     ua,

			Rid:    p.Rid,
			RootID: app.Dao().FindCommentRootID(p.Rid),

			IsPending:   false,
			IsCollapsed: false,
			IsPinned:    false,
			IsVerified:  isVerified,
		}

		// Set the default pending status
		// (if not admin and the `PendingDefault` is enabled)
		if !isAdmin && app.Conf().Moderator.PendingDefault {
			comment.IsPending = true
		}

		// Save the comment
		if err := app.Dao().CreateComment(&comment); err != nil {
			log.Error("Save Comment error: ", err)
			return common.RespError(c, 500, i18n.T("Comment failed"))
		}

		// Async jobs after comment created
		go commentCreatedJobs(app, comment, parentComment, commentCreatedJobsArguments{
			ip, ua, referer, isAdmin, isVerified, page,
		})

		// Response the comment data
		cookedComment := app.Dao().CookComment(&comment)
		cookedComment = fetchIPRegionForComment(app, cookedComment)

		return common.RespData(c, ResponseCommentCreate{
			CookedComment: cookedComment,
		})
	}))
}

// Fetch IP Region for Comment
func fetchIPRegionForComment(app *core.App, comment entity.CookedComment) entity.CookedComment {
	if app.Conf().IPRegion.Enabled {
		if ipRegionService, err := core.AppService[*core.IPRegionService](app); err == nil {
			comment.IPRegion = ipRegionService.Query(comment.IP)
		} else {
			log.Error("[IPRegionService] err: ", err)
		}
	}
	return comment
}

func isAllowComment(app *core.App, c *fiber.Ctx, name string, email string, pageAdminOnly bool) (bool, error) {
	// if the user is an admin user or page is admin only
	isAdminUser := app.Dao().IsAdminUserByNameEmail(name, email)
	if isAdminUser || pageAdminOnly {
		// then check has admin access
		if !common.CheckIsAdminReq(app, c) {
			return false, common.RespError(c, 403, i18n.T("Admin access required"), Map{"need_login": true})
		}
	}

	// if token is provided, then check token is valid
	user, err := common.GetUserByReq(app, c)
	if !errors.Is(err, common.ErrTokenNotProvided) && user.IsEmpty() {
		// need_auth_login is a hook for frontend to show login modal (new Auth api)
		return false, common.RespError(c, 401, i18n.T("Login required"), Map{"need_auth_login": true})
	}

	return true, nil
}

// Get and update anonymous user
//
// Call this function will create or find existing user, and update the user profile from the parameters.
// If auth is enabled, but anonymous is disabled, it will return an error.
// The column `link` will be updated only when social login is disabled. (GitHub issue #921)
func getUpdateAnonymousUser(app *core.App, name string, email string, link string, ip string, ua string) (entity.User, error) {
	// Check anonymous user is allowed
	if app.Conf().Auth.Enabled && !app.Conf().Auth.Anonymous {
		return entity.User{}, fmt.Errorf("anonymous user is not allowed")
	}

	// Create or find existing user
	user, err := app.Dao().FindCreateUser(name, email, link)
	if err != nil {
		log.Error("[CommentCreate] Create user error: ", err)
		return entity.User{}, fmt.Errorf("anonymous user create error")
	}

	// Update anonymous user profile
	// (only works when social login is disabled, @see GitHub issue #921)
	if !app.Conf().Auth.Enabled {
		user.Link = link
		user.LastIP = ip
		user.LastUA = ua

		// for modify the case of name but with the same name
		user.Name = name
		user.Email = email

		app.Dao().UpdateUser(&user)
	}

	return user, nil
}

type commentCreatedJobsArguments struct {
	IP         string
	UA         string
	Referer    string
	IsAdmin    bool
	IsVerified bool
	Page       entity.Page
}

// The jobs after comment created (should call in async)
//
// It will do the following jobs:
//  1. Page Update
//  2. AntiSpam Check
//  3. Send Notify (email, webhook, telegram, etc.)
func commentCreatedJobs(app *core.App, comment entity.Comment, parentComment entity.Comment, args commentCreatedJobsArguments) {
	// Page Update (if the original page title is empty and the URL is given)
	if app.Dao().CookPage(&args.Page).URL != "" && args.Page.Title == "" {
		go app.Dao().FetchPageFromURL(&args.Page)
	}

	// AntiSpam Check
	if !args.IsAdmin { // if the user is an admin, skip the anti-spam check
		if antiSpamService, err := core.AppService[*core.AntiSpamService](app); err == nil {
			// The check and block is sync, if comment is blocked, the message will not be sent
			antiSpamService.CheckAndBlock(&core.AntiSpamCheckPayload{
				Comment:      &comment,
				ReqReferer:   args.Referer,
				ReqIP:        args.IP,
				ReqUserAgent: args.UA,
			})
		} else {
			log.Error("[AntiSpamService] err: ", err)
		}
	}

	// Send Notify (email, webhook, telegram, etc.)
	if notifyService, err := core.AppService[*core.NotifyService](app); err == nil {
		if err := notifyService.Push(&comment, &parentComment); err != nil {
			log.Error("[NotifyService] notify push err: ", err)
		}
	} else {
		log.Error("[NotifyService] err: ", err)
	}
}
