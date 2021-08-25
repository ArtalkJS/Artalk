package config

import (
	"fmt"
	"strings"

	"github.com/ArtalkJS/Artalk-API-Go/model"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Instance 配置实例
var Instance *Config

// Config 配置
// @link https://godoc.org/github.com/mitchellh/mapstructure
type Config struct {
	SiteName     string        `mapstructure:"site_name"`    // 网站名
	AppKey       string        `mapstructure:"app_key"`      // 加密密钥
	HttpAddr     string        `mapstructure:"http_addr"`    // HTTP Server 监听地址
	DB           DBConf        `mapstructure:"db"`           // 数据文件
	Log          LogConf       `mapstructure:"log"`          // 日志文件
	AllowOrigin  []string      `mapstructure:"allow_origin"` // 允许跨域访问的域名
	AdminUsers   []model.User  `mapstructure:"admin_users"`
	LoginTimeout int           `mapstructure:"login_timeout"`
	Moderator    ModeratorConf `mapstructure:"moderator"`
	Captcha      CaptchaConf   `mapstructure:"captcha"`
	Email        EmailConf     `mapstructure:"email"`
}

type DBConf struct {
	Type DBType `mapstructure:"type"`
	Dsn  string `mapstructure:"dsn"`
}

type LogConf struct {
	Enabled  bool   `mapstructure:"enabled"`
	Filename string `mapstructure:"filename"`
}

type ModeratorConf struct {
	PendingDefault bool `mapstructure:"pending_default"`
}

type CaptchaConf struct {
	Enabled      bool `mapstructure:"enabled"`
	Timeout      uint `mapstructure:"timeout"`
	CommentLimit uint `mapstructure:"comment_limit"`
}

type EmailConf struct {
	Enabled          bool            `mapstructure:"enabled"`
	AdminAddr        string          `mapstructure:"admin_addr"`
	SenderType       EmailSenderType `mapstructure:"sender_type"`
	MailTitle        string          `mapstructure:"mail_title"`
	MailTitleToAdmin string          `mapstructure:"mail_title_to_admin"`
	MailTplName      string          `mapstructure:"mail_tpl_name"`
	SMTP             SMTPConf        `mapstructure:"smtp"`
	AliDM            AliDMConf       `mapstructure:"ali_dm"`
}

type SMTPConf struct {
	Host       string `mapstructure:"host"`
	Port       uint   `mapstructure:"port"`
	SMTPAuth   bool   `mapstructure:"smtp_auth"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	SMTPSecure string `mapstructure:"smtp_secure"`
	FromAddr   string `mapstructure:"from_addr"`
	FromName   string `mapstructure:"from_name"`
}

type AliDMConf struct {
	AccessKeyId     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	AccountName     string `mapstructure:"account_name"`
}

type DBType string

const (
	TypeMySql  DBType = "mysql"
	TypeSqlite DBType = "sqlite"
)

type EmailSenderType string

const (
	SMTP  EmailSenderType = "smpt"
	AliDM EmailSenderType = "ali_dm"
)

// Init 初始化配置
func Init(cfgFile string) {
	viper.SetConfigType("yaml")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find config file in path.
		viper.AddConfigPath(".")
		viper.SetConfigName("artalk-go.yaml")
	}

	viper.SetEnvPrefix("ATG")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	Instance = &Config{}
	err := viper.Unmarshal(&Instance)
	if err != nil {
		logrus.Errorf("unable to decode into struct, %v", err)
	}
}
