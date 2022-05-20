package model

import (
	"strings"

	"github.com/ArtalkJS/ArtalkGo/config"
)

func CookAllComments(comments []Comment) []CookedComment {
	cookedComments := []CookedComment{}
	for _, c := range comments {
		cookedComments = append(cookedComments, c.ToCooked())
	}
	return cookedComments
}

func CookAllPages(pages []Page) []CookedPage {
	cookedPages := []CookedPage{}
	for _, p := range pages {
		cookedPages = append(cookedPages, p.ToCooked())
	}
	return cookedPages
}

func GetUserAllCommentIDs(userID uint) []uint {
	userAllCommentIDs := []uint{}
	DB().Model(&Comment{}).Select("id").Where("user_id = ?", userID).Find(&userAllCommentIDs)
	return userAllCommentIDs
}

func ContainsComment(comments []Comment, targetID uint) bool {
	for _, c := range comments {
		if c.ID == targetID {
			return true
		}
	}
	return false
}

func ContainsCookedComment(comments []CookedComment, targetID uint) bool {
	for _, c := range comments {
		if c.ID == targetID {
			return true
		}
	}
	return false
}

func SyncFromConf() {
	// 初始化默认站点
	FindCreateSite(config.Instance.SiteDefault)

	// 导入配置文件的管理员用户
	for _, admin := range config.Instance.AdminUsers {
		user := FindUser(admin.Name, admin.Email)
		receiveEmail := true // 默认允许接收邮件
		if admin.ReceiveEmail != nil {
			receiveEmail = *admin.ReceiveEmail
		}
		if user.IsEmpty() {
			// create
			user = User{
				Name:         admin.Name,
				Email:        admin.Email,
				Link:         admin.Link,
				Password:     admin.Password,
				BadgeName:    admin.BadgeName,
				BadgeColor:   admin.BadgeColor,
				IsAdmin:      true,
				IsInConf:     true,
				ReceiveEmail: receiveEmail,
				SiteNames:    strings.Join(admin.Sites, ","),
			}
			CreateUser(&user)
		} else {
			// update
			user.Name = admin.Name
			user.Email = admin.Email
			user.Link = admin.Link
			user.Password = admin.Password
			user.BadgeName = admin.BadgeName
			user.BadgeColor = admin.BadgeColor
			user.IsAdmin = true
			user.IsInConf = true
			user.ReceiveEmail = receiveEmail
			user.SiteNames = strings.Join(admin.Sites, ",")
			UpdateUser(&user)
		}
	}

	// 清理配置文件中不存在的用户
	// var dbAdminUsers []User
	// lib.DB.Model(&User{}).Where(&User{IsInConf: true}).Find(&dbAdminUsers)
	// for _, dbU := range dbAdminUsers {
	// 	isUserExist := func() bool {
	// 		for _, confU := range config.Instance.AdminUsers {
	// 			// 忽略大小写比较
	// 			if strings.EqualFold(confU.Name, dbU.Name) && strings.EqualFold(confU.Email, dbU.Email) {
	// 				return true
	// 			}
	// 		}
	// 		return false
	// 	}

	// 	if !isUserExist() {
	// 		DelUser(&dbU)
	// 	}
	// }
}
