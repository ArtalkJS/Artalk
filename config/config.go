package config

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Instance 配置实例
var Instance *Config

// Config 配置
// @link https://godoc.org/github.com/mitchellh/mapstructure
type Config struct {
	SiteName string `mapstructure:"site_name"` // 网站名
	HttpAddr string `mapstructure:"http_addr"` // HTTP Server 监听地址

	DbFile  string `mapstructure:"db_file"`  // 数据文件
	LogFile string `mapstructure:"log_file"` // 日志文件
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
