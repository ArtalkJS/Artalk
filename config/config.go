package config

import (
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
	AllowOrigins   []string        `mapstructure:"allow_origins" json:"allow_origins"`     // 允许跨域访问的域名
	TrustedDomains []string        `mapstructure:"trusted_domains" json:"trusted_domains"` // 可信任的域名 (新)
	SSL            SSLConf         `mapstructure:"ssl" json:"ssl"`                         // SSL
	SiteDefault    string          `mapstructure:"site_default" json:"site_default"`       // 默认站点名（当请求无指定 site_name 时使用）
	AdminUsers     []AdminUserConf `mapstructure:"admin_users" json:"admin_users"`         // 管理员账户
	LoginTimeout   int             `mapstructure:"login_timeout" json:"login_timeout"`     // 登陆超时
	Moderator      ModeratorConf   `mapstructure:"moderator" json:"moderator"`             // 评论审查
	Captcha        CaptchaConf     `mapstructure:"captcha" json:"captcha"`                 // 验证码
	Email          EmailConf       `mapstructure:"email" json:"email"`                     // 邮箱提醒
	ImgUpload      ImgUploadConf   `mapstructure:"img_upload" json:"img_upload"`           // 图片上传
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
	PendingDefault bool   `mapstructure:"pending_default" json:"pending_default"`
	AkismetKey     string `mapstructure:"akismet_key" json:"akismet_key"`
}

type CaptchaConf struct {
	Enabled       bool `mapstructure:"enabled" json:"enabled"`
	Always        bool `mapstructure:"always" json:"always"`
	ActionTimeout int  `mapstructure:"action_timeout" json:"action_timeout"`
	ActionLimit   int  `mapstructure:"action_limit" json:"action_limit"`
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
	Enabled bool      `mapstructure:"enabled" json:"enabled"`   // 总开关
	Path    string    `mapstructure:"path" json:"path"`         // 图片存放路径
	MaxSize string    `mapstructure:"max_size" json:"max_size"` // 图片大小限制
	Quality string    `mapstructure:"quality" json:"quality"`   // 图片质量
	Upgit   UpgitConf `mapstructure:"upgit" json:"upgit"`       // upgit
}

type UpgitConf struct {
	Enabled bool   `mapstructure:"enabled" json:"enabled"` // 启用 upgit
	Exec    string `mapstructure:"exec" json:"exec"`       // 启动命令
}

// Init 初始化配置
func Init(cfgFile string) {
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

	if err := viper.ReadInConfig(); err == nil {
		// fmt.Print("\n")
		// fmt.Println("- Using ArtalkGo config file:", viper.ConfigFileUsed())
	}

	Instance = &Config{}
	err := viper.Unmarshal(&Instance)
	if err != nil {
		logrus.Errorf("unable to decode into struct, %v", err)
	}
}
