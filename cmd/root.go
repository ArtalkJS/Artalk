package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/lib/email"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var Banner = `
 ________  ________  _________  ________  ___       ___  __       
|\   __  \|\   __  \|\___   ___\\   __  \|\  \     |\  \|\  \     
\ \  \|\  \ \  \|\  \|___ \  \_\ \  \|\  \ \  \    \ \  \/  /|_   
 \ \   __  \ \   _  _\   \ \  \ \ \   __  \ \  \    \ \   ___  \  
  \ \  \ \  \ \  \\  \|   \ \  \ \ \  \ \  \ \  \____\ \  \\ \  \ 
   \ \__\ \__\ \__\\ _\    \ \__\ \ \__\ \__\ \_______\ \__\\ \__\
    \|__|\|__|\|__|\|__|    \|__|  \|__|\|__|\|_______|\|__| \|__|
 
Artalk: A Fast, Slight & Delightful Comment System.
More details on https://github.com/ArtalkJS/Artalk
(c) 2021 artalk.js.org`

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "artalk-go",
	Short: "Artalk: A Fast, Slight & Delightful Comment System",
	Long:  Banner,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Banner)
		fmt.Println()
		fmt.Println("NOTE: add `-h` flag to show help about any command.")
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initLog)
	cobra.OnInitialize(initDB)
	cobra.OnInitialize(syncConfWithDB)
	cobra.OnInitialize(initCache)
	cobra.OnInitialize(email.InitQueue) // 初始化邮件队列

	rootCmd.SetVersionTemplate("Artalk-GO {{printf \"version %s\" .Version}}\n")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件路径 (defaults are './artalk-go.yml')")
}

// 1. 初始化配置
func initConfig() {
	config.Init(cfgFile)

	// 设置时区
	if strings.TrimSpace(config.Instance.TimeZone) == "" {
		logrus.Fatal("请检查配置文件，并设置 timezone")
		os.Exit(1)
	}
	denverLoc, _ := time.LoadLocation(config.Instance.TimeZone)
	time.Local = denverLoc

	// 检查 app_key 是否设置
	if strings.TrimSpace(config.Instance.AppKey) == "" {
		logrus.Fatal("请检查配置文件，并设置一个 app_key (任意字符串) 用于数据加密")
		os.Exit(1)
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
	err := lib.OpenDB()
	if err != nil {
		logrus.Error("数据库初始化发生错误 ", err)
		os.Exit(1)
	}

	// Migrate the schema
	lib.DB.AutoMigrate(&model.Site{}, &model.Page{}, &model.User{}, &model.Comment{}) // 注意表的创建顺序，因为有关联字段
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
	lib.DB.Where("is_in_conf = 1").Find(&dbAdminUsers)
	for _, dbU := range dbAdminUsers {
		isUserExist := func() bool {
			for _, confU := range config.Instance.AdminUsers {
				if confU.Name == dbU.Name && confU.Email == dbU.Email {
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
