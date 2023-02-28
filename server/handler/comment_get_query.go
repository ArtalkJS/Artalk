package handler

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/db"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// 获取评论查询实例
func GetCommentQuery(c *fiber.Ctx, p ParamsGet, siteID uint, scopes ...func(*gorm.DB) *gorm.DB) *gorm.DB {
	q := db.DB().Model(&entity.Comment{})

	q.Scopes(SiteIsolationScope(c, p), AllowedComment(c, p))
	q.Order(GetSortRuleSQL(p.SortBy, "created_at DESC")) // 排序规则
	q.Scopes(scopes...)

	switch {
	case !p.IsMsgCenter:
		// ==========
		//  非通知中心
		// ==========

		// 只看管理员功能
		if p.ViewOnlyAdmin {
			adminIDs := query.GetAllAdminIDs() // 获取管理员列表
			q.Where("user_id IN ?", adminIDs)  // 只允许管理员 user_id
		}

		q.Where("page_key = ?", p.PageKey)

	case p.IsMsgCenter:
		// ==========
		//  通知中心
		// ==========
		user := p.User
		if user == nil || user.IsEmpty() {
			return q.Where("id = 0") // user not found
		}

		// admin_only 检测
		if strings.HasPrefix(p.Type, "admin_") {
			if !p.IsAdminReq || !common.IsAdminHasSiteAccess(c, p.SiteName) {
				return q.Where("id = 0")
			}
		}

		switch p.Type {
		case "all":
			q.Where("rid IN (?) OR user_id = ?", query.GetUserAllCommentIDs(user.ID), user.ID)
		case "mentions":
			q.Where("rid IN (?) AND user_id != ?", query.GetUserAllCommentIDs(user.ID), user.ID)
		case "mine":
			q.Where("user_id = ?", user.ID)
		case "pending":
			q.Where("user_id = ? AND is_pending = ?", user.ID, true)
		case "admin_all":

		case "admin_pending":
			q.Where("is_pending = ?", true)
		default:
			return q.Where("id = 0")
		}
	}

	return q
}

// 根评论 Root comments
func RootComments() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("rid = 0")
	}
}

// 允许的评论
func AllowedComment(c *fiber.Ctx, p ParamsGet) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if common.CheckIsAdminReq(c) {
			return db // 管理员显示全部
		}

		// 通知中心允许显示个人的待审状态的评论
		if p.IsMsgCenter && !p.User.IsEmpty() {
			return db.Where("is_pending = ? OR (is_pending = ? AND user_id = ?)", false, true, p.User.ID)
		}

		return db.Where("is_pending = ?", false) // 不允许待审评论
	}
}

func AllowedCommentChecker(c *fiber.Ctx, p ParamsGet) func(*entity.Comment) bool {
	return func(comment *entity.Comment) bool {
		if common.CheckIsAdminReq(c) {
			return true // 管理员显示全部
		}

		// 通知中心允许显示个人的待审状态的评论
		if p.IsMsgCenter && p.User.ID == comment.UserID {
			return true
		}

		// 不允许待审评论
		if comment.IsPending {
			return false
		}

		return true
	}
}

func SiteIsolationChecker(c *fiber.Ctx, p ParamsGet) func(*entity.Comment) bool {
	return func(comment *entity.Comment) bool {
		if common.CheckIsAdminReq(c) && p.SiteAll {
			return true // 仅管理员支持取消站点隔离
		}

		if comment.SiteName != p.SiteName {
			return false
		}

		return true
	}
}

// 站点隔离
func SiteIsolationScope(c *fiber.Ctx, p ParamsGet) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if common.CheckIsAdminReq(c) && p.SiteAll {
			return db // 仅管理员支持取消站点隔离
		}

		return db.Where("site_name = ?", p.SiteName)
	}
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

// 评论计数
func CountComments(db *gorm.DB) int64 {
	var count int64
	db.Count(&count)
	return count
}

// 评论搜索
func CommentSearchScope(p ParamsGet) func(d *gorm.DB) *gorm.DB {
	var userIds []uint
	db.DB().Model(&entity.User{}).Where(
		"LOWER(name) = LOWER(?) OR LOWER(email) = LOWER(?)", p.Search, p.Search,
	).Pluck("id", &userIds)

	return func(d *gorm.DB) *gorm.DB {
		return d.Where("user_id IN (?) OR content LIKE ? OR page_key = ? OR ip = ? OR ua = ?",
			userIds, "%"+p.Search+"%", p.Search, p.Search, p.Search)
	}
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

// 插入置顶的评论
func prependPinnedComments(c *fiber.Ctx, p ParamsGet, comments *[]entity.Comment) {
	if p.IsMsgCenter || p.Offset != 0 {
		return // 通知中心关闭置顶 & 仅在分页的首页加入置顶评论
	}

	pinnedComments := []entity.Comment{}
	GetCommentQuery(c, p, p.SiteID).Where("is_pinned = ?", true).Find(&pinnedComments)
	if len(pinnedComments) == 0 {
		return // 没有置顶评论
	}

	// 去掉已 pinned 且重复存在于原列表中的评论
	filteredComments := []entity.Comment{}
	for _, co := range *comments {
		if !entity.ContainsComment(pinnedComments, co.ID) {
			filteredComments = append(filteredComments, co)
		}
	}

	// prepend
	*comments = append(pinnedComments, filteredComments...)
}

// 处理来自配置文件的 fronted 配置项
func UseCfgFrontend(p *ParamsGet) {
	feConf := config.Instance.Frontend
	if feConf == nil || reflect.ValueOf(feConf).Kind() != reflect.Map {
		return
	}

	// pagination
	// TODO consider ignore pageSize configuration in sidebar
	(func() {
		pagination, isExist := feConf["pagination"]
		if !isExist {
			return
		}
		if reflect.ValueOf(pagination).Kind() != reflect.Map {
			return
		}

		// pagination.pageSize
		cfgPageSizeStr := fmt.Sprintf("%v", pagination.(common.Map)["pageSize"])
		confPageSize, err := strconv.Atoi(cfgPageSizeStr)
		if err == nil && confPageSize > 0 {
			p.Limit = confPageSize
		}
	})()
}
