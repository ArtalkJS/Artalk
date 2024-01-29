package config

import (
	"time"
)

// Config 配置
// @link https://github.com/knadh/koanf
type Config struct {
	AppKey         string                 `koanf:"app_key" json:"app_key"`                 // 加密密钥
	Debug          bool                   `koanf:"debug" json:"debug"`                     // 调试模式
	Locale         string                 `koanf:"locale" json:"locale"`                   // 语言
	TimeZone       string                 `koanf:"timezone" json:"timezone"`               // 时区
	Host           string                 `koanf:"host" json:"host"`                       // HTTP Server 监听 IP
	Port           int                    `koanf:"port" json:"port"`                       // HTTP Server 监听 Port
	DB             DBConf                 `koanf:"db" json:"db"`                           // 数据库配置
	Cache          CacheConf              `koanf:"cache" json:"cache"`                     // 缓存
	Log            LogConf                `koanf:"log" json:"log"`                         // 日志文件
	TrustedDomains []string               `koanf:"trusted_domains" json:"trusted_domains"` // 可信任的域名 (新)
	SSL            SSLConf                `koanf:"ssl" json:"ssl"`                         // SSL
	SiteDefault    string                 `koanf:"site_default" json:"site_default"`       // 默认站点名（当请求无指定 site_name 时使用）
	AdminUsers     []AdminUserConf        `koanf:"admin_users" json:"admin_users"`         // 管理员账户
	LoginTimeout   int                    `koanf:"login_timeout" json:"login_timeout"`     // 登录超时
	Moderator      ModeratorConf          `koanf:"moderator" json:"moderator"`             // 评论审查
	Captcha        CaptchaConf            `koanf:"captcha" json:"captcha"`                 // 验证码
	Email          EmailConf              `koanf:"email" json:"email"`                     // 邮箱提醒
	IPRegion       IPRegionConf           `koanf:"ip_region" json:"ip_region"`             // IP 归属地展示
	ImgUpload      ImgUploadConf          `koanf:"img_upload" json:"img_upload"`           // 图片上传
	AdminNotify    AdminNotifyConf        `koanf:"admin_notify" json:"admin_notify"`       // 其他通知方式
	Frontend       map[string]interface{} `koanf:"frontend" json:"frontend"`

	// deprecated options
	// (only for unmarshal, please not reference)
	// ---------------------------

	AllowOrigins []string         `koanf:"allow_origins" json:"-"` // @deprecated 已废弃 (请使用 TrustedDomains)
	Notify       *AdminNotifyConf `koanf:"notify" json:"-"`        // @deprecated 已废弃 (请使用 AdminNotify)

	// system runtime produce data
	// ---------------------------

	cfgFile string
}

type DBConf struct {
	Type DBType `koanf:"type" json:"type"`
	Dsn  string `koanf:"dsn" json:"dsn"` // 最高优先级

	File string `koanf:"file" json:"file"`
	Name string `koanf:"name" json:"name"`

	Host     string `koanf:"host" json:"host"`
	Port     int    `koanf:"port" json:"port"`
	User     string `koanf:"user" json:"user"`
	Password string `koanf:"password" json:"password"`

	TablePrefix string `koanf:"table_prefix" json:"table_prefix"`
	Charset     string `koanf:"charset" json:"charset"`
	SSL         bool   `koanf:"ssl" json:"ssl"`
}

type CacheConf struct {
	Enabled bool      `koanf:"enabled" json:"enabled"`
	Type    CacheType `koanf:"type" json:"type"`
	Expires int       `koanf:"expires" json:"expires"` // 过期时间
	WarmUp  bool      `koanf:"warm_up" json:"warm_up"` // 启动时缓存预热
	Server  string    `koanf:"server" json:"server"`   // 缓存服务器
	Redis   RedisConf `koanf:"redis" json:"redis"`
}

func (c *CacheConf) GetExpiresTime() int64 {
	if c.Expires == 0 {
		return int64(30 * time.Minute) // 默认 30min
	}

	if c.Expires == -1 {
		return -1 // Redis.KeepTTL = -1
	}

	return int64(time.Duration(c.Expires) * time.Minute)
}

