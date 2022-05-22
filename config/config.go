package config

import "time"

// Config 配置
// @link https://godoc.org/github.com/mitchellh/mapstructure
type Config struct {
	AppKey         string           `mapstructure:"app_key" json:"app_key"`                 // 加密密钥
	Debug          bool             `mapstructure:"debug" json:"debug"`                     // 调试模式
	TimeZone       string           `mapstructure:"timezone" json:"timezone"`               // 时区
	Host           string           `mapstructure:"host" json:"host"`                       // HTTP Server 监听 IP
	Port           int              `mapstructure:"port" json:"port"`                       // HTTP Server 监听 Port
	DB             DBConf           `mapstructure:"db" json:"db"`                           // 数据文件
	Cache          CacheConf        `mapstructure:"cache" json:"cache"`                     // 缓存
	Log            LogConf          `mapstructure:"log" json:"log"`                         // 日志文件
	AllowOrigins   []string         `mapstructure:"allow_origins" json:"-"`                 // @deprecated 已废弃 (请使用 TrustedDomains)
	TrustedDomains []string         `mapstructure:"trusted_domains" json:"trusted_domains"` // 可信任的域名 (新)
	SSL            SSLConf          `mapstructure:"ssl" json:"ssl"`                         // SSL
	SiteDefault    string           `mapstructure:"site_default" json:"site_default"`       // 默认站点名（当请求无指定 site_name 时使用）
	AdminUsers     []AdminUserConf  `mapstructure:"admin_users" json:"admin_users"`         // 管理员账户
	LoginTimeout   int              `mapstructure:"login_timeout" json:"login_timeout"`     // 登陆超时
	Moderator      ModeratorConf    `mapstructure:"moderator" json:"moderator"`             // 评论审查
	Captcha        CaptchaConf      `mapstructure:"captcha" json:"captcha"`                 // 验证码
	Email          EmailConf        `mapstructure:"email" json:"email"`                     // 邮箱提醒
	ImgUpload      ImgUploadConf    `mapstructure:"img_upload" json:"img_upload"`           // 图片上传
	AdminNotify    AdminNotifyConf  `mapstructure:"admin_notify" json:"admin_notify"`       // 其他通知方式
	Notify         *AdminNotifyConf `mapstructure:"notify" json:"-"`                        // @deprecated 已废弃 (请使用 AdminNotify)
	Frontend       FrontendConf     `mapstructure:"frontend" json:"frontend"`
}

type DBConf struct {
	Type DBType `mapstructure:"type" json:"type"`
	Dsn  string `mapstructure:"dsn" json:"dsn"` // 最高优先级

	File string `mapstructure:"file" json:"file"`
	Name string `mapstructure:"name" json:"name"`

	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`

	TablePrefix string `mapstructure:"table_prefix" json:"table_prefix"`
	Charset     string `mapstructure:"charset" json:"charset"`
}

type CacheConf struct {
	Enabled bool      // 配置文件不允许修改
	Type    CacheType `mapstructure:"type" json:"type"`
	Expires int       `mapstructure:"expires" json:"expires"` // 过期时间
	WarmUp  bool      `mapstructure:"warm_up" json:"warm_up"` // 启动时缓存预热
	Server  string    `mapstructure:"server" json:"server"`   // 缓存服务器
	Redis   RedisConf `mapstructure:"redis" json:"redis"`
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
	Enabled  bool   `mapstructure:"enabled" json:"enabled"`
	Filename string `mapstructure:"filename" json:"filename"`
}

type SSLConf struct {
	Enabled  bool   `mapstructure:"enabled" json:"enabled"`
	CertPath string `mapstructure:"cert_path" json:"cert_path"`
	KeyPath  string `mapstructure:"key_path" json:"key_path"`
}

type AdminUserConf struct {
	Name         string   `mapstructure:"name" json:"name"`
	Email        string   `mapstructure:"email" json:"email"`
	Link         string   `mapstructure:"link" json:"link"`
	Password     string   `mapstructure:"password" json:"password"`
	BadgeName    string   `mapstructure:"badge_name" json:"badge_name"`
	BadgeColor   string   `mapstructure:"badge_color" json:"badge_color"`
	ReceiveEmail *bool    `mapstructure:"receive_email" json:"receive_email"`
	Sites        []string `mapstructure:"sites" json:"sites"`
}

type ModeratorConf struct {
	PendingDefault bool                 `mapstructure:"pending_default" json:"pending_default"`
	ApiFailBlock   bool                 `mapstructure:"api_fail_block" json:"api_fail_block"` // API 请求错误仍然拦截
	AkismetKey     string               `mapstructure:"akismet_key" json:"akismet_key"`
	Tencent        TencentAntispamConf  `mapstructure:"tencent" json:"tencent"`
	Aliyun         AliyunAntispamConf   `mapstructure:"aliyun" json:"aliyun"`
	Keywords       KeyWordsAntispamConf `mapstructure:"keywords" json:"keywords"`
}

// 腾讯云反垃圾
type TencentAntispamConf struct {
	Enabled   bool   `mapstructure:"enabled" json:"enabled"`
	SecretID  string `mapstructure:"secret_id" json:"secret_id"`
	SecretKey string `mapstructure:"secret_key" json:"secret_key"`
	Region    string `mapstructure:"region" json:"region"`
}

// 阿里云反垃圾
type AliyunAntispamConf struct {
	Enabled         bool   `mapstructure:"enabled" json:"enabled"`
	AccessKeyID     string `mapstructure:"access_key_id" json:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret" json:"access_key_secret"`
	Region          string `mapstructure:"region" json:"region"`
}

