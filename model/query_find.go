package model

import (
	"fmt"
	"strings"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"gorm.io/gorm"
)

func FindComment(id uint) Comment {
	var comment Comment

	if cacher, err := FindCache(fmt.Sprintf("comment#id=%d", id), &comment); err != nil {
		cacher.StoreCache(func() interface{} {
			lib.DB.Where("id = ?", id).First(&comment)
			return &comment
		})
	}

	return comment
}

// TODO (!!no cache)
func FindCommentScopes(id uint, filters ...func(db *gorm.DB) *gorm.DB) Comment {
	var comment Comment
	lib.DB.Where("id = ?", id).Scopes(filters...).First(&comment)
	return comment
}

func FindCommentRules(id uint, rules ...func(*Comment) bool) Comment {
	comment := FindComment(id)
	for _, r := range rules {
		if !r(&comment) {
			return Comment{}
		}
	}
	return comment
}

// (Cached：parent-comments)
func FindCommentChildren(parentID uint, checkers ...func(*Comment) bool) []Comment {
	var children []Comment
	var childIDs []uint

	if cacher, err := FindCache(fmt.Sprintf("parent-comments#pid=%d", parentID), &childIDs); err != nil {
		cacher.StoreCache(func() interface{} {
			lib.DB.Model(&Comment{}).Where(&Comment{Rid: parentID}).Select("id").Find(&childIDs)
			return &childIDs
		})
	}

	for _, childID := range childIDs {
		comment := FindComment(childID)
		if comment.IsEmpty() {
			continue
		}

		// 规则过滤
		if len(checkers) > 0 {
			for _, r := range checkers {
				if !r(&comment) {
					continue
				}
			}
		}

		children = append(children, comment)
	}

	return children
}

func GetUserAllCommentIDs(userID uint) []uint {
	userAllCommentIDs := []uint{}
	lib.DB.Model(&Comment{}).Select("id").Where("user_id = ?", userID).Find(&userAllCommentIDs)
	return userAllCommentIDs
}

// 查找用户 (精确查找 name & email)
func FindUser(name string, email string) User {
	var user User

	// 查询缓存
	if cacher, err := FindCache(fmt.Sprintf("user#name=%s;email=%s", strings.ToLower(name), strings.ToLower(email)), &user); err != nil {
		cacher.StoreCache(func() interface{} {
			// 不区分大小写
			lib.DB.Where("LOWER(name) = LOWER(?) AND LOWER(email) = LOWER(?)", name, email).First(&user)
			return &user
		})
	}

	return user
}

// 查找用户 (通过 ID)
func FindUserByID(id uint) User {
	var user User

	// 查询缓存
	if cacher, err := FindCache(fmt.Sprintf("user#id=%d", id), &user); err != nil {
		cacher.StoreCache(func() interface{} {
			lib.DB.Where("id = ?", id).First(&user)
			return &user
		})
	}

	return user
}

func FindPage(key string, siteName string) Page {
	var page Page

	if cacher, err := FindCache(fmt.Sprintf("page#key=%s;site_name=%s", key, siteName), &page); err != nil {
		cacher.StoreCache(func() interface{} {
			lib.DB.Where(&Page{Key: key, SiteName: siteName}).First(&page)
			return &page
		})
	}

	return page
}

func FindPageByID(id uint) Page {
	var page Page

	if cacher, err := FindCache(fmt.Sprintf("page#id=%d", id), &page); err != nil {
		cacher.StoreCache(func() interface{} {
			lib.DB.Where("id = ?", id).First(&page)
			return &page
		})
	}

	return page
}

func FindSite(name string) Site {
	var site Site

	// 查询缓存
	if cacher, err := FindCache(fmt.Sprintf("site#name=%s", name), &site); err != nil {
		cacher.StoreCache(func() interface{} {
			lib.DB.Where("name = ?", name).First(&site)
			return &site
		})
	}

	return site
}

func FindSiteByID(id uint) Site {
	var site Site

	if cacher, err := FindCache(fmt.Sprintf("site#id=%d", id), &site); err != nil {
		cacher.StoreCache(func() interface{} {
			lib.DB.Where("id = ?", id).First(&site)
			return &site
		})
	}

	return site
}

func GetAllCookedSites() []CookedSite {
	var sites []Site
	lib.DB.Model(&Site{}).Find(&sites)

	var cookedSites []CookedSite
	for _, s := range sites {
		cookedSites = append(cookedSites, s.ToCooked())
	}

	return cookedSites
}

//#region Notify
func FindNotify(userID uint, commentID uint) Notify {
	var notify Notify
	lib.DB.Where("user_id = ? AND comment_id = ?", userID, commentID).First(&notify)
	return notify
}

func FindNotifyByKey(key string) Notify {
	var notify Notify
	lib.DB.Where(Notify{Key: key}).First(&notify)
	return notify
}

func FindNotifyByID(id uint) Notify {
	var notify Notify
	lib.DB.First(&notify, id)
	return notify
}

func FindUnreadNotifies(userID uint) []CookedNotify {
	if userID == 0 {
		return []CookedNotify{}
	}

	var notifies []Notify
	lib.DB.Where("user_id = ? AND is_read = ?", userID, false).Find(&notifies)

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
	lib.DB.Model(&Vote{}).Where("target_id = ? AND type = ?", targetID, voteType).Count(&num)
	return int(num)
}

func GetVoteNumUpDown(targetID uint, voteTo string) (int, int) {
	var up int64
	var down int64
	lib.DB.Model(&Vote{}).Where("target_id = ? AND type = ?", targetID, voteTo+"_up").Count(&up)
	lib.DB.Model(&Vote{}).Where("target_id = ? AND type = ?", targetID, voteTo+"_down").Count(&down)
	return int(up), int(down)
}

//#endregion

//#region 管理员账号检测
var allAdmins *[]User = nil

func GetAllAdmins() []User {
	if allAdmins == nil {
		var admins []User
		lib.DB.Where(&User{IsAdmin: true}).Find(&admins)
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
