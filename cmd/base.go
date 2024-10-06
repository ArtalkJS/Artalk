package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/adrg/xdg"
	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/log"
	"github.com/artalkjs/artalk/v2/internal/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var Banner = `
 ________  ________  _________  ________  ___       ___  __
|\   __  \|\   __  \|\___   ___\\   __  \|\  \     |\  \|\  \
\ \  \|\  \ \  \|\  \|___ \  \_\ \  \|\  \ \  \    \ \  \/  /|_
 \ \   __  \ \   _  _\   \ \  \ \ \   __  \ \  \    \ \   ___  \
  \ \  \ \  \ \  \\  \|   \ \  \ \ \  \ \  \ \  \____\ \  \\ \  \
   \ \__\ \__\ \__\\ _\    \ \__\ \ \__\ \__\ \_______\ \__\\ \__\
    \|__|\|__|\|__|\|__|    \|__|  \|__|\|__|\|_______|\|__| \|__|

Artalk (` + config.VersionString() + `)

 -> A Self-hosted Comment System.
 -> https://artalk.js.org
`

type ArtalkCmd struct {
	*core.App

	cfgFile string
	workDir string

	RootCmd *cobra.Command
}

const BootModeKey = "BootMode"

const (
	MODE_FULL_BOOT = "FULL_BOOT"
	MODE_MINI_BOOT = "MINI_BOOT"
)

func New() *ArtalkCmd {
	cmd := &ArtalkCmd{
		RootCmd: &cobra.Command{
			Use:     "artalk",
			Short:   "Artalk: A self-hosted comment system",
			Long:    Banner,
			Version: config.VersionString(),
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println(Banner)
				fmt.Print("-------------------------------\n\n")
				fmt.Println("NOTE: add `-h` flag to show help about any command.")
			},
		},
	}

	cmd.RootCmd.SetVersionTemplate("Artalk ({{printf \"%s\" .Version}})\n")

	// Parse base flags
	cmd.eagerParseFlags()

	// Init data directory (work directory)
	if workDir, err := initDataDir(cmd.workDir); err == nil {
		cmd.workDir = workDir
	} else {
		log.Fatal("Data directory fail: ", err)
	}

	return cmd
}

// Parses the global app flags before calling atk.RootCmd.Execute().
// so we can have all flags ready for use on initialization.
func (atk *ArtalkCmd) eagerParseFlags() error {
	atk.RootCmd.PersistentFlags().StringVarP(&atk.cfgFile, "config", "c", "", "config file path (defaults are './artalk.yml')")
	atk.RootCmd.PersistentFlags().StringVarP(&atk.workDir, "workdir", "w", "", "program working directory (defaults are './')")

	return atk.RootCmd.ParseFlags(os.Args[1:])
}

func (atk *ArtalkCmd) addCommand(cmd *cobra.Command) {
	originalPreRunFunc := cmd.PreRun

	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		// ================
		//  APP Bootstrap
		// ================
		if cmd.Annotations[BootModeKey] != string(MODE_MINI_BOOT) {
			// Load config
			config, err := getConfig(atk.cfgFile)
			if err != nil {
				log.Fatal("Config fail: ", err)
			}

			// Create new instance
			atk.App = core.NewApp(config)

			// Bootstrap APP
			if err := atk.App.Bootstrap(); err != nil {
				panic(err)
			}
		}

		if originalPreRunFunc != nil {
			originalPreRunFunc(cmd, args) // extends original pre-run logic
		}
	}

	atk.RootCmd.AddCommand(cmd)
}

func (atk *ArtalkCmd) mountCommands() {
	atk.addCommand(NewServeCommand(atk))
	atk.addCommand(NewAdminCommand(atk))
	atk.addCommand(NewExportCommand(atk))
	atk.addCommand(NewImportCommand(atk))
	atk.addCommand(NewConfigCommand())
	atk.addCommand(NewGenCommand())
	atk.addCommand(NewUpgradeCommand())
	atk.addCommand(NewVersionCommand())
}

func (atk *ArtalkCmd) Launch() error {
	// ===================
	//  1. Prepare Commands
	// ===================
	atk.mountCommands()

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
	if atk.App != nil {
		return atk.App.OnTerminate().Trigger(&core.TerminateEvent{
			App: atk.App,
		})
	}

	return nil
}

func initDataDir(workDir string) (string, error) {
	// If work directory is not specified, try to find it
	if workDir == "" {
		// If `data` in current directory exists, not change work directory
		if utils.CheckDirExist("./data") {
			// log.Warn("[DEPRECATED] The `data` directory is deprecated, please move to `~/.local/share/artalk/data`")
			return "", nil
		}

		// Retrieve data directory assuming it's already exists
		workDir = config.RetrieveDataDir()
	}

	// Create new data directory if not exists
	if workDir == "" {
		// default data path is `~/.local/share/artalk`
		workDir = path.Join(xdg.DataHome, "artalk")
		if err := utils.EnsureDir(workDir); err != nil {
			return "", fmt.Errorf("data directory creation error: %v", err)
		}
	}

	// Change work directory
	if workDir != "" {
		if err := os.Chdir(workDir); err != nil {
			return "", fmt.Errorf("working directory change error: %v", err)
		}
	}

	return workDir, nil
}

// Create new config instance by specific config filename
func getConfig(cfgFile string) (*config.Config, error) {
	// If config file is not specified, try to find it
	if cfgFile == "" {
		// Retrieve config file assuming it's already exists
		cfgFile = config.RetrieveConfigFile()
	}

	// Generate new config file if not exists
	if cfgFile == "" {
		// default config path is `~/.config/artalk/artalk.yml`
		cfgFile = path.Join(xdg.ConfigHome, "artalk", config.CONF_DEFAULT_FILENAMES[0])
		core.Gen("config", cfgFile, false)
	}

	// Create new config instance and return
	return config.NewFromFile(cfgFile)
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
