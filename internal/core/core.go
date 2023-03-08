package core

import (
	"io"
	"os"
	"strconv"
	"sync"
	"unicode"

	"github.com/ArtalkJS/Artalk/internal/cache"
	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/db"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/ip_region"
	"github.com/ArtalkJS/Artalk/internal/notify_launcher"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var firstLoad = true
var mutex = sync.Mutex{}

// 装载核心功能
func LoadCore(cfgFile string, workDir string) {
	mutex.Lock()
	defer mutex.Unlock()

	firstLoad = false

	initConfig(cfgFile, workDir)
	initI18n()
	initLog()
	initCache()
	initDB()
	notify_launcher.Init() // 初始化 Notify 发射台
	initIPRegion()

	// 首次 Load
	if firstLoad {
		// 缓存预热
		if config.Instance.Cache.Enabled && config.Instance.Cache.WarmUp {
			cache.CacheWarmUp()
		}
		// 异步加载
		// go func() {
		// 	entity.CacheWarmUp()
		// }()
	}
}

// 仅装载配置
func LoadConfOnly(cfgFile string, workDir string) {
	initConfig(cfgFile, workDir)
}

// 1. 初始化配置
func initConfig(cfgFile string, workDir string) {
	// 切换工作目录
	if workDir != "" {
		if err := os.Chdir(workDir); err != nil {
			logrus.Fatal("Working directory change error: ", err)
		}
	}

	if cfgFile == "" {
		cfgFile = config.DEFAULT_CONF_FILE

		// 默认配置文件名 "artalk-go.yml"（for 向下兼容）
		if _, err := os.Stat("artalk-go.yml"); err == nil {
			cfgFile = "artalk-go.yml"
		}
	}

	// 自动生成新配置文件
	if !CheckFileExist(cfgFile) {
		Gen("config", cfgFile, false)
	}

	config.Init(cfgFile)
}

func initI18n() {
	if config.Instance.Locale == "" {
		config.Instance.Locale = "en"

		// zh-CN default patch
		// 判断配置文件中是否有中文，若有中文则将 locale 设置为 zh-CN
		if confRaw, err := os.ReadFile(config.GetCfgFileLoaded()); err == nil {
			containsHan := false
			for _, runeValue := range string(confRaw) {
				if unicode.Is(unicode.Han, runeValue) {
					containsHan = true
					break
				}
			}
			if containsHan {
				config.Instance.Locale = "zh-CN"
			}
		}

		logrus.Warn("config `locale` is not set, now it is: " + strconv.Quote(config.Instance.Locale))
	} else if config.Instance.Locale == "zh" {
		config.Instance.Locale = "zh-CN"
	}

	i18n.Init(config.Instance.Locale)
}

// 2. 初始化日志
func initLog() {
	logrus.New()
	if !config.Instance.Log.Enabled {
		logrus.SetOutput(io.Discard)
		return
	}

	// 命令行输出格式
	stdFormatter := &prefixed.TextFormatter{
		DisableTimestamp: true,
		ForceFormatting:  true,
		ForceColors:      true,
		DisableColors:    false,
	}

	logrus.SetFormatter(stdFormatter)
	logrus.SetOutput(os.Stdout)

	if config.Instance.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	// 日志输出到文件
	if config.Instance.Log.Filename != "" {
		fileFormatter := &prefixed.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02.15:04:05.000000",
			ForceFormatting: true,
			ForceColors:     false,
			DisableColors:   true,
		}

		pathMap := lfshook.PathMap{
			logrus.InfoLevel:  config.Instance.Log.Filename,
			logrus.DebugLevel: config.Instance.Log.Filename,
			logrus.ErrorLevel: config.Instance.Log.Filename,
		}

		newHooks := make(logrus.LevelHooks)
		newHooks.Add(lfshook.NewHook(
			pathMap,
			fileFormatter,
		))

		//logrus.AddHook(lfshook.NewHook()) // 使用 Replace 而不使用 Add
		logrus.StandardLogger().ReplaceHooks(newHooks)
	}
}

// 3. 初始化缓存
func initCache() {
	err := cache.OpenCache()
	if err != nil {
		logrus.Error("[Cache] ", "Init cache error: ", err)
		os.Exit(1)
	}
}

// 4. 初始化数据库
func initDB() {
	db.InitDB()
	db.MigrateModels()
	SyncFromConf()
}

func initIPRegion() {
	irConf := &config.Instance.IPRegion
	if !irConf.Enabled {
		return
	}

	if irConf.DBPath == "" {
		irConf.DBPath = "./data/ip2region.xdb"
	}
	if irConf.Precision == "" {
		// 默认精确到省
		irConf.Precision = string(ip_region.Province)
	}

	ip_region.Init(irConf.DBPath)                                 // 初始化 IP 归属地查询
	ip_region.SetPrecision(ip_region.Precision(irConf.Precision)) // 精度设置
}
