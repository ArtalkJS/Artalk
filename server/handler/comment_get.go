package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ParamsGet struct {
	PageKey  string `query:"page_key" json:"page_key" validate:"required"` // The comment page_key
	SiteName string `query:"site_name" json:"site_name"`                   // The site name of your content scope

	Limit  int `query:"limit" json:"limit"`   // The limit for pagination
	Offset int `query:"offset" json:"offset"` // The offset for pagination

	FlatMode      bool   `query:"flat_mode" json:"flat_mode"`                             // Enable flat_mode
	SortBy        string `query:"sort_by" json:"sort_by" enums:"date_asc,date_desc,vote"` // Sort by condition
	ViewOnlyAdmin bool   `query:"view_only_admin" json:"view_only_admin"`                 // Only show comments by admin

	Search string `query:"search" json:"search"` // Search keywords

	// Message Center

	Type  string `query:"type" json:"type" enums:"all,mentions,mine,pending,admin_all,admin_pending"` // Message center show type
	Name  string `query:"name" json:"name"`                                                           // The username
	Email string `query:"email" json:"email"`                                                         // The user email
}

type CommentGetOptions struct {
	Params ParamsGet

	SiteID  uint // The site id of your content scope
	SiteAll bool // The site all of your content scope

	IsMsgCenter bool         // Is message center
	User        *entity.User // The user
	IsAdminReq  bool         // Is admin request
}

type ResponseGet struct {
	Data        []entity.CookedComment `json:"comments"`
	Total       int64                  `json:"total"`
	TotalRoots  int64                  `json:"total_roots"`
	Page        entity.CookedPage      `json:"page"`
	Unread      []entity.CookedNotify  `json:"unread"`
	UnreadCount int                    `json:"unread_count"`
	ApiVersion  common.ApiVersionData  `json:"api_version"`
}

// @Summary      Get Comment List
// @Description  Get a list of comments by some conditions
// @Tags         Comment
// @Security     ApiKeyAuth
// @Param        options  query  ParamsGet  true  "The options"
// @Produce      json
// @Success      200  {object}  ResponseGet
// @Failure      500  {object}  Map{msg=string}
// @Router       /comments  [get]
func CommentGet(app *core.App, router fiber.Router) {
	router.Get("/comments", func(c *fiber.Ctx) error {
		var p ParamsGet
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		var opts CommentGetOptions
		opts.Params = p

		// use site
		site := common.GetSiteInfo(c)
		opts.SiteID = site.ID
		opts.SiteAll = site.All

		// handle params
		UseCfgFrontend(app, &p)

		// find page
		var page entity.Page
		if !opts.SiteAll {
			page = app.Dao().FindPage(p.PageKey, p.SiteName)
			if page.IsEmpty() { // if page not found
				page = entity.Page{
					Key:      p.PageKey,
					SiteName: p.SiteName,
				}
			}
		}

		// find user
		var user entity.User
		if p.Name != "" && p.Email != "" {
			user = app.Dao().FindUser(p.Name, p.Email)
			opts.User = &user // init params user field
		}

		// check if admin
		if common.CheckIsAdminReq(app, c) {
			opts.IsAdminReq = true
		}

		// check if msg center
		if p.Type != "" && p.Name != "" && p.Email != "" {
			opts.IsMsgCenter = true
			p.FlatMode = true // 通知中心强制平铺模式
		}

		// prepare the first query
		findScopes := []func(*gorm.DB) *gorm.DB{}
		if !p.FlatMode {
			// nested_mode prepare the root comments as first query result
			findScopes = append(findScopes, RootComments())
		}
		if !opts.IsMsgCenter {
			// pinned comments ignore
			findScopes = append(findScopes, func(d *gorm.DB) *gorm.DB {
				return d.Where("is_pinned = ?", false) // 因为置顶是独立的查询，这里就不再查)
			})
		}

		// search function
		if p.Search != "" {
			findScopes = append(findScopes, CommentSearchScope(app, p))
		}

		// get comments for the first query
		var comments []entity.Comment
		GetCommentQuery(app, c, opts, opts.SiteID, findScopes...).Scopes(Paginate(p.Offset, p.Limit)).Find(&comments)

		// prepend the pinned comments
		prependPinnedComments(app, c, opts, &comments)

		// cook
		cookedComments := app.Dao().CookAllComments(comments)

		switch {
		case !p.FlatMode:
			// ==========
			// 层级嵌套模式
			// ==========

			// 获取 comment 子评论
			for _, parent := range cookedComments { // TODO: Read more children, pagination for children comment
				children := app.Dao().FindCommentChildren(parent.ID, SiteIsolationChecker(app, c, opts), AllowedCommentChecker(app, c, opts))
				cookedComments = append(cookedComments, app.Dao().CookAllComments(children)...)
			}

		case p.FlatMode:
			// ==========
			// 平铺模式
			// ==========

			// find linked comments (被引用的评论，不单独显示)
			for _, comment := range comments {
				if comment.Rid == 0 || entity.ContainsCookedComment(cookedComments, comment.Rid) {
					continue
				}

				rComment := app.Dao().FindComment(comment.Rid, SiteIsolationChecker(app, c, opts)) // 查找被回复的评论
				if rComment.IsEmpty() {
					continue
				}

				rCooked := app.Dao().CookComment(&rComment)
				rCooked.Visible = false // 设置为不可见
				cookedComments = append(cookedComments, rCooked)
			}
		}

		// count comments
		total := CountComments(GetCommentQuery(app, c, opts, opts.SiteID))
		totalRoots := CountComments(GetCommentQuery(app, c, opts, opts.SiteID, RootComments()))

		// mark all as read in msg center
		if opts.IsMsgCenter {
			app.Dao().UserNotifyMarkAllAsRead(opts.User.ID)
		}

		// unread notifies
		var unreadNotifies = []entity.CookedNotify{}
		if opts.User != nil {
			unreadNotifies = app.Dao().CookAllNotifies(app.Dao().FindUnreadNotifies(opts.User.ID))
		}

		// IP region query
		if app.Conf().IPRegion.Enabled {
			ipRegionService, err := core.AppService[*core.IPRegionService](app)
			if err == nil {
				nCookedComments := []entity.CookedComment{}
				for _, c := range cookedComments {
					rawC := app.Dao().FindComment(c.ID)

					c.IPRegion = ipRegionService.Query(rawC.IP)
					nCookedComments = append(nCookedComments, c)
				}
				cookedComments = nCookedComments
			} else {
				log.Error("[IPRegionService] err: ", err)
			}
		}

		resp := ResponseGet{
			Data:        cookedComments,
			Total:       total,
			TotalRoots:  totalRoots,
			Page:        app.Dao().CookPage(&page),
			Unread:      unreadNotifies,
			UnreadCount: len(unreadNotifies),
			ApiVersion:  common.GetApiVersionDataMap(),
		}

		return common.RespData(c, resp)
	})
}
