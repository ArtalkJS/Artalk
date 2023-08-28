package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var Version = config.Version + `/` + config.CommitHash

var Banner = `
 ________  ________  _________  ________  ___       ___  __
|\   __  \|\   __  \|\___   ___\\   __  \|\  \     |\  \|\  \
\ \  \|\  \ \  \|\  \|___ \  \_\ \  \|\  \ \  \    \ \  \/  /|_
 \ \   __  \ \   _  _\   \ \  \ \ \   __  \ \  \    \ \   ___  \
  \ \  \ \  \ \  \\  \|   \ \  \ \ \  \ \  \ \  \____\ \  \\ \  \
   \ \__\ \__\ \__\\ _\    \ \__\ \ \__\ \__\ \_______\ \__\\ \__\
    \|__|\|__|\|__|\|__|    \|__|  \|__|\|__|\|_______|\|__| \|__|

Artalk (` + Version + `)

 -> A Selfhosted Comment System.
 -> https://artalk.js.org
`

type ArtalkCmd struct {
	*core.App

	cfgFile string
	workDir string

	RootCmd *cobra.Command
}

func New() *ArtalkCmd {
	cmd := &ArtalkCmd{
		RootCmd: &cobra.Command{
			Use:     "artalk",
			Short:   "Artalk: Your self-hosted comment system",
			Long:    Banner,
			Version: Version,
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println(Banner)
				fmt.Print("-------------------------------\n\n")
				fmt.Println("NOTE: add `-h` flag to show help about any command.")
			},
		},
	}

	cmd.RootCmd.SetVersionTemplate("Artalk ({{printf \"%s\" .Version}})\n")
	cmd.RootCmd.PersistentFlags().StringVarP(&cmd.cfgFile, "config", "c", "", "config file path (defaults are './artalk.yml')")
	cmd.RootCmd.PersistentFlags().StringVarP(&cmd.workDir, "workdir", "w", "", "program working directory (defaults are './')")

	// 切换工作目录
	if cmd.workDir != "" {
		if err := os.Chdir(cmd.workDir); err != nil {
			log.Fatal("Working directory change error: ", err)
		}
	}

	// initialize app instance
	config := getConfig(cmd.cfgFile)
	cmd.App = core.NewApp(config)

	return cmd
}

func (atk *ArtalkCmd) mountCommands() {
	atk.RootCmd.AddCommand(NewServeCommand(atk.App))
	atk.RootCmd.AddCommand(NewAdminCommand(atk.App))
	atk.RootCmd.AddCommand(NewExportCommand(atk.App))
	atk.RootCmd.AddCommand(NewImportCommand(atk.App))

	// TODO running without bootstrapped env
	atk.RootCmd.AddCommand(NewConfigCommand(atk.App))
	atk.RootCmd.AddCommand(NewGenCommand(atk.App))
	atk.RootCmd.AddCommand(NewUpgradeCommand(atk.App))
	atk.RootCmd.AddCommand(NewVersionCommand(atk.App))
}

func (atk *ArtalkCmd) Launch() error {
	// mount system commands
	atk.mountCommands()

	// ===================
	//  1. App Bootstrap
	// ===================
	if err := atk.Bootstrap(); err != nil {
		return err
	}

	done := make(chan bool, 1) // shutdown signal

	// listen for interrupt signal to gracefully shutdown the application
	go func() {
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
		<-sigch

		done <- true
	}()

	// ===================
	//  2. Command Execute
	// ===================
	go func() {
		if err := atk.RootCmd.Execute(); err != nil {
			color.Red(err.Error())
		}

		done <- true
	}()

	<-done

	// ===================
	//  3. App Cleanups
	// ===================
	return atk.OnTerminate().Trigger(&core.TerminateEvent{
		App: atk.App,
	})
}

// 获取配置
func getConfig(cfgFile string) *config.Config {
	// 尝试查找配置文件
	if cfgFile == "" {
		cfgFile = config.RetrieveConfigFile()
	}

	// 自动生成新配置文件
	if cfgFile == "" {
		cfgFile = config.CONF_DEFAULT_FILENAMES[0]
		core.Gen("config", cfgFile, false)
	}

	conf := config.NewFromFile(cfgFile)
	return conf
}

// -------------------------------------------------------------------
//  Shortcut Functions
// -------------------------------------------------------------------

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
}

func flagV(cmd *cobra.Command, name string, defaultVal interface{}, usage string) {
	flag(cmd, name, defaultVal, usage)
}

func flagPV(cmd *cobra.Command, name, shorthand string, defaultVal interface{}, usage string) {
	flagP(cmd, name, shorthand, defaultVal, usage)
}
