package dao

import (
	"fmt"
	"strings"

	"github.com/artalkjs/artalk/v2/internal/cache"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"golang.org/x/sync/singleflight"
)

func QueryDBWithCache[T any](dao *Dao, name string, queryDB func() (T, error)) (T, error) {
	if dao.cache == nil {
		return QueryDBWithoutCache(dao, name, queryDB) // directly call queryDB while cache is disabled
	}
	return cache.QueryDBWithCache(dao.cache.Cache, name, queryDB)
}

var noCacheFindSingleflightGroup = new(singleflight.Group)

func QueryDBWithoutCache[T any](dao *Dao, name string, queryDB func() (T, error)) (T, error) {
	v, err, _ := noCacheFindSingleflightGroup.Do(name, func() (any, error) {
		return queryDB() // Use singleflight to prevent high concurrency pressure on DB
	})
	if err != nil {
		return *new(T), err
	}
	return v.(T), err
}

func (dao *Dao) FindComment(id uint, checkers ...func(*entity.Comment) bool) entity.Comment {
	comment, _ := QueryDBWithCache(dao, fmt.Sprintf(CommentByIDKey, id), func() (comment entity.Comment, err error) {
		dao.DB().Where("id = ?", id).First(&comment)
		return comment, nil
	})

	// the case with checkers
	for _, c := range checkers {
		if !c(&comment) {
			return entity.Comment{}
		}
	}

	return comment
}

func (dao *Dao) FindCommentRootID(rid uint) uint {
	visited := map[uint]bool{}
	rootId := rid
	for rootId != 0 && !visited[rootId] {
		visited[rootId] = true // avoid infinite loop (rid = id)

		var comment entity.Comment
		dao.DB().First(&comment, rootId)
		if comment.Rid == 0 { // if comment is root
			return rootId
		}

		rootId = comment.Rid // continue to find root
	}
	return rootId
}

// (Cached：parent-comments)
func (dao *Dao) FindCommentChildrenShallow(parentID uint, checkers ...func(*entity.Comment) bool) []entity.Comment {
	var children []entity.Comment

	childIDs, _ := QueryDBWithCache(dao, fmt.Sprintf(CommentChildIDsByIDKey, parentID), func() ([]uint, error) {
		childIDs := []uint{}
		dao.DB().Model(&entity.Comment{}).Where(&entity.Comment{Rid: parentID}).Select("id").Find(&childIDs)
		return childIDs, nil
	})

	for _, childID := range childIDs {
		child := dao.FindComment(childID, checkers...)
		if !child.IsEmpty() {
			children = append(children, child)
		}
	}

	return children
}

func (dao *Dao) FindCommentChildren(parentID uint, checkers ...func(*entity.Comment) bool) []entity.Comment {
	allChildren := []entity.Comment{}
	dao._findCommentChildrenOnce(&allChildren, parentID, checkers...) // TODO: children 数量限制
	return allChildren
}

func (dao *Dao) _findCommentChildrenOnce(source *[]entity.Comment, parentID uint, checkers ...func(*entity.Comment) bool) {
	// TODO 子评论排序问题
	children := dao.FindCommentChildrenShallow(parentID, checkers...)

	for _, child := range children {
		*source = append(*source, child)
		dao._findCommentChildrenOnce(source, child.ID, checkers...) // recurse
	}
}

// 查找用户 (精确查找 name & email)
func (dao *Dao) FindUser(name string, email string) entity.User {
	user, _ := QueryDBWithCache(dao, fmt.Sprintf(UserByNameEmailKey, strings.ToLower(name), strings.ToLower(email)), func() (user entity.User, err error) {
		dao.DB().Where("LOWER(name) = LOWER(?) AND LOWER(email) = LOWER(?)", name, email).First(&user) // 不区分大小写
		return user, nil
	})

	return user
}

// 查找用户 ID (仅根据 email)
func (dao *Dao) FindUserIdsByEmail(email string) []uint {
	userIDs, _ := QueryDBWithCache(dao, fmt.Sprintf(UserIDByEmailKey, strings.ToLower(email)), func() ([]uint, error) {
		userIDs := []uint{}
		dao.DB().Model(&entity.User{}).Where("LOWER(email) = LOWER(?)", email).Pluck("id", &userIDs)
		return userIDs, nil
	})

	return userIDs
}

// 查找用户 (仅根据 email)
func (dao *Dao) FindUsersByEmail(email string) []entity.User {
	userIds := dao.FindUserIdsByEmail(email)

	users := []entity.User{}
	for _, id := range userIds {
		users = append(users, dao.FindUserByID(id))
	}

	return users
}

// 查找用户 (通过 ID)
func (dao *Dao) FindUserByID(id uint) entity.User {
	user, _ := QueryDBWithCache(dao, fmt.Sprintf("user#id=%d", id), func() (user entity.User, err error) {
		dao.DB().Where("id = ?", id).First(&user)
		return user, nil
	})
	return user
}

