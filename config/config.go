package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jeremywohl/flatten"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Instance 配置实例
var Instance *Config

var Flat map[string]interface{}

// Config 配置
// @link https://godoc.org/github.com/mitchellh/mapstructure
type Config struct {
	SiteName     string          `mapstructure:"site_name"`    // 网站名
	AppKey       string          `mapstructure:"app_key"`      // 加密密钥
	Debug        bool            `mapstructure:"debug"`        // 调试模式
	HttpAddr     string          `mapstructure:"http_addr"`    // HTTP Server 监听地址
	DB           DBConf          `mapstructure:"db"`           // 数据文件
	Log          LogConf         `mapstructure:"log"`          // 日志文件
	AllowOrigin  []string        `mapstructure:"allow_origin"` // 允许跨域访问的域名
	AdminUsers   []AdminUserConf `mapstructure:"admin_users"`
	LoginTimeout int             `mapstructure:"login_timeout"`
	Moderator    ModeratorConf   `mapstructure:"moderator"`
	Captcha      CaptchaConf     `mapstructure:"captcha"`
	Email        EmailConf       `mapstructure:"email"`
}

type DBConf struct {
	Type DBType `mapstructure:"type"`
	Dsn  string `mapstructure:"dsn"`
}

type LogConf struct {
	Enabled  bool   `mapstructure:"enabled"`
	Filename string `mapstructure:"filename"`
}

type AdminUserConf struct {
	Name       string `mapstructure:"name"`
	Email      string `mapstructure:"email"`
	Link       string `mapstructure:"link"`
	Password   string `mapstructure:"password"`
	BadgeName  string `mapstructure:"badge_name"`
	BadgeColor string `mapstructure:"badge_color"`
	SiteID     uint   `mapstructure:"site_id"`
}

type ModeratorConf struct {
	PendingDefault bool `mapstructure:"pending_default"`
}

type CaptchaConf struct {
	Enabled       bool `mapstructure:"enabled"`
	Always        bool `mapstructure:"always"`
	ActionTimeout int  `mapstructure:"action_timeout"`
	ActionLimit   int  `mapstructure:"action_limit"`
}

type EmailConf struct {
	Enabled            bool            `mapstructure:"enabled"`               // 总开关
	SendType           EmailSenderType `mapstructure:"send_type"`             // 发送方式
	SendName           string          `mapstructure:"send_name"`             // 发件人名
	SendAddr           string          `mapstructure:"send_addr"`             // 发件人地址
	MailSubject        string          `mapstructure:"mail_subject"`          // 邮件标题
	MailSubjectToAdmin string          `mapstructure:"mail_subject_to_admin"` // 邮件标题 (发送给管理员用)
	MailTpl            string          `mapstructure:"mail_tpl"`              // 邮件模板
	SMTP               SMTPConf        `mapstructure:"smtp"`                  // SMTP 配置
	AliDM              AliDMConf       `mapstructure:"ali_dm"`                // 阿里云邮件配置
}

type SMTPConf struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

type AliDMConf struct {
	AccessKeyId     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	AccountName     string `mapstructure:"account_name"`
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
	TypeSMTP     EmailSenderType = "smpt"
	TypeAliDM    EmailSenderType = "ali_dm"
	TypeSendmail EmailSenderType = "sendmail"
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
		viper.SetConfigName("artalk-go.yml")
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

	Flat = StructToFlatDotMap(&Instance)

	if strings.TrimSpace(Instance.AppKey) == "" {
		logrus.Fatal("请检查配置文件，并设置一个 app_key")
	}
}

func StructToMap(s interface{}) map[string]interface{} {
	b, _ := json.Marshal(s)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	return m
}

func StructToFlatDotMap(s interface{}) map[string]interface{} {
	m := StructToMap(s)
	mainFlat, err := flatten.Flatten(m, "", flatten.DotStyle)
	if err != nil {
		return map[string]interface{}{}
	}
	return mainFlat
}
