package config

import (
	"strings"
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/sirupsen/logrus"
)

const DEFAULT_CONF_FILE = "artalk-go.yml"

var (
	kf     = koanf.New(".")
	parser = yaml.Parser()

	// 配置实例
	Instance      *Config
	cfgFileLoaded string
)

func GetCfgFileLoaded() string {
	return cfgFileLoaded
}

// Init 初始化配置
func Init(cfgFile string) {
	if cfgFile == "" {
		cfgFile = DEFAULT_CONF_FILE
	}

	// load yaml config
	if err := kf.Load(file.Provider(cfgFile), parser); err != nil {
		logrus.Errorln(err)
		logrus.Fatal("配置文件读取错误")
	}

	Instance = &Config{}

	if err := kf.Unmarshal("", Instance); err != nil {
		logrus.Errorln(err)
		logrus.Fatal("配置文件解析错误")
	}

	cfgFileLoaded = cfgFile

	// 后续处理
	postInit()
}

func postInit() {
	// 检查 app_key 是否设置
	if strings.TrimSpace(Instance.AppKey) == "" {
		logrus.Fatal("请检查配置文件，并设置一个 app_key (任意字符串) 用于数据加密")
	}

	// 设置时区
	if strings.TrimSpace(Instance.TimeZone) == "" {
		logrus.Fatal("请检查配置文件，并设置 timezone")
	}
	denverLoc, _ := time.LoadLocation(Instance.TimeZone)
	time.Local = denverLoc

	// 默认站点配置
	Instance.SiteDefault = strings.TrimSpace(Instance.SiteDefault)
	if Instance.SiteDefault == "" {
		logrus.Fatal("请设置 SiteDefault 默认站点，不能为空")
	}

	// 缓存配置
	if Instance.Cache.Type == "" {
		// 默认使用内建缓存
		Instance.Cache.Type = CacheTypeBuiltin
	}
	if Instance.Cache.Type != CacheTypeDisabled {
		// 非缓存禁用模式，Enabled = true
		Instance.Cache.Enabled = true
	}

	// 配置文件 alias 处理
	if Instance.Captcha.ActionLimit == 0 {
		Instance.Captcha.Always = true
	}

	/* 检查废弃需更新配置 */
	if Instance.Captcha.ActionTimeout != 0 {
		logrus.Warn("captcha.action_timeout 配置项已废弃，请使用 captcha.action_reset 代替")
		if Instance.Captcha.ActionReset == 0 {
			Instance.Captcha.ActionReset = Instance.Captcha.ActionTimeout
		}
	}
	if len(Instance.AllowOrigins) != 0 {
		logrus.Warn("allow_origins 配置项已废弃，请使用 trusted_domains 代替")
		if len(Instance.TrustedDomains) == 0 {
			Instance.TrustedDomains = Instance.AllowOrigins
		}
	}

	// @version < 2.2.0
	if Instance.Notify != nil {
		logrus.Warn("notify 配置项已废弃，请使用 admin_notify 代替")
		Instance.AdminNotify = *Instance.Notify
	}
	if Instance.AdminNotify.Email == nil {
		Instance.AdminNotify.Email = &AdminEmailConf{
			Enabled: true, // 默认开启管理员邮件通知
		}
	}
	if Instance.Email.MailSubjectToAdmin != "" {
		logrus.Warn("email.mail_subject_to_admin 配置项已废弃，请使用 admin_notify.email.mail_subject 代替")
		Instance.AdminNotify.Email.MailSubject = Instance.Email.MailSubjectToAdmin
	}

	// 管理员邮件通知配置继承
	if Instance.AdminNotify.Email.MailSubject == "" {
		if Instance.AdminNotify.NotifySubject != "" {
			Instance.AdminNotify.Email.MailSubject = Instance.AdminNotify.NotifySubject
		} else if Instance.Email.MailSubject != "" {
			Instance.AdminNotify.Email.MailSubject = Instance.Email.MailSubject
		}
	}

	// 默认待审模式下开启管理员通知嘈杂模式，保证管理员能看到待审核文章
	if Instance.Moderator.PendingDefault {
		Instance.AdminNotify.NoiseMode = true
	}
}