func (dao *Dao) FindPage(key string, siteName string) entity.Page {
	page, _ := QueryDBWithCache(dao, fmt.Sprintf(PageByKeySiteNameKey, key, siteName), func() (page entity.Page, err error) {
		dao.DB().Where(&entity.Page{Key: key, SiteName: siteName}).First(&page)
		return page, nil
	})

	return page
}

func (dao *Dao) FindPageByID(id uint) entity.Page {
	page, _ := QueryDBWithCache(dao, fmt.Sprintf(PageByIDKey, id), func() (page entity.Page, err error) {
		dao.DB().Where("id = ?", id).First(&page)
		return page, nil
	})

	return page
}

func (dao *Dao) FindSite(name string) entity.Site {
	site, _ := QueryDBWithCache(dao, fmt.Sprintf(SiteByNameKey, name), func() (site entity.Site, err error) {
		dao.DB().Where("name = ?", name).First(&site)
		return site, nil
	})
	return site
}

func (dao *Dao) FindSiteByID(id uint) entity.Site {
	site, _ := QueryDBWithCache(dao, fmt.Sprintf(SiteByIDKey, id), func() (site entity.Site, err error) {
		dao.DB().Where("id = ?", id).First(&site)
		return site, nil
	})
	return site
}

func (dao *Dao) FindAllSites() []entity.Site {
	var sites []entity.Site
	dao.DB().Model(&entity.Site{}).Find(&sites)

	return sites
}

// #region Notify
func (dao *Dao) FindNotify(userID uint, commentID uint) entity.Notify {
	var notify entity.Notify
	dao.DB().Where("user_id = ? AND comment_id = ?", userID, commentID).First(&notify)
	return notify
}

func (dao *Dao) FindNotifyForComment(commentID uint, key string) entity.Notify {
	var notify entity.Notify
	dao.DB().Where(entity.Notify{CommentID: commentID, Key: key}).First(&notify)
	return notify
}

func (dao *Dao) FindUnreadNotifies(userID uint) []entity.Notify {
	if userID == 0 {
		return []entity.Notify{}
	}

	var notifies []entity.Notify
	dao.DB().Where("user_id = ? AND is_read = ?", userID, false).Find(&notifies)

	return notifies
}

func (dao *Dao) FindNotifyParentComment(n *entity.Notify) entity.Comment {
	comment := dao.FetchCommentForNotify(n)
	if comment.Rid == 0 {
		return entity.Comment{}
	}

	return dao.FindComment(comment.Rid)
}

//#endregion

// #region Vote
func (dao *Dao) GetVoteNum(targetID uint, voteType string) int {
	var num int64
	dao.DB().Model(&entity.Vote{}).Where("target_id = ? AND type = ?", targetID, voteType).Count(&num)
	return int(num)
}

func (dao *Dao) GetVoteNumUpDown(targetName string, targetID uint) (int, int) {
	var up int64
	var down int64
	dao.DB().Model(&entity.Vote{}).Where("target_id = ? AND type = ?", targetID, targetName+"_up").Count(&up)
	dao.DB().Model(&entity.Vote{}).Where("target_id = ? AND type = ?", targetID, targetName+"_down").Count(&down)
	return int(up), int(down)
}

//#endregion

// #region 管理员账号检测
func (dao *Dao) GetAllAdmins() []entity.User {
	// TODO add cache and flush cache when admin changed
	var admins []entity.User
	dao.DB().Where(&entity.User{IsAdmin: true}).Find(&admins)
	return admins
}

func (dao *Dao) GetAllAdminIDs() []uint {
	admins := dao.GetAllAdmins()
	ids := []uint{}
	for _, a := range admins {
		ids = append(ids, a.ID)
	}
	return ids
}

func (dao *Dao) IsAdminUser(userID uint) bool {
	admins := dao.GetAllAdmins()
	for _, admin := range admins {
		if admin.ID == userID {
			return true
		}
	}

	return false
}

func (dao *Dao) IsAdminUserByNameEmail(name string, email string) bool {
	admins := dao.GetAllAdmins()
	for _, admin := range admins {
		// Name 和 Email 都匹配才是管理员
		if strings.EqualFold(admin.Name, name) &&
			strings.EqualFold(admin.Email, email) {
			return true
		}
	}

	return false
}

//#endregion

func (dao *Dao) FindAuthIdentityByToken(provider string, token string) entity.AuthIdentity {
	var identity entity.AuthIdentity
	dao.DB().Where("provider = ? AND token = ?", provider, token).First(&identity)
	return identity
}

func (dao *Dao) FindAuthIdentityByRemoteUID(provider string, remoteUID string) entity.AuthIdentity {
	var identity entity.AuthIdentity
	dao.DB().Where("provider = ? AND remote_uid = ?", provider, remoteUID).First(&identity)
	return identity
}

func (dao *Dao) FindAuthIdentityByUserID(provider string, userID uint) entity.AuthIdentity {
	var identity entity.AuthIdentity
	dao.DB().Where("provider = ? AND user_id = ?", provider, userID).First(&identity)
	return identity
}