type LogConf struct {
	Enabled  bool   `koanf:"enabled" json:"enabled"`
	Filename string `koanf:"filename" json:"filename"`
}

type SSLConf struct {
	Enabled  bool   `koanf:"enabled" json:"enabled"`
	CertPath string `koanf:"cert_path" json:"cert_path"`
	KeyPath  string `koanf:"key_path" json:"key_path"`
}

type AdminUserConf struct {
	Name         string `koanf:"name" json:"name"`
	Email        string `koanf:"email" json:"email"`
	Link         string `koanf:"link" json:"link"`
	Password     string `koanf:"password" json:"password"`
	BadgeName    string `koanf:"badge_name" json:"badge_name"`
	BadgeColor   string `koanf:"badge_color" json:"badge_color"`
	ReceiveEmail *bool  `koanf:"receive_email" json:"receive_email"`
}

type ModeratorConf struct {
	PendingDefault bool                 `koanf:"pending_default" json:"pending_default"`
	ApiFailBlock   bool                 `koanf:"api_fail_block" json:"api_fail_block"` // API 请求错误仍然拦截
	AkismetKey     string               `koanf:"akismet_key" json:"akismet_key"`
	Tencent        TencentAntispamConf  `koanf:"tencent" json:"tencent"`
	Aliyun         AliyunAntispamConf   `koanf:"aliyun" json:"aliyun"`
	Keywords       KeyWordsAntispamConf `koanf:"keywords" json:"keywords"`
}

// 腾讯云反垃圾
type TencentAntispamConf struct {
	Enabled   bool   `koanf:"enabled" json:"enabled"`
	SecretID  string `koanf:"secret_id" json:"secret_id"`
	SecretKey string `koanf:"secret_key" json:"secret_key"`
	Region    string `koanf:"region" json:"region"`
}

// 阿里云反垃圾
type AliyunAntispamConf struct {
	Enabled         bool   `koanf:"enabled" json:"enabled"`
	AccessKeyID     string `koanf:"access_key_id" json:"access_key_id"`
	AccessKeySecret string `koanf:"access_key_secret" json:"access_key_secret"`
	Region          string `koanf:"region" json:"region"`
}

// 关键词词库过滤
type KeyWordsAntispamConf struct {
	Enabled  bool     `koanf:"enabled" json:"enabled"`
	Pending  bool     `koanf:"pending" json:"pending"`
	Files    []string `koanf:"files" json:"files"`
	FileSep  string   `koanf:"file_sep" json:"file_sep"`
	ReplacTo string   `koanf:"replac_to" json:"replac_to"`
}

type CaptchaConf struct {
	Enabled       bool          `koanf:"enabled" json:"enabled"`
	Always        bool          `koanf:"always" json:"always"`
	CaptchaType   CaptchaType   `koanf:"captcha_type" json:"captcha_type"`
	ActionTimeout int           `koanf:"action_timeout" json:"-"` // @deprecated 已废弃 (请使用 ActionReset)
	ActionReset   int           `koanf:"action_reset" json:"action_reset"`
	ActionLimit   int           `koanf:"action_limit" json:"action_limit"`
	Turnstile     TurnstileConf `koanf:"turnstile" json:"turnstile"`
	ReCaptcha     ReCaptchaConf `koanf:"recaptcha" json:"recaptcha"`
	HCaptcha      HCaptchaConf  `koanf:"hcaptcha" json:"hcaptcha"`
	Geetest       GeetestConf   `koanf:"geetest" json:"geetest"`
}

type CaptchaType string

const (
	TypeImage     CaptchaType = "image"
	TypeTurnstile CaptchaType = "turnstile"
	TypeReCaptcha CaptchaType = "recaptcha"
	TypeHCaptcha  CaptchaType = "hcaptcha"
	TypeGeetest   CaptchaType = "geetest"
)

type TurnstileConf struct {
	SiteKey   string `koanf:"site_key" json:"site_key"`
	SecretKey string `koanf:"secret_key" json:"secret_key"`
}