// 关键词词库过滤
type KeyWordsAntispamConf struct {
	Enabled  bool     `mapstructure:"enabled" json:"enabled"`
	Pending  bool     `mapstructure:"pending" json:"pending"`
	Files    []string `mapstructure:"files" json:"files"`
	FileSep  string   `mapstructure:"file_sep" json:"file_sep"`
	ReplacTo string   `mapstructure:"replac_to" json:"replac_to"`
}

type CaptchaConf struct {
	Enabled       bool        `mapstructure:"enabled" json:"enabled"`
	Always        bool        `mapstructure:"always" json:"always"`
	ActionTimeout int         `mapstructure:"action_timeout" json:"-"` // @deprecated 已废弃 (请使用 ActionReset)
	ActionReset   int         `mapstructure:"action_reset" json:"action_reset"`
	ActionLimit   int         `mapstructure:"action_limit" json:"action_limit"`
	Geetest       GeetestConf `mapstructure:"geetest" json:"geetest"`
}

type GeetestConf struct {
	Enabled    bool   `mapstructure:"enabled" json:"enabled"`
	CaptchaID  string `mapstructure:"captcha_id" json:"captcha_id"`
	CaptchaKey string `mapstructure:"captcha_key" json:"captcha_key"`
}

type EmailConf struct {
	Enabled            bool            `mapstructure:"enabled" json:"enabled"`           // 总开关
	SendType           EmailSenderType `mapstructure:"send_type" json:"send_type"`       // 发送方式
	SendName           string          `mapstructure:"send_name" json:"send_name"`       // 发件人名
	SendAddr           string          `mapstructure:"send_addr" json:"send_addr"`       // 发件人地址
	MailSubject        string          `mapstructure:"mail_subject" json:"mail_subject"` // 邮件标题
	MailSubjectToAdmin string          `mapstructure:"mail_subject_to_admin" json:"-"`   // @deprecated 已废弃 (请使用 AdminNotify.Email.MailSubject) - 邮件标题 (发送给管理员用)
	MailTpl            string          `mapstructure:"mail_tpl" json:"mail_tpl"`         // 邮件模板
	SMTP               SMTPConf        `mapstructure:"smtp" json:"smtp"`                 // SMTP 配置
	AliDM              AliDMConf       `mapstructure:"ali_dm" json:"ali_dm"`             // 阿里云邮件配置
}

type SMTPConf struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
	From     string `mapstructure:"from" json:"from"`
}

type AliDMConf struct {
	AccessKeyId     string `mapstructure:"access_key_id" json:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret" json:"access_key_secret"`
	AccountName     string `mapstructure:"account_name" json:"account_name"`
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
	CacheTypeDisabled CacheType = "disabled" // 关闭缓存
)

type EmailSenderType string

const (
	TypeSMTP     EmailSenderType = "smtp"
	TypeAliDM    EmailSenderType = "ali_dm"
	TypeSendmail EmailSenderType = "sendmail"
)

// # Redis 配置
// redis:
//   network: "tcp"
//   username: ""
//   password: ""
//   db: 0
type RedisConf struct {
	Network  string `mapstructure:"network" json:"network"` // tcp or unix
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
	DB       int    `mapstructure:"db" json:"db"` // Redis 默认数据库 0
}

type ImgUploadConf struct {
	Enabled    bool      `mapstructure:"enabled" json:"enabled"`         // 总开关
	Path       string    `mapstructure:"path" json:"path"`               // 图片存放路径
	MaxSize    int64     `mapstructure:"max_size" json:"max_size"`       // 图片大小限制
	Quality    string    `mapstructure:"quality" json:"quality"`         // 图片质量
	PublicPath string    `mapstructure:"public_path" json:"public_path"` // 图片 URL 基础路径
	Upgit      UpgitConf `mapstructure:"upgit" json:"upgit"`             // upgit
}

type UpgitConf struct {
	Enabled  bool   `mapstructure:"enabled" json:"enabled"`     // 启用 upgit
	Exec     string `mapstructure:"exec" json:"exec"`           // 启动命令
	DelLocal bool   `mapstructure:"del_local" json:"del_local"` // 上传后删除本地的图片
}

