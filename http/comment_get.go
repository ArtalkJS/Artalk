package http

import (
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ParamsGet struct {
	PageKey  string `mapstructure:"page_key" param:"required"`
	SiteName string

	Limit  int `mapstructure:"limit"`
	Offset int `mapstructure:"offset"`

	FlatMode      bool   `mapstructure:"flat_mode"`
	SortBy        string `mapstructure:"sort_by"`         // date_asc, date_desc, vote
	ViewOnlyAdmin bool   `mapstructure:"view_only_admin"` // 只看 admin

	Search string `mapstructure:"search"`

	// Message Center
	Type  string `mapstructure:"type"` // ["", "all", "mentions", "mine", "pending", "admin_all", "admin_pending"]
	Name  string `mapstructure:"name"`
	Email string `mapstructure:"email"`

	SiteID  uint
	SiteAll bool

	IsMsgCenter bool
	User        *model.User
	IsAdminReq  bool
}

type ResponseGet struct {
	Comments    []model.CookedComment `json:"comments"`
	Total       int64                 `json:"total"`
	TotalRoots  int64                 `json:"total_roots"`
	Page        model.CookedPage      `json:"page"`
	Unread      []model.CookedNotify  `json:"unread"`
	UnreadCount int                   `json:"unread_count"`
	ApiVersion  Map                   `json:"api_version"`
	Conf        Map                   `json:"conf,omitempty"`
}

func (a *action) Get(c echo.Context) error {
	var p ParamsGet
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	// use site
	UseSite(c, &p.SiteName, &p.SiteID, &p.SiteAll)

	// handle params
	UseCfgFrontend(&p)

	// find page
	var page model.Page
	if !p.SiteAll {
		page = model.FindPage(p.PageKey, p.SiteName)
		if page.IsEmpty() { // if page not found
			page = model.Page{
				Key:      p.PageKey,
				SiteName: p.SiteName,
			}
		}
	}

	// find user
	var user model.User
	if p.Name != "" && p.Email != "" {
		user = model.FindUser(p.Name, p.Email)
		p.User = &user // init params user field
	}

	// check if admin
	if CheckIsAdminReq(c) {
		p.IsAdminReq = true
	}

	// check if msg center
	if p.Type != "" && p.Name != "" && p.Email != "" {
		p.IsMsgCenter = true
		p.FlatMode = true // 通知中心强制平铺模式
	}

	// prepare the first query
	findScopes := []func(*gorm.DB) *gorm.DB{}
	if !p.FlatMode {
		// nested_mode prepare the root comments as first query result
		findScopes = append(findScopes, RootComments())
	}
	if !p.IsMsgCenter {
		// pinned comments ignore
		findScopes = append(findScopes, func(d *gorm.DB) *gorm.DB {
			return d.Where("is_pinned = ?", false) // 因为置顶是独立的查询，这里就不再查)
		})
	}

	// search function
	if p.Search != "" {
		findScopes = append(findScopes, CommentSearchScope(a, p))
	}

	// get comments for the first query
	var comments []model.Comment
	GetCommentQuery(a, c, p, p.SiteID, findScopes...).Scopes(Paginate(p.Offset, p.Limit)).Find(&comments)

	// prepend the pinned comments
	prependPinnedComments(a, c, p, &comments)

	// cook
	cookedComments := model.CookAllComments(comments)

	switch {
	case !p.FlatMode:
		// ==========
		// 层级嵌套模式
		// ==========

		// 获取 comment 子评论
		for _, parent := range cookedComments { // TODO: Read more children, pagination for children comment
			children := model.FindCommentChildren(parent.ID, SiteIsolationChecker(c, p), AllowedCommentChecker(c, p))
			cookedComments = append(cookedComments, model.CookAllComments(children)...)
		}

	case p.FlatMode:
		// ==========
		// 平铺模式
		// ==========

		// find linked comments (被引用的评论，不单独显示)
		for _, comment := range comments {
			if comment.Rid == 0 || model.ContainsCookedComment(cookedComments, comment.Rid) {
				continue
			}

			rComment := model.FindComment(comment.Rid, SiteIsolationChecker(c, p)) // 查找被回复的评论
			if rComment.IsEmpty() {
				continue
			}

			rCooked := rComment.ToCooked()
			rCooked.Visible = false // 设置为不可见
			cookedComments = append(cookedComments, rCooked)
		}
	}

	// count comments
	total := CountComments(GetCommentQuery(a, c, p, p.SiteID))
	totalRoots := CountComments(GetCommentQuery(a, c, p, p.SiteID, RootComments()))

	// mark all as read in msg center
	if p.IsMsgCenter {
		model.UserNotifyMarkAllAsRead(p.User.ID)
	}

	// unread notifies
	var unreadNotifies = []model.CookedNotify{}
	if p.User != nil {
		unreadNotifies = model.FindUnreadNotifies(p.User.ID)
	}

	resp := ResponseGet{
		Comments:    cookedComments,
		Total:       total,
		TotalRoots:  totalRoots,
		Page:        page.ToCooked(),
		Unread:      unreadNotifies,
		UnreadCount: len(unreadNotifies),
		ApiVersion:  GetApiVersionDataMap(),
	}

	if p.Offset == 0 {
		resp.Conf = GetApiPublicConfDataMap(c)
	}

	return RespData(c, resp)
}
