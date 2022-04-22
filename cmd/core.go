package cmd

import (
	"fmt"
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
)

// 装载核心功能
func loadCore() {
	initConfig()
	initLog()
	initDB()
	initCache()
	syncConfWithDB()
	notify_launcher.Init() // 初始化 Notify 发射台

	makeCache()
	// TODO 异步加载开关
	// go func() {
	// 	makeCache()
	// }()
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

	// 配置文件 alias 处理
	if config.Instance.Captcha.ActionLimit == 0 {
		config.Instance.Captcha.Always = true
	}

	// 检查废弃需更新配置
	if config.Instance.Captcha.ActionTimeout != 0 {
		logrus.Warn("captcha.action_timeout 配置项已废弃，请使用 captcha.action_reset 代替")
		if config.Instance.Captcha.ActionReset == 0 {
			config.Instance.Captcha.ActionReset = config.Instance.Captcha.ActionTimeout
		}
	}
	if len(config.Instance.AllowOrigins) != 0 {
		logrus.Warn("allow_origins 配置项已废弃，请使用 trusted_domains 代替")
		if len(config.Instance.TrustedDomains) == 0 {
			config.Instance.TrustedDomains = config.Instance.AllowOrigins
		}
	}
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
	model.InitDB()
}

// 4. 初始化缓存
func initCache() {
	err := lib.OpenCache()
	if err != nil {
		logrus.Error("缓存初始化发生错误 ", err)
		os.Exit(1)
	}
}

// 5. 同步配置文件与数据库
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
			model.CreateUser(&user)
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
			model.UpdateUser(&user)
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
			model.DelUser(&dbU)
		}
	}
}

// 制备缓存
func makeCache() {
	// Users
	{
		start := time.Now()

		var items []model.User
		lib.DB.Find(&items)

		for _, item := range items {
			model.UserCacheSave(&item)
		}

		logrus.Debug(fmt.Sprintf("[Users] 缓存完毕 (共 %d 个，耗时：%s)", len(items), time.Since(start)))
	}

	// Sites
	{
		start := time.Now()

		var items []model.Site
		lib.DB.Find(&items)

		for _, item := range items {
			model.SiteCacheSave(&item)
		}

		logrus.Debug(fmt.Sprintf("[Sites] 缓存完毕 (共 %d 个，耗时：%s)", len(items), time.Since(start)))
	}

	// Pages
	{
		start := time.Now()

		var items []model.Page
		lib.DB.Find(&items)

		for _, item := range items {
			model.PageCacheSave(&item)
		}

		logrus.Debug(fmt.Sprintf("[Pages] 缓存完毕 (共 %d 个，耗时：%s)", len(items), time.Since(start)))
	}

	// Comments
	{
		start := time.Now()

		var items []model.Comment
		lib.DB.Find(&items)

		for _, item := range items {
			model.CommentCacheSave(&item)
		}

		logrus.Debug(fmt.Sprintf("[Comments] 缓存完毕 (共 %d 个，耗时：%s)", len(items), time.Since(start)))
	}
}
