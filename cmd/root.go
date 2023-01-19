package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/core"
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

var (
	cfgFile string
	workDir string
)

var rootCmd = &cobra.Command{
	Use:     "artalk",
	Short:   "Artalk: A Fast, Slight & Delightful Comment System",
	Long:    Banner,
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Banner)
		fmt.Print("-------------------------------\n\n")
		fmt.Println("NOTE: add `-h` flag to show help about any command.")
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.SetVersionTemplate("Artalk ({{printf \"%s\" .Version}})\n")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file path (defaults are './artalk.yml')")
	rootCmd.PersistentFlags().StringVarP(&workDir, "workdir", "w", "", "program working directory (defaults are './')")

	// Version Command
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Output Version Information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Artalk (" + Version + ")")
		},
	}
	rootCmd.AddCommand(versionCmd)

	// Config Command
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Output Config Information",
		Run: func(cmd *cobra.Command, args []string) {
			core.LoadConfOnly(cfgFile, workDir)
			buf, _ := json.MarshalIndent(config.Instance, "", "    ")
			fmt.Println(string(buf))
		},
	}
	rootCmd.AddCommand(configCmd)

}