type ReCaptchaConf struct {
	SiteKey   string `koanf:"site_key" json:"site_key"`
	SecretKey string `koanf:"secret_key" json:"secret_key"`
}

type HCaptchaConf struct {
	SiteKey   string `koanf:"site_key" json:"site_key"`
	SecretKey string `koanf:"secret_key" json:"secret_key"`
}

type GeetestConf struct {
	Enabled    bool   `koanf:"enabled" json:"-"` // @deprecated 已废弃 (请使用 captcha.captcha_type)
	CaptchaID  string `koanf:"captcha_id" json:"captcha_id"`
	CaptchaKey string `koanf:"captcha_key" json:"captcha_key"`
}

type EmailConf struct {
	Enabled            bool            `koanf:"enabled" json:"enabled"`           // 总开关
	SendType           EmailSenderType `koanf:"send_type" json:"send_type"`       // 发送方式
	SendName           string          `koanf:"send_name" json:"send_name"`       // 发件人名
	SendAddr           string          `koanf:"send_addr" json:"send_addr"`       // 发件人地址
	MailSubject        string          `koanf:"mail_subject" json:"mail_subject"` // 邮件标题
	MailSubjectToAdmin string          `koanf:"mail_subject_to_admin" json:"-"`   // @deprecated 已废弃 (请使用 AdminNotify.Email.MailSubject) - 邮件标题 (发送给管理员用)
	MailTpl            string          `koanf:"mail_tpl" json:"mail_tpl"`         // 邮件模板
	SMTP               SMTPConf        `koanf:"smtp" json:"smtp"`                 // SMTP 配置
	AliDM              AliDMConf       `koanf:"ali_dm" json:"ali_dm"`             // 阿里云邮件配置
	Queue              EmailQueueConf  `koanf:"queue" json:"queue"`               // 邮件发送队列配置
}

type SMTPConf struct {
	Host     string `koanf:"host" json:"host"`
	Port     int    `koanf:"port" json:"port"`
	Username string `koanf:"username" json:"username"`
	Password string `koanf:"password" json:"password"`
	From     string `koanf:"from" json:"from"`
}

type AliDMConf struct {
	AccessKeyId     string `koanf:"access_key_id" json:"access_key_id"`
	AccessKeySecret string `koanf:"access_key_secret" json:"access_key_secret"`
	AccountName     string `koanf:"account_name" json:"account_name"`
	Region          string `koanf:"region" json:"region"`
}

type EmailQueueConf struct {
	BufferSize int `koanf:"buffer_size" json:"buffer_size"` // Channel buffer size (default is zero that not create buffer)
}

type DBType string

const (
	TypeMySql      DBType = "mysql"
	TypeSQLite     DBType = "sqlite"
	TypePostgreSQL DBType = "pgsql"
	TypeMSSQL      DBType = "mssql"
)

type CacheType string

const (
	CacheTypeBuiltin  CacheType = "builtin" // 内建缓存
	CacheTypeRedis    CacheType = "redis"
	CacheTypeMemcache CacheType = "memcache"
)

type EmailSenderType string

const (
	TypeSMTP     EmailSenderType = "smtp"
	TypeAliDM    EmailSenderType = "ali_dm"
	TypeSendmail EmailSenderType = "sendmail"
)

// # Redis 配置
// redis:
//
//	network: "tcp"
//	username: ""
//	password: ""
//	db: 0
type RedisConf struct {
	Network  string `koanf:"network" json:"network"` // tcp or unix
	Username string `koanf:"username" json:"username"`
	Password string `koanf:"password" json:"password"`
	DB       int    `koanf:"db" json:"db"` // Redis 默认数据库 0
}

type IPRegionConf struct {
	Enabled   bool   `koanf:"enabled" json:"enabled"`     // 启用 IP 归属地展示
	DBPath    string `koanf:"db_path" json:"db_path"`     // 数据文件路径
	Precision string `koanf:"precision" json:"precision"` // 显示精度
}

type IPRegionPrecision string

const (
	IPRegionProvince IPRegionPrecision = "province"
	IPRegionCity     IPRegionPrecision = "city"
	IPRegionCountry  IPRegionPrecision = "country"
)

