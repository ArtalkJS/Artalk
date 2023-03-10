package config

import (
	"strings"
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/sirupsen/logrus"
)

// 默认配置文件名
const DEFAULT_CONF_FILE = "artalk.yml"

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
	// load yaml config
	if err := kf.Load(file.Provider(cfgFile), parser); err != nil {
		logrus.Errorln(err)
		logrus.Fatal("Config file read error")
	}

	Instance = &Config{}

	if err := kf.Unmarshal("", Instance); err != nil {
		logrus.Errorln(err)
		logrus.Fatal("Config file parse error")
	}

	cfgFileLoaded = cfgFile

	// 后续处理
	postInit()
}

func postInit() {
	// 检查 app_key 是否设置
	if strings.TrimSpace(Instance.AppKey) == "" {
		logrus.Fatal("Please check config file and set an `app_key` for data encryption")
	}

	// 设置时区
	if strings.TrimSpace(Instance.TimeZone) == "" {
		logrus.Fatal("Please check config file and set `timezone`")
	}
	denverLoc, _ := time.LoadLocation(Instance.TimeZone)
	time.Local = denverLoc

	// 默认站点配置
	Instance.SiteDefault = strings.TrimSpace(Instance.SiteDefault)
	if Instance.SiteDefault == "" {
		logrus.Fatal("Please check config file and set `site_default`")
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
		logrus.Warn("The config option `captcha.action_timeout` is deprecated, please use `captcha.action_reset` instead")
		if Instance.Captcha.ActionReset == 0 {
			Instance.Captcha.ActionReset = Instance.Captcha.ActionTimeout
		}
	}
	if len(Instance.AllowOrigins) != 0 {
		logrus.Warn("The config option `allow_origins` is deprecated, please use `trusted_domains` instead")
		if len(Instance.TrustedDomains) == 0 {
			Instance.TrustedDomains = Instance.AllowOrigins
		}
	}

	// @version < 2.2.0
	if Instance.Notify != nil {
		logrus.Warn("The config option `notify` is deprecated, please use `admin_notify` instead")
		Instance.AdminNotify = *Instance.Notify
	}
	if Instance.AdminNotify.Email == nil {
		Instance.AdminNotify.Email = &AdminEmailConf{
			Enabled: true, // 默认开启管理员邮件通知
		}
	}
	if Instance.Email.MailSubjectToAdmin != "" {
		logrus.Warn("The config option `email.mail_subject_to_admin` is deprecated, please use `admin_notify.email.mail_subject` instead")
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

	// 默认将验证码类型设置为 image
	if strings.TrimSpace(string(Instance.Captcha.CaptchaType)) == "" {
		Instance.Captcha.CaptchaType = TypeImage
	}
}
