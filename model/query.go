package model

import (
	"encoding/json"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/eko/gocache/v2/store"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func FindComment(id uint) Comment {
	var comment Comment
	lib.DB.Preload(clause.Associations).Where("id = ?", id).First(&comment)
	return comment
}

func FindCommentScopes(id uint, filters ...func(db *gorm.DB) *gorm.DB) Comment {
	var comment Comment
	lib.DB.Preload(clause.Associations).Where("id = ?", id).Scopes(filters...).First(&comment)
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

func IsAdminUser(name string, email string) bool {
	var user User
	lib.DB.Where("LOWER(name) = LOWER(?) AND LOWER(email) = LOWER(?) AND is_admin = 1", name, email).First(&user)
	return !user.IsEmpty()
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