type ImgUploadConf struct {
	Enabled    bool      `koanf:"enabled" json:"enabled"`         // 总开关
	Path       string    `koanf:"path" json:"path"`               // 图片存放路径
	MaxSize    int64     `koanf:"max_size" json:"max_size"`       // 图片大小限制
	Quality    string    `koanf:"quality" json:"quality"`         // 图片质量
	PublicPath string    `koanf:"public_path" json:"public_path"` // 图片 URL 基础路径
	Upgit      UpgitConf `koanf:"upgit" json:"upgit"`             // upgit
}

type UpgitConf struct {
	Enabled  bool   `koanf:"enabled" json:"enabled"`     // 启用 upgit
	Exec     string `koanf:"exec" json:"exec"`           // 启动命令
	DelLocal bool   `koanf:"del_local" json:"del_local"` // 上传后删除本地的图片
}

// 其他通知方式
type AdminNotifyConf struct {
	NotifyTpl     string             `koanf:"notify_tpl" json:"notify_tpl"`         // 通知模板
	NotifySubject string             `koanf:"notify_subject" json:"notify_subject"` // 通知标题
	Email         *AdminEmailConf    `koanf:"email" json:"email"`                   // 邮件通知
	Telegram      NotifyTelegramConf `koanf:"telegram" json:"telegram"`             // TG
	Lark          NotifyLarkConf     `koanf:"lark" json:"lark"`                     // 飞书
	DingTalk      NotifyDingTalkConf `koanf:"ding_talk" json:"ding_talk"`           // 钉钉
	Bark          NotifyBarkConf     `koanf:"bark" json:"bark"`                     // bark
	Slack         NotifySlackConf    `koanf:"slack" json:"slack"`                   // slack
	LINE          NotifyLINEConf     `koanf:"line" json:"line"`                     // LINE
	WebHook       NotifyWebHookConf  `koanf:"webhook" json:"webhook"`               // WebHook
	NotifyPending bool               `koanf:"notify_pending" json:"notify_pending"` // 待审核评论通知
	NoiseMode     bool               `koanf:"noise_mode" json:"noise_mode"`         // 嘈杂模式 (非回复管理员的评论也发送通知)
}

type AdminEmailConf struct {
	Enabled     bool   `koanf:"enabled" json:"enabled"`           // 管理员总开关
	MailSubject string `koanf:"mail_subject" json:"mail_subject"` // 管理员邮件标题
	MailTpl     string `koanf:"mail_tpl" json:"mail_tpl"`         // 管理员专用邮件模板
}

type NotifyTelegramConf struct {
	Enabled   bool    `koanf:"enabled" json:"enabled"`
	ApiToken  string  `koanf:"api_token" json:"api_token"`
	Receivers []int64 `koanf:"receivers" json:"receivers"`
}

type NotifyDingTalkConf struct {
	Enabled bool   `koanf:"enabled" json:"enabled"`
	Token   string `koanf:"token" json:"token"`
	Secret  string `koanf:"secret" json:"secret"`
}

type NotifyLarkConf struct {
	Enabled    bool   `koanf:"enabled" json:"enabled"`
	WebhookURL string `koanf:"webhook_url" json:"webhook_url"`
}

type NotifyBarkConf struct {
	Enabled bool   `koanf:"enabled" json:"enabled"`
	Server  string `koanf:"server" json:"server"`
}

type NotifySlackConf struct {
	Enabled    bool     `koanf:"enabled" json:"enabled"`
	OauthToken string   `koanf:"oauth_token" json:"oauth_token"`
	Receivers  []string `koanf:"receivers" json:"receivers"`
}

type NotifyLINEConf struct {
	Enabled            bool     `koanf:"enabled" json:"enabled"`
	ChannelSecret      string   `koanf:"channel_secret" json:"channel_secret"`
	ChannelAccessToken string   `koanf:"channel_access_token" json:"channel_access_token"`
	Receivers          []string `koanf:"receivers" json:"receivers"`
}

type NotifyWebHookConf struct {
	Enabled bool   `koanf:"enabled" json:"enabled"`
	URL     string `koanf:"url" json:"url"`
}
