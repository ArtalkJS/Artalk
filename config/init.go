package config

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Instance 配置实例
var Instance *Config

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
