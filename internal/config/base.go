package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

func New() *Config {
	conf := &Config{}

	return conf
}

// 从文件中创建配置实例
func NewFromFile(cfgFile string) (*Config, error) {
	kf := koanf.New(".")

	// load yaml config
	if err := kf.Load(file.Provider(cfgFile), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("config file read error: %w", err)
	}

	// load environment variables and merge into the loaded config
	const envPrefix = "ATK_"
	if err := kf.Load(env.Provider(envPrefix, ".", func(s string) string {
		// FOO__BAR -> foo_bar to handle dash in config names
		s = strings.ToLower(strings.TrimPrefix(s, envPrefix))
		s = strings.ReplaceAll(s, "__", "-")
		s = strings.ReplaceAll(s, "_", ".")
		return strings.ReplaceAll(s, "-", "_")
	}), nil); err != nil {
		return nil, fmt.Errorf("config environment variable parse error: %w", err)
	}

	// create new config instance
	conf := &Config{
		cfgFile: cfgFile,
	}

	// use koanf parser to decode config file to instance
	if err := kf.Unmarshal("", conf); err != nil {
		return nil, fmt.Errorf("config file parse error: %w", err)
	}

	// patch config
	{
		conf.historyPatch()
		conf.normalPatch()
		conf.i18nPatch()
		conf.ipRegionPatch()
	}

	return conf, nil
}

func (conf *Config) GetCfgFileLoaded() string {
	return conf.cfgFile
}

// 配置修补
func (conf *Config) normalPatch() {
	// 检查 app_key 是否设置
	if strings.TrimSpace(conf.AppKey) == "" {
		conf.AppKey = utils.RandomString(16)
		log.Warn("config `app_key` is not set, now it is random value")
	}

	// 检查时区
	conf.TimeZone = strings.TrimSpace(conf.TimeZone)
	if conf.TimeZone == "" {
		conf.TimeZone = "Local"
		log.Warn("config `timezone` is not set, now it is: " + strconv.Quote(conf.TimeZone))
	}

	// 默认站点配置
	conf.SiteDefault = strings.TrimSpace(conf.SiteDefault)
	if conf.SiteDefault == "" {
		conf.SiteDefault = "Default Site"
		log.Warn("config `site_default` is not set, now it is: " + strconv.Quote(conf.SiteDefault))
	}

	// 缓存配置
	if conf.Cache.Type == "" {
		// 默认使用内建缓存
		conf.Cache.Type = CacheTypeBuiltin
	}

	// 配置文件 alias 处理
	if conf.Captcha.ActionLimit == 0 {
		conf.Captcha.Always = true
	}

	// 管理员邮件通知配置继承
	if conf.AdminNotify.Email.MailSubject == "" {
		if conf.AdminNotify.NotifySubject != "" {
			conf.AdminNotify.Email.MailSubject = conf.AdminNotify.NotifySubject
		} else if conf.Email.MailSubject != "" {
			conf.AdminNotify.Email.MailSubject = conf.Email.MailSubject
		}
	}

	// 默认待审模式下开启管理员通知嘈杂模式，保证管理员能看到待审核文章
	if conf.Moderator.PendingDefault {
		conf.AdminNotify.NoiseMode = true
	}

	// 默认将验证码类型设置为 image
	if strings.TrimSpace(string(conf.Captcha.CaptchaType)) == "" {
		conf.Captcha.CaptchaType = TypeImage
	}

	// 图片上传存放路径默认设置
	if conf.ImgUpload.Path == "" {
		conf.ImgUpload.Path = "./data/artalk-img/"
		log.Warn("[Image Upload] img_upload.path is not configured, using the default value: " + strconv.Quote(conf.ImgUpload.Path))
	}
}

// 多语言配置修补
func (conf *Config) i18nPatch() {
	if conf.Locale == "" {
		conf.Locale = "en"

		// zh-CN default patch (for 历史兼容)
		// 判断配置文件中是否有中文，若有中文则将 locale 设置为 zh-CN
		if confRaw, err := os.ReadFile(conf.GetCfgFileLoaded()); err == nil {
			containsHan := false
			for _, runeValue := range string(confRaw) {
				if unicode.Is(unicode.Han, runeValue) {
					containsHan = true
					break
				}
			}
			if containsHan {
				conf.Locale = "zh-CN"
			}
		}

		log.Warn("config `locale` is not set, now it is: " + strconv.Quote(conf.Locale))
	} else if conf.Locale == "zh" {
		conf.Locale = "zh-CN"
	}

	// Case Convert
	// follow Unicode BCP 47
	// @see https://www.techonthenet.com/js/language_tags.php
	parts := strings.Split(conf.Locale, "-")
	if len(parts) == 2 {
		conf.Locale = strings.ToLower(parts[0]) + "-" + strings.ToUpper(parts[1])
	} else if len(parts) == 1 {
		conf.Locale = strings.ToLower(parts[0])
	}

	// Temporary convert `en-*` to `en`
	if parts[0] == "en" {
		conf.Locale = "en"
	}
}

// IP属地功能配置修补
func (conf *Config) ipRegionPatch() {
	// IP 属地默认数据文件
	if conf.IPRegion.DBPath == "" {
		conf.IPRegion.DBPath = "./data/ip2region.xdb"
	}

	// 检测配置文件是否存在
	if !utils.CheckFileExist(conf.IPRegion.DBPath) {
		log.Warn("未找到 IP 数据库文件：" + strconv.Quote(conf.IPRegion.DBPath) + "，IP 属地功能已禁用，" +
			"参考链接：https://artalk.js.org/guide/frontend/ip-region.html")
		conf.IPRegion.Enabled = false
	}

	// 默认精确到省
	if conf.IPRegion.Precision == "" {
		conf.IPRegion.Precision = string(IPRegionProvince)
	}
}

// 配置修补 for 历史版本
func (conf *Config) historyPatch() {
	if conf.Captcha.ActionTimeout != 0 {
		log.Warn("The config option `captcha.action_timeout` is deprecated, please use `captcha.action_reset` instead")
		if conf.Captcha.ActionReset == 0 {
			conf.Captcha.ActionReset = conf.Captcha.ActionTimeout
		}
	}
	if len(conf.AllowOrigins) != 0 {
		log.Warn("The config option `allow_origins` is deprecated, please use `trusted_domains` instead")
		if len(conf.TrustedDomains) == 0 {
			conf.TrustedDomains = conf.AllowOrigins
		}
	}

	// @version < 2.2.0
	if conf.Notify != nil {
		log.Warn("The config option `notify` is deprecated, please use `admin_notify` instead")
		conf.AdminNotify = *conf.Notify
	}
	if conf.AdminNotify.Email == nil {
		conf.AdminNotify.Email = &AdminEmailConf{
			Enabled: true, // 默认开启管理员邮件通知
		}
	}
	if conf.Email.MailSubjectToAdmin != "" {
		log.Warn("The config option `email.mail_subject_to_admin` is deprecated, please use `admin_notify.email.mail_subject` instead")
		conf.AdminNotify.Email.MailSubject = conf.Email.MailSubjectToAdmin
	}
}

// 尝试查找配置文件
func RetrieveConfigFile() string {
	for _, v := range CONF_DEFAULT_FILENAMES {
		if utils.CheckFileExist(v) {
			return v
		}
	}
	return ""
}
