package cmd

import (
	"fmt"
	"os"

	"github.com/ArtalkJS/Artalk-API-Go/config"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var Banner = `
 ________  ________  _________  ________  ___       ___  __       
|\   __  \|\   __  \|\___   ___\\   __  \|\  \     |\  \|\  \     
\ \  \|\  \ \  \|\  \|___ \  \_\ \  \|\  \ \  \    \ \  \/  /|_   
 \ \   __  \ \   _  _\   \ \  \ \ \   __  \ \  \    \ \   ___  \  
  \ \  \ \  \ \  \\  \|   \ \  \ \ \  \ \  \ \  \____\ \  \\ \  \ 
   \ \__\ \__\ \__\\ _\    \ \__\ \ \__\ \__\ \_______\ \__\\ \__\
    \|__|\|__|\|__|\|__|    \|__|  \|__|\|__|\|_______|\|__| \|__|
 
Artalk: A Fast, Slight & Funny Comment System.
More detail on https://github.com/ArtalkJS/Artalk
(c) 2021 artalk.js.org`

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "artalk-go",
	Short: "Artalk: A Fast, Slight & Funny Comment System",
	Long:  Banner,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Banner)
		fmt.Println()
		fmt.Println("NOTE: add `-h` flag to show help about any command.")
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initLog)
	cobra.OnInitialize(initDB)

	rootCmd.SetVersionTemplate("Artalk-GO {{printf \"version %s\" .Version}}\n")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件路径 (defaults are './artalk-go.conf.yaml', './config.yaml')")
}

func initConfig() {
	config.Init(cfgFile)
}

func initLog() {
	// 初始化日志
	stdFormatter := &prefixed.TextFormatter{
		DisableTimestamp: true,
		ForceFormatting:  true,
		ForceColors:      true,
		DisableColors:    false,
	} // 命令行输出格式
	fileFormatter := &prefixed.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02.15:04:05.000000",
		ForceFormatting: true,
		ForceColors:     false,
		DisableColors:   true,
	} // 文件输出格式

	// logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(stdFormatter)
	logrus.SetOutput(os.Stdout)

	// 文件保存
	pathMap := lfshook.PathMap{
		logrus.InfoLevel:  config.Instance.LogFile,
		logrus.DebugLevel: config.Instance.LogFile,
		logrus.ErrorLevel: config.Instance.LogFile,
	}
	logrus.AddHook(lfshook.NewHook(
		pathMap,
		fileFormatter,
	))
}

func initDB() {
}

//// 捷径函数 ////

func flag(cmd *cobra.Command, name string, defaultVal interface{}, usage string) {
	f := cmd.PersistentFlags()
	switch y := defaultVal.(type) {
	case bool:
		f.Bool(name, y, usage)
	case int:
		f.Int(name, y, usage)
	case string:
		f.String(name, y, usage)
	}
	viper.SetDefault(name, defaultVal)
}

func flagP(cmd *cobra.Command, name, shorthand string, defaultVal interface{}, usage string) {
	f := cmd.PersistentFlags()
	switch y := defaultVal.(type) {
	case bool:
		f.BoolP(name, shorthand, y, usage)
	case int:
		f.IntP(name, shorthand, y, usage)
	case string:
		f.StringP(name, shorthand, y, usage)
	}
	viper.SetDefault(name, defaultVal)
}

func flagV(cmd *cobra.Command, name string, defaultVal interface{}, usage string) {
	flag(cmd, name, defaultVal, usage)
	viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
}

func flagPV(cmd *cobra.Command, name, shorthand string, defaultVal interface{}, usage string) {
	flagP(cmd, name, shorthand, defaultVal, usage)
	viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
}
