package http

import (
	"strings"

	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ParamsGet struct {
	PageKey string `mapstructure:"page_key" param:"required"`
	Limit   int    `mapstructure:"limit"`
	Offset  int    `mapstructure:"offset"`

	// Message Center
	Type  string `mapstructure:"type"`
	Name  string `mapstructure:"name"`
	Email string `mapstructure:"email"`
	User  *model.User

	SiteName string `mapstructure:"site_name"`
	SiteID   uint
	SiteAll  bool

	FlatMode bool `mapstructure:"flat_mode"`

	IsAdminReq bool

	// Sort By
	SortBy string `mapstructure:"sort_by"` // date_asc, date_desc, vote

	// 只看 admin
	ViewOnlyAdmin bool `mapstructure:"view_only_admin"`
}

type ResponseGet struct {
	Comments    []model.CookedComment `json:"comments"`
	Total       int64                 `json:"total"`
	TotalRoots  int64                 `json:"total_roots"`
	Page        model.CookedPage      `json:"page"`
	Unread      []model.CookedNotify  `json:"unread"`
	UnreadCount int                   `json:"unread_count"`
	ApiVersion  Map                   `json:"api_version"`
	Conf        Map                   `json:"conf"`
}

// 获取评论查询实例
func GetCommentQuery(a *action, c echo.Context, p ParamsGet, siteID uint) *gorm.DB {
	query := a.db.Model(&model.Comment{}).Order(GetSortRuleSQL(p.SortBy, "created_at DESC"))

	if IsMsgCenter(p) {
		return query.Scopes(MsgCenter(a, c, p, siteID), SiteIsolation(c, p))
	}

	return query.Where("page_key = ?", p.PageKey).Scopes(SiteIsolation(c, p), AllowedComment(c, p))
}

func (a *action) Get(c echo.Context) error {
	var p ParamsGet
	if isOK, resp := ParamsDecode(c, ParamsGet{}, &p); !isOK {
		return resp
	}

	// find site
	if isOK, resp := CheckSite(c, &p.SiteName, &p.SiteID, &p.SiteAll); !isOK {
		return resp
	}

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
	isMsgCenter := IsMsgCenter(p)
	if isMsgCenter {
		p.FlatMode = true // 通知中心强制平铺模式
	}

	// comment parents
	var comments []model.Comment

	query := GetCommentQuery(a, c, p, p.SiteID).Scopes(Paginate(p.Offset, p.Limit))
	query = query.Scopes(ViewOnlyAdmin(c, p))       // 装载只看管理员功能
	query = query.Scopes(PinnedCommentsScope(c, p)) // 评论置顶功能总控制
	cookedComments := []model.CookedComment{}

	if !p.FlatMode {
		// 层级嵌套模式
		query = query.Scopes(RootComments())
		query.Find(&comments)

		for _, c := range comments {
			cookedComments = append(cookedComments, c.ToCooked())
		}

		pinnedCommentsFunction(a, c, p, &cookedComments) // 插入置顶评论

		// 获取 comment 子评论
		for _, parent := range cookedComments { // TODO: Read more children, pagination for children comment
			children := parent.FetchChildrenWithRules(SiteIsolationRule(c, p), AllowedCommentRule(c, p))
			for _, child := range children {
				cookedComments = append(cookedComments, child)
			}
		}
	} else {
		// 平铺模式
		query.Find(&comments)

		for _, c := range comments {
			cookedComments = append(cookedComments, c.ToCooked())
		}

		pinnedCommentsFunction(a, c, p, &cookedComments) // 插入置顶评论

		containsComment := func(cid uint) bool {
			for _, c := range cookedComments {
				if c.ID == cid {
					return true
				}
			}
			return false
		}

		// find linked comments (被引用的评论，不单独显示)
		for _, comment := range comments {
			if comment.Rid == 0 || containsComment(comment.Rid) {
				continue
			}

			rComment := model.FindCommentRules(comment.Rid, SiteIsolationRule(c, p)) // 查找被回复的评论
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
	totalRoots := CountComments(GetCommentQuery(a, c, p, p.SiteID).Scopes(RootComments()))

	if isMsgCenter {
		// mark all as read
		model.UserNotifyMarkAllAsRead(p.User.ID)
	}

	// unread notifies
	var unreadNotifies = []model.CookedNotify{}
	if p.User != nil {
		unreadNotifies = model.FindUnreadNotifies(p.User.ID)
	}

	return RespData(c, ResponseGet{
		Comments:    cookedComments,
		Total:       total,
		TotalRoots:  totalRoots,
		Page:        page.ToCooked(),
		Unread:      unreadNotifies,
		UnreadCount: len(unreadNotifies),
		ApiVersion:  GetApiVersionDataMap(),
		Conf:        GetApiPublicConfDataMap(c),
	})
}

// 请求是否为 通知中心数据
func IsMsgCenter(p ParamsGet) bool {
	return p.Type != "" && p.Name != "" && p.Email != ""
}

// TODO: 重构 MsgCenter
func MsgCenter(a *action, c echo.Context, p ParamsGet, siteID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		user := p.User

		if user == nil || user.IsEmpty() { // user not found
			return db.Where("id = 0")
		}

		isAdminReq := CheckIsAdminReq(c)

		// admin_only 检测
		if strings.HasPrefix(p.Type, "admin_") && !isAdminReq {
			db = db.Where("id = 0")
			return db
		}

		// 获取我的 commentIDs
		getMyCommentIDs := func() []int {
			myCommentIDs := []int{}
			var myComments []model.Comment
			a.db.Where("user_id = ?", user.ID).Find(&myComments)
			for _, comment := range myComments {
				myCommentIDs = append(myCommentIDs, int(comment.ID))
			}
			return myCommentIDs
		}

		switch p.Type {
		case "all":
			return db.Where("rid IN (?) OR user_id = ?", getMyCommentIDs(), user.ID)
		case "mentions":
			return db.Where("rid IN (?) AND user_id != ?", getMyCommentIDs(), user.ID)
		case "mine":
			return db.Where("user_id = ?", user.ID)
		case "pending":
			return db.Where("user_id = ? AND is_pending = 1", user.ID)
		case "admin_all":
			return db
		case "admin_pending":
			return db.Where("is_pending = 1")
		}

		return db.Where("id = 0")
	}
}

// 评论计数
func CountComments(db *gorm.DB) int64 {
	var count int64
	db.Count(&count)
	return count
}

// 分页
func Paginate(offset int, limit int) func(db *gorm.DB) *gorm.DB {
	if offset < 0 {
		offset = 0
	}

	if limit > 100 {
		limit = 100
	} else if limit <= 0 {
		limit = 15
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(limit)
	}
}

// 允许的评论
func AllowedComment(c echo.Context, p ParamsGet) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if CheckIsAdminReq(c) {
			return db // 管理员显示全部
		}

		// 显示个人全部评论
		if p.Name != "" && p.Email != "" {
			if !p.User.IsEmpty() {
				return db.Where("is_pending = 0 OR (is_pending = 1 AND user_id = ?)", p.User.ID)
			}
		}

		return db.Where("is_pending = 0") // 不允许待审评论
	}
}