// 其他通知方式
type AdminNotifyConf struct {
	Email    *EmailConf         `mapstructure:"email" json:"email"`         // 邮件通知
	Telegram NotifyTelegramConf `mapstructure:"telegram" json:"telegram"`   // TG
	Lark     NotifyLarkConf     `mapstructure:"lark" json:"lark"`           // 飞书
	DingTalk NotifyDingTalkConf `mapstructure:"ding_talk" json:"ding_talk"` // 钉钉
	Bark     NotifyBarkConf     `mapstructure:"bark" json:"bark"`           // bark
	Slack    NotifySlackConf    `mapstructure:"slack" json:"slack"`         // slack
	LINE     NotifyLINEConf     `mapstructure:"line" json:"line"`           // LINE
}

type NotifyTelegramConf struct {
	Enabled   bool    `mapstructure:"enabled" json:"enabled"`
	ApiToken  string  `mapstructure:"api_token" json:"api_token"`
	Receivers []int64 `mapstructure:"receivers" json:"receivers"`
}

type NotifyDingTalkConf struct {
	Enabled bool   `mapstructure:"enabled" json:"enabled"`
	Token   string `mapstructure:"token" json:"token"`
	Secret  string `mapstructure:"secret" json:"secret"`
}

type NotifyLarkConf struct {
	Enabled    bool   `mapstructure:"enabled" json:"enabled"`
	WebhookURL string `mapstructure:"webhook_url" json:"webhook_url"`
}

type NotifyBarkConf struct {
	Enabled bool   `mapstructure:"enabled" json:"enabled"`
	Server  string `mapstructure:"server" json:"server"`
}

type NotifySlackConf struct {
	Enabled    bool     `mapstructure:"enabled" json:"enabled"`
	OauthToken string   `mapstructure:"oauth_token" json:"oauth_token"`
	Receivers  []string `mapstructure:"receivers" json:"receivers"`
}

type NotifyLINEConf struct {
	Enabled            bool     `mapstructure:"enabled" json:"enabled"`
	ChannelSecret      string   `mapstructure:"channel_secret" json:"channel_secret"`
	ChannelAccessToken string   `mapstructure:"channel_access_token" json:"channel_access_token"`
	Receivers          []string `mapstructure:"receivers" json:"receivers"`
}

// 使用转换 @link https://transform.tools/json-to-go
type FrontendConf struct {
	Placeholder  *string `mapstructure:"placeholder" json:"placeholder,omitempty"`
	NoComment    *string `mapstructure:"noComment" json:"noComment,omitempty"`
	SendBtn      *string `mapstructure:"sendBtn" json:"sendBtn,omitempty"`
	DarkMode     *bool   `mapstructure:"darkMode" json:"darkMode,omitempty"`
	EditorTravel *bool   `mapstructure:"editorTravel" json:"editorTravel,omitempty"`
	Emoticons    *string `mapstructure:"emoticons" json:"emoticons,omitempty"`
	Vote         *bool   `mapstructure:"vote" json:"vote,omitempty"`
	VoteDown     *bool   `mapstructure:"voteDown" json:"voteDown,omitempty"`
	UaBadge      *bool   `mapstructure:"uaBadge" json:"uaBadge,omitempty"`
	ListSort     *bool   `mapstructure:"listSort" json:"listSort,omitempty"`
	PvEl         *string `mapstructure:"pvEl" json:"pvEl,omitempty"`
	CountEl      *string `mapstructure:"countEl" json:"countEl,omitempty"`
	Preview      *bool   `mapstructure:"preview" json:"preview,omitempty"`
	FlatMode     *string `mapstructure:"flatMode" json:"flatMode,omitempty"`
	NestMax      *int    `mapstructure:"nestMax" json:"nestMax,omitempty"`
	NestSort     *string `mapstructure:"nestSort" json:"nestSort,omitempty"`
	Gravatar     *struct {
		Default *string `mapstructure:"default" json:"default,omitempty"`
		Mirror  *string `mapstructure:"mirror" json:"mirror,omitempty"`
	} `mapstructure:"gravatar" json:"gravatar,omitempty"`
	Pagination *struct {
		PageSize *int  `mapstructure:"pageSize" json:"pageSize,omitempty"`
		ReadMore *bool `mapstructure:"readMore" json:"readMore,omitempty"`
		AutoLoad *bool `mapstructure:"autoLoad" json:"autoLoad,omitempty"`
	} `mapstructure:"pagination" json:"pagination,omitempty"`
	HeightLimit *struct {
		Content  *int `mapstructure:"content" json:"content,omitempty"`
		Children *int `mapstructure:"children" json:"children,omitempty"`
	} `mapstructure:"heightLimit" json:"heightLimit,omitempty"`
	ImgUpload    *bool   `mapstructure:"imgUpload"  json:"imgUpload,omitempty"`
	ReqTimeout   *int    `mapstructure:"reqTimeout" json:"reqTimeout,omitempty"`
	VersionCheck *bool   `mapstructure:"versionCheck" json:"versionCheck,omitempty"`
	Locale       *string `mapstructure:"locale" json:"locale,omitempty"`
}
