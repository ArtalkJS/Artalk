package cmd

import (
	"os"
	"strings"
	"time"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/ArtalkJS/ArtalkGo/model/notify_launcher"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"gorm.io/gorm"
)

// 装载核心功能
func loadCore() {
	initConfig()
	initLog()
	initDB()
	syncConfWithDB()
	initCache()
	notify_launcher.Init() // 初始化 Notify 发射台
}

// 1. 初始化配置
func initConfig() {
	config.Init(cfgFile, workDir)

	// 检查 app_key 是否设置
	if strings.TrimSpace(config.Instance.AppKey) == "" {
		logrus.Fatal("请检查配置文件，并设置一个 app_key (任意字符串) 用于数据加密")
	}

	// 设置时区
	if strings.TrimSpace(config.Instance.TimeZone) == "" {
		logrus.Fatal("请检查配置文件，并设置 timezone")
	}
	denverLoc, _ := time.LoadLocation(config.Instance.TimeZone)
	time.Local = denverLoc
}

// 2. 初始化日志
func initLog() {
	if !config.Instance.Log.Enabled {
		return
	}

	// 命令行输出格式
	stdFormatter := &prefixed.TextFormatter{
		DisableTimestamp: true,
		ForceFormatting:  true,
		ForceColors:      true,
		DisableColors:    false,
	}

	// 文件输出格式
	fileFormatter := &prefixed.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02.15:04:05.000000",
		ForceFormatting: true,
		ForceColors:     false,
		DisableColors:   true,
	}

	// logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(stdFormatter)
	logrus.SetOutput(os.Stdout)

	if config.Instance.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if config.Instance.Log.Filename != "" {
		// 文件保存
		pathMap := lfshook.PathMap{
			logrus.InfoLevel:  config.Instance.Log.Filename,
			logrus.DebugLevel: config.Instance.Log.Filename,
			logrus.ErrorLevel: config.Instance.Log.Filename,
		}
		logrus.AddHook(lfshook.NewHook(
			pathMap,
			fileFormatter,
		))
	}
}

// 3. 初始化数据库
func initDB() {
	var db *gorm.DB
	db, err := lib.OpenDB(config.DBType(config.Instance.DB.Type), config.Instance.DB.Dsn)
	if err != nil {
		logrus.Error("数据库初始化发生错误 ", err)
		os.Exit(1)
	}

	lib.DB = db

	// Migrate the schema
	lib.DB.AutoMigrate(&model.Site{}, &model.Page{}, &model.User{},
		&model.Comment{}, &model.Notify{}, &model.Vote{}, &model.PV{}) // 注意表的创建顺序，因为有关联字段
}

// 4. 同步配置文件与数据库
func syncConfWithDB() {
	// 初始化默认站点
	siteDefault := strings.TrimSpace(config.Instance.SiteDefault)
	if siteDefault == "" {
		logrus.Error("请设置 SiteDefault 默认站点，不能为空")
		os.Exit(1)
	}
	model.FindCreateSite(siteDefault)

	// 导入配置文件的管理员用户
	for _, admin := range config.Instance.AdminUsers {
		user := model.FindUser(admin.Name, admin.Email)
		if user.IsEmpty() {
			// create
			user = model.User{
				Name:       admin.Name,
				Email:      admin.Email,
				Link:       admin.Link,
				Password:   admin.Password,
				BadgeName:  admin.BadgeName,
				BadgeColor: admin.BadgeColor,
				IsAdmin:    true,
				IsInConf:   true,
			}
			lib.DB.Create(&user)
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
			lib.DB.Save(&user)
		}
	}

	// 清理配置文件中不存在的用户
	var dbAdminUsers []model.User
	lib.DB.Model(&model.User{}).Where(&model.User{IsInConf: true}).Find(&dbAdminUsers)
	for _, dbU := range dbAdminUsers {
		isUserExist := func() bool {
			for _, confU := range config.Instance.AdminUsers {
				// 忽略大小写比较
				if strings.EqualFold(confU.Name, dbU.Name) && strings.EqualFold(confU.Email, dbU.Email) {
					return true
				}
			}
			return false
		}

		if !isUserExist() {
			lib.DB.Unscoped().Delete(&dbU)
		}
	}
}

// 5. 初始化缓存
func initCache() {
	err := lib.OpenCache()
	if err != nil {
		logrus.Error("缓存初始化发生错误 ", err)
		os.Exit(1)
	}
}
