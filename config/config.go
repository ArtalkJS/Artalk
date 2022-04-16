package config

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Instance 配置实例
var Instance *Config

// Config 配置
// @link https://godoc.org/github.com/mitchellh/mapstructure
type Config struct {
	AppKey         string          `mapstructure:"app_key" json:"app_key"`                 // 加密密钥
	Debug          bool            `mapstructure:"debug" json:"debug"`                     // 调试模式
	TimeZone       string          `mapstructure:"timezone" json:"timezone"`               // 时区
	Host           string          `mapstructure:"host" json:"host"`                       // HTTP Server 监听 IP
	Port           int             `mapstructure:"port" json:"port"`                       // HTTP Server 监听 Port
	DB             DBConf          `mapstructure:"db" json:"db"`                           // 数据文件
	Log            LogConf         `mapstructure:"log" json:"log"`                         // 日志文件
	AllowOrigins   []string        `mapstructure:"allow_origins" json:"allow_origins"`     // @deprecated 已废弃 (请使用 TrustedDomains)
	TrustedDomains []string        `mapstructure:"trusted_domains" json:"trusted_domains"` // 可信任的域名 (新)
	SSL            SSLConf         `mapstructure:"ssl" json:"ssl"`                         // SSL
	SiteDefault    string          `mapstructure:"site_default" json:"site_default"`       // 默认站点名（当请求无指定 site_name 时使用）
	AdminUsers     []AdminUserConf `mapstructure:"admin_users" json:"admin_users"`         // 管理员账户
	LoginTimeout   int             `mapstructure:"login_timeout" json:"login_timeout"`     // 登陆超时
	Moderator      ModeratorConf   `mapstructure:"moderator" json:"moderator"`             // 评论审查
	Captcha        CaptchaConf     `mapstructure:"captcha" json:"captcha"`                 // 验证码
	Email          EmailConf       `mapstructure:"email" json:"email"`                     // 邮箱提醒
	ImgUpload      ImgUploadConf   `mapstructure:"img_upload" json:"img_upload"`           // 图片上传
	Notify         NotifyConf      `mapstructure:"notify" json:"notify"`                   // 其他通知方式
}

type DBConf struct {
	Type DBType `mapstructure:"type" json:"type"`
	Dsn  string `mapstructure:"dsn" json:"dsn"`
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
	Name       string `mapstructure:"name" json:"name"`
	Email      string `mapstructure:"email" json:"email"`
	Link       string `mapstructure:"link" json:"link"`
	Password   string `mapstructure:"password" json:"password"`
	BadgeName  string `mapstructure:"badge_name" json:"badge_name"`
	BadgeColor string `mapstructure:"badge_color" json:"badge_color"`
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
	ActionTimeout int         `mapstructure:"action_timeout" json:"action_timeout"` // @deprecated 已废弃 (请使用 ActionReset)
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
	Enabled            bool            `mapstructure:"enabled" json:"enabled"`                             // 总开关
	SendType           EmailSenderType `mapstructure:"send_type" json:"send_type"`                         // 发送方式
	SendName           string          `mapstructure:"send_name" json:"send_name"`                         // 发件人名
	SendAddr           string          `mapstructure:"send_addr" json:"send_addr"`                         // 发件人地址
	MailSubject        string          `mapstructure:"mail_subject" json:"mail_subject"`                   // 邮件标题
	MailSubjectToAdmin string          `mapstructure:"mail_subject_to_admin" json:"mail_subject_to_admin"` // 邮件标题 (发送给管理员用)
	MailTpl            string          `mapstructure:"mail_tpl" json:"mail_tpl"`                           // 邮件模板
	SMTP               SMTPConf        `mapstructure:"smtp" json:"smtp"`                                   // SMTP 配置
	AliDM              AliDMConf       `mapstructure:"ali_dm" json:"ali_dm"`                               // 阿里云邮件配置
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
	TypeSqlServer  DBType = "sqlserver"
)

type EmailSenderType string

const (
	TypeSMTP     EmailSenderType = "smtp"
	TypeAliDM    EmailSenderType = "ali_dm"
	TypeSendmail EmailSenderType = "sendmail"
)

type ImgUploadConf struct {
	Enabled    bool      `mapstructure:"enabled" json:"enabled"`         // 总开关
	Path       string    `mapstructure:"path" json:"path"`               // 图片存放路径
	MaxSize    int64     `mapstructure:"max_size" json:"max_size"`       // 图片大小限制
	Quality    string    `mapstructure:"quality" json:"quality"`         // 图片质量
	PublicPath string    `mapstructure:"public_path" json:"public_path"` // 图片 URL 基础路径
	Upgit      UpgitConf `mapstructure:"upgit" json:"upgit"`             // upgit
}

type UpgitConf struct {
	Enabled bool   `mapstructure:"enabled" json:"enabled"` // 启用 upgit
	Exec    string `mapstructure:"exec" json:"exec"`       // 启动命令
}

// 其他通知方式
type NotifyConf struct {
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

// Init 初始化配置
func Init(cfgFile string, workDir string) {
	viper.SetConfigType("yaml")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find config file in path.
		viper.AddConfigPath(".")
		viper.SetConfigName("artalk-go.yml")
	}

	viper.SetEnvPrefix("ATG")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 切换工作目录
	if workDir != "" {
		viper.AddConfigPath(workDir) // must before
		if err := os.Chdir(workDir); err != nil {
			logrus.Fatal("工作目录切换错误 ", err)
		}
	}

	if err := viper.ReadInConfig(); err == nil {
		// fmt.Print("\n")
		// fmt.Println("- Using ArtalkGo config file:", viper.ConfigFileUsed())
	} else {
		logrus.Fatal("找不到配置文件，使用 `-h` 参数查看帮助")
	}

	Instance = &Config{}
	err := viper.Unmarshal(&Instance)
	if err != nil {
		logrus.Errorf("unable to decode into struct, %v", err)
	}
}
