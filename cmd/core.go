package cmd

import (
	"os"

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
	initCache()
	initDB()
	notify_launcher.Init() // 初始化 Notify 发射台

	// 缓存预热
	if config.Instance.Cache.Enabled && config.Instance.Cache.WarmUp {
		model.CacheWarmUp()
	}
	// 异步加载
	// go func() {
	// 	model.CacheWarmUp()
	// }()
}

// 1. 初始化配置
func initConfig() {
	config.Init(cfgFile, workDir)
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

// 3. 初始化缓存
func initCache() {
	err := lib.OpenCache()
	if err != nil {
		logrus.Error("缓存初始化发生错误 ", err)
		os.Exit(1)
	}
}

// 4. 初始化数据库
func initDB() {
	lib.InitDB()
	model.SetDB(lib.DB)
	model.MigrateModels()
	model.SyncFromConf()
}