// 站点隔离
func SiteIsolation(c echo.Context, p ParamsGet) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if CheckIsAdminReq(c) && p.SiteAll {
			return db // 仅管理员支持取消站点隔离
		}

		return db.Where("site_name = ?", p.SiteName)
	}
}

// TODO 实验，以后全部 Scopes 方法换成 Rules 方法
func AllowedCommentRule(c echo.Context, p ParamsGet) func(*model.Comment) bool {
	return func(comment *model.Comment) bool {
		if CheckIsAdminReq(c) {
			return true // 管理员显示全部
		}

		// 显示个人全部评论
		if (p.Name != "" && p.Email != "") && !p.User.IsEmpty() {
			if !comment.IsPending || (comment.IsPending && comment.UserID == p.User.ID) {
				return true
			}
		}

		// 不允许待审评论
		if comment.IsPending {
			return false
		}

		return true
	}
}

// TODO
func SiteIsolationRule(c echo.Context, p ParamsGet) func(*model.Comment) bool {
	return func(comment *model.Comment) bool {
		if CheckIsAdminReq(c) && p.SiteAll {
			return true // 仅管理员支持取消站点隔离
		}

		if comment.SiteName != p.SiteName {
			return false
		}

		return true
	}
}

// 只看管理员功能
func ViewOnlyAdmin(c echo.Context, p ParamsGet) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 是否启用
		if !p.ViewOnlyAdmin {
			return db
		}

		// 获取管理员列表
		adminIDs := model.GetAllAdminIDs()

		// 只允许管理员 user_id
		return db.Where("user_id IN ?", adminIDs)
	}
}

func PinnedCommentsScope(c echo.Context, p ParamsGet) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if IsMsgCenter(p) {
			// 通知中心关闭置顶
			return db
		}

		if p.Offset == 0 {
			return db
		} else {
			// 其他页面不再显示置顶内容
			return db.Where("is_pinned = 0")
		}
	}
}

func pinnedCommentsFunction(a *action, c echo.Context, p ParamsGet, cookedComments *[]model.CookedComment) {
	if IsMsgCenter(p) {
		// 通知中心关闭置顶
		return
	}

	// 仅在分页的首页加入置顶评论
	if p.Offset != 0 {
		return
	}

	pinnedComments := []model.Comment{}
	GetCommentQuery(a, c, p, p.SiteID).Where("is_pinned = 1").Find(&pinnedComments)

	if len(pinnedComments) == 0 {
		return // 没有置顶评论
	}

	// cook
	pinnedCookedComments := []model.CookedComment{}
	for _, pCo := range pinnedComments {
		pinnedCookedComments = append(pinnedCookedComments, pCo.ToCooked())
	}

	// 去掉已 pinned 且重复存在于原列表中的评论
	filteredCookedComments := []model.CookedComment{}
	for _, co := range *cookedComments {
		isExistInPinnedList := false
		for _, pCo := range pinnedComments {
			if co.ID == pCo.ID {
				isExistInPinnedList = true
				break
			}
		}
		if !isExistInPinnedList {
			filteredCookedComments = append(filteredCookedComments, co)
		}
	}

	// prepend
	*cookedComments = append(pinnedCookedComments, filteredCookedComments...)
}

// 排序规则
func GetSortRuleSQL(sortBy string, defaultSQL string) string {
	switch sortBy {
	case "date_desc":
		return "created_at DESC"
	case "date_asc":
		return "created_at ASC"
	case "vote":
		return "vote_up DESC, created_at DESC"
	}

	return defaultSQL
}

func RootComments() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("rid = 0")
	}
}
