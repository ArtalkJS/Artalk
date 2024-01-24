package core

import (
	"github.com/ArtalkJS/Artalk/internal/entity"
)

func (app *App) syncFromConf() {
	// 初始化默认站点
	app.Dao().FindCreateSite(app.Conf().SiteDefault)

	// 导入配置文件的管理员用户
	for _, admin := range app.Conf().AdminUsers {
		user := app.Dao().FindUser(admin.Name, admin.Email)
		receiveEmail := true // 默认允许接收邮件
		if admin.ReceiveEmail != nil {
			receiveEmail = *admin.ReceiveEmail
		}
		if user.IsEmpty() {
			// create
			user = entity.User{
				Name:         admin.Name,
				Email:        admin.Email,
				Link:         admin.Link,
				Password:     admin.Password,
				BadgeName:    admin.BadgeName,
				BadgeColor:   admin.BadgeColor,
				IsAdmin:      true,
				IsInConf:     true,
				ReceiveEmail: receiveEmail,
			}
			app.dao.CreateUser(&user)
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
			app.dao.UpdateUser(&user)
		}
	}

	// 清理配置文件中不存在的用户
	// var dbAdminUsers []User
	// lib.DB.Model(&User{}).Where(&User{IsInConf: true}).Find(&dbAdminUsers)
	// for _, dbU := range dbAdminUsers {
	// 	isUserExist := func() bool {
	// 		for _, confU := range app.Conf().AdminUsers {
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
