package model

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/eko/gocache/v2/store"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func FindComment(id uint) Comment {
	var comment Comment
	lib.DB.Where("id = ?", id).First(&comment)
	return comment
}

func FindCommentScopes(id uint, filters ...func(db *gorm.DB) *gorm.DB) Comment {
	var comment Comment
	lib.DB.Where("id = ?", id).Scopes(filters...).First(&comment)
	return comment
}

// 查找用户（返回：精确查找 AND）
func FindUser(name string, email string) User {
	var user User

	// 不区分大小写
	lib.DB.Where("LOWER(name) = LOWER(?) AND LOWER(email) = LOWER(?)", name, email).First(&user)
	return user
}

func FindUserByID(id uint) User {
	var user User
	lib.DB.First(&user, id)
	return user
}

func FindCache(name string, destStruct interface{}) error {
	entry, err := lib.CACHE.Get(lib.Ctx, name)
	if err != nil {
		return err
	}

	str := entry.([]byte)
	err = json.Unmarshal(str, destStruct)
	if err != nil {
		return err
	}

	return nil
}

func StoreCache(name string, srcStruct interface{}) error {
	str, err := json.Marshal(srcStruct)
	if err != nil {
		return err
	}

	err = lib.CACHE.Set(lib.Ctx, name, []byte(str), &store.Options{})
	if err != nil {
		return err
	}

	return nil
}

func UpdateComment(comment *Comment) error {
	err := lib.DB.Save(comment).Error
	if err != nil {
		logrus.Error("Update Comment error: ", err)
	}
	return err
}

func FindSite(name string) Site {
	var site Site
	lib.DB.Where("name = ?", name).First(&site)
	return site
}

func FindSiteByID(id uint) Site {
	var site Site
	lib.DB.First(&site, id)
	return site
}

func FindNotify(userID uint, commentID uint) Notify {
	var notify Notify
	lib.DB.Where("user_id = ? AND comment_id = ?", userID, commentID).First(&notify)
	return notify
}

func FindNotifyByKey(key string) Notify {
	var notify Notify
	lib.DB.Where("`key` = ?", key).First(&notify)
	return notify
}

func FindNotifyByID(id uint) Notify {
	var notify Notify
	lib.DB.First(&notify, id)
	return notify
}

func FindCreateSite(siteName string) Site {
	site := FindSite(siteName)
	if site.IsEmpty() {
		site = NewSite(siteName, "")
	}
	return site
}

func NewSite(name string, urls string) Site {
	site := Site{
		Name: name,
		Urls: urls,
	}

	err := lib.DB.Create(&site).Error
	if err != nil {
		logrus.Error("Create Site error: ", err)
	}

	return site
}

func UpdateSite(site *Site) error {
	err := lib.DB.Save(site).Error
	if err != nil {
		logrus.Error("Update Site error: ", err)
	}
	return err
}

func FindCreatePage(pageKey string, pageTitle string, siteName string) Page {
	page := FindPage(pageKey, siteName)
	if page.IsEmpty() {
		page = NewPage(pageKey, pageTitle, siteName)
	}
	return page
}

func FindCreateUser(name string, email string, link string) User {
	user := FindUser(name, email)
	if user.IsEmpty() {
		user = NewUser(name, email, link) // save a new user
	}
	return user
}

func FindCreateNotify(userID uint, lookCommentID uint) Notify {
	notify := FindNotify(userID, lookCommentID)
	if notify.IsEmpty() {
		notify = NewNotify(userID, lookCommentID)
	}
	return notify
}

func NewUser(name string, email string, link string) User {
	user := User{
		Name:  name,
		Email: email,
		Link:  link,
	}

	err := lib.DB.Create(&user).Error
	if err != nil {
		logrus.Error("Create User error: ", err)
	}

	return user
}

func UpdateUser(user *User) error {
	err := lib.DB.Save(user).Error
	if err != nil {
		logrus.Error("Update User error: ", err)
	}

	return err
}

func FindPage(key string, siteName string) Page {
	var page Page
	lib.DB.Where("`key` = ? AND `site_name` = ?", key, siteName).First(&page)
	return page
}

func FindPageByID(id uint) Page {
	var page Page
	lib.DB.Where("id = ?", id).First(&page)
	return page
}

func NewPage(key string, pageTitle string, siteName string) Page {
	page := Page{
		Key:      key,
		Title:    pageTitle,
		SiteName: siteName,
	}

	err := lib.DB.Create(&page).Error
	if err != nil {
		logrus.Error("Create Page error: ", err)
	}

	return page
}

func UpdatePage(page *Page) error {
	err := lib.DB.Save(page).Error
	if err != nil {
		logrus.Error("Update Page error: ", err)
	}
	return err
}

func NewNotify(userID uint, commentID uint) Notify {
	notify := Notify{
		UserID:    userID,
		CommentID: commentID,
		IsRead:    false,
		IsEmailed: false,
	}
	notify.GenerateKey()

	err := lib.DB.Create(&notify).Error
	if err != nil {
		logrus.Error("Create Notify error: ", err)
	}

	return notify
}

func FindUnreadNotifies(userID uint) []CookedNotify {
	if userID == 0 {
		return []CookedNotify{}
	}

	var notifies []Notify
	lib.DB.Where("user_id = ? AND is_read = 0", userID).Find(&notifies)

	cookedNotifies := []CookedNotify{}
	for _, n := range notifies {
		cookedNotifies = append(cookedNotifies, n.ToCooked())
	}

	return cookedNotifies
}

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

func DelComment(commentID uint) error {
	// 清除 notify
	if err := lib.DB.Unscoped().Where("comment_id = ?", commentID).Delete(&Notify{}).Error; err != nil {
		return err
	}

	// 清除 vote
	if err := lib.DB.Unscoped().Where(
		"target_id = ? AND (type = ? OR type = ?)",
		commentID,
		string(VoteTypeCommentUp),
		string(VoteTypeCommentDown),
	).Delete(&Vote{}).Error; err != nil {
		return err
	}

	// 删除 comment
	err := lib.DB.Delete(&Comment{}, commentID).Error
	if err != nil {
		return err
	}

	return nil
}

func DelPage(page *Page) error {
	err := lib.DB.Unscoped().Delete(page).Error
	if err != nil {
		return err
	}

	// 删除所有相关内容
	var comments []Comment
	lib.DB.Where("page_key = ? AND site_name = ?", page.Key, page.SiteName).Find(&comments)

	for _, c := range comments {
		DelComment(c.ID)
	}

	// 删除 vote
	lib.DB.Unscoped().Where(
		"target_id = ? AND (type = ? OR type = ?)",
		page.ID,
		string(VoteTypePageUp),
		string(VoteTypePageDown),
	).Delete(&Vote{})

	return nil
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

func UserNotifyMarkAllAsRead(userID uint) error {
	if userID == 0 {
		return errors.New("user not found")
	}

	lib.DB.Model(&Notify{}).Where("user_id = ?", userID).Updates(&Notify{
		IsRead: true,
		ReadAt: time.Now(),
	})

	return nil
}
