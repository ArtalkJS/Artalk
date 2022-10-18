package model

import (
	"fmt"
	"strings"
)

func FindComment(id uint, checkers ...func(*Comment) bool) Comment {
	var comment Comment

	FindAndStoreCache(fmt.Sprintf("comment#id=%d", id), &comment, func() interface{} {
		DB().Where("id = ?", id).First(&comment)
		return &comment
	})

	// the case with checkers
	for _, c := range checkers {
		if !c(&comment) {
			return Comment{}
		}
	}

	return comment
}

// (Cached：parent-comments)
func FindCommentChildrenShallow(parentID uint, checkers ...func(*Comment) bool) []Comment {
	var children []Comment
	var childIDs []uint

	FindAndStoreCache(fmt.Sprintf("parent-comments#pid=%d", parentID), &childIDs, func() interface{} {
		DB().Model(&Comment{}).Where(&Comment{Rid: parentID}).Select("id").Find(&childIDs)
		return &childIDs
	})

	for _, childID := range childIDs {
		child := FindComment(childID, checkers...)
		if !child.IsEmpty() {
			children = append(children, child)
		}
	}

	return children
}

func FindCommentChildren(parentID uint, checkers ...func(*Comment) bool) []Comment {
	allChildren := []Comment{}
	_findCommentChildrenOnce(&allChildren, parentID, checkers...) // TODO: children 数量限制
	return allChildren
}

func _findCommentChildrenOnce(source *[]Comment, parentID uint, checkers ...func(*Comment) bool) {
	// TODO 子评论排序问题
	children := FindCommentChildrenShallow(parentID, checkers...)

	for _, child := range children {
		*source = append(*source, child)
		_findCommentChildrenOnce(source, child.ID, checkers...) // recurse
	}
}

// 查找用户 (精确查找 name & email)
func FindUser(name string, email string) User {
	var user User

	// 查询缓存
	FindAndStoreCache(fmt.Sprintf("user#name=%s;email=%s", strings.ToLower(name), strings.ToLower(email)), &user, func() interface{} {
		// 不区分大小写
		DB().Where("LOWER(name) = LOWER(?) AND LOWER(email) = LOWER(?)", name, email).First(&user)
		return &user
	})

	return user
}

// 查找用户 ID (仅根据 email)
func FindUserIdsByEmail(email string) []uint {
	var userIds = []uint{}

	// 查询缓存
	FindAndStoreCache(fmt.Sprintf("user_id#email=%s", strings.ToLower(email)), &userIds, func() interface{} {
		DB().Model(&User{}).Where("LOWER(email) = LOWER(?)", email).Pluck("id", &userIds)

		return &userIds
	})

	return userIds
}

// 查找用户 (仅根据 email)
func FindUsersByEmail(email string) []User {
	userIds := FindUserIdsByEmail(email)

	users := []User{}
	for _, id := range userIds {
		users = append(users, FindUserByID(id))
	}

	return users
}

// 查找用户 (通过 ID)
func FindUserByID(id uint) User {
	var user User

	// 查询缓存
	FindAndStoreCache(fmt.Sprintf("user#id=%d", id), &user, func() interface{} {
		DB().Where("id = ?", id).First(&user)
		return &user
	})

	return user
}

func FindPage(key string, siteName string) Page {
	var page Page

	FindAndStoreCache(fmt.Sprintf("page#key=%s;site_name=%s", key, siteName), &page, func() interface{} {
		DB().Where(&Page{Key: key, SiteName: siteName}).First(&page)
		return &page
	})

	return page
}

func FindPageByID(id uint) Page {
	var page Page

	FindAndStoreCache(fmt.Sprintf("page#id=%d", id), &page, func() interface{} {
		DB().Where("id = ?", id).First(&page)
		return &page
	})

	return page
}

func FindSite(name string) Site {
	var site Site

	// 查询缓存
	FindAndStoreCache(fmt.Sprintf("site#name=%s", name), &site, func() interface{} {
		DB().Where("name = ?", name).First(&site)
		return &site
	})

	return site
}

func FindSiteByID(id uint) Site {
	var site Site

	FindAndStoreCache(fmt.Sprintf("site#id=%d", id), &site, func() interface{} {
		DB().Where("id = ?", id).First(&site)
		return &site
	})

	return site
}

func FindAllSitesCooked() []CookedSite {
	var sites []Site
	DB().Model(&Site{}).Find(&sites)

	var cookedSites []CookedSite
	for _, s := range sites {
		cookedSites = append(cookedSites, s.ToCooked())
	}

	return cookedSites
}

//#region Notify
func FindNotify(userID uint, commentID uint) Notify {
	var notify Notify
	DB().Where("user_id = ? AND comment_id = ?", userID, commentID).First(&notify)
	return notify
}

func FindNotifyByKey(key string) Notify {
	var notify Notify
	DB().Where(Notify{Key: key}).First(&notify)
	return notify
}

func FindUnreadNotifies(userID uint) []CookedNotify {
	if userID == 0 {
		return []CookedNotify{}
	}

	var notifies []Notify
	DB().Where("user_id = ? AND is_read = ?", userID, false).Find(&notifies)

	cookedNotifies := []CookedNotify{}
	for _, n := range notifies {
		cookedNotifies = append(cookedNotifies, n.ToCooked())
	}

	return cookedNotifies
}

//#endregion

//#region Vote
func GetVoteNum(targetID uint, voteType string) int {
	var num int64
	DB().Model(&Vote{}).Where("target_id = ? AND type = ?", targetID, voteType).Count(&num)
	return int(num)
}

func GetVoteNumUpDown(targetID uint, voteTo string) (int, int) {
	var up int64
	var down int64
	DB().Model(&Vote{}).Where("target_id = ? AND type = ?", targetID, voteTo+"_up").Count(&up)
	DB().Model(&Vote{}).Where("target_id = ? AND type = ?", targetID, voteTo+"_down").Count(&down)
	return int(up), int(down)
}

//#endregion

//#region 管理员账号检测
var allAdmins *[]User = nil

func GetAllAdmins() []User {
	if allAdmins == nil {
		var admins []User
		DB().Where(&User{IsAdmin: true}).Find(&admins)
		allAdmins = &admins
	}

	return *allAdmins
}

func GetAllAdminIDs() []uint {
	admins := GetAllAdmins()
	ids := []uint{}
	for _, a := range admins {
		ids = append(ids, a.ID)
	}
	return ids
}

func IsAdminUser(userID uint) bool {
	admins := GetAllAdmins()
	for _, admin := range admins {
		if admin.ID == userID {
			return true
		}
	}

	return false
}

func IsAdminUserByNameEmail(name string, email string) bool {
	admins := GetAllAdmins()
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
