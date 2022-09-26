package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/lib/core"
	"github.com/spf13/cobra"
)

var Version = lib.Version + `/` + lib.CommitHash

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
	Use:     "artalk-go",
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
	rootCmd.SetVersionTemplate("ArtalkGo ({{printf \"%s\" .Version}})\n")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "配置文件路径 (defaults are './artalk-go.yml')")
	rootCmd.PersistentFlags().StringVarP(&workDir, "workdir", "w", "", "程序工作目录 (defaults are './')")

	// Version Command
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "输出版本信息",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("ArtalkGo (" + Version + ")")
		},
	}
	rootCmd.AddCommand(versionCmd)

	// Config Command
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "输出配置信息",
		Run: func(cmd *cobra.Command, args []string) {
			core.LoadConfOnly(cfgFile, workDir)
			buf, _ := json.MarshalIndent(config.Instance, "", "    ")
			fmt.Println(string(buf))
		},
	}
	rootCmd.AddCommand(configCmd)

}
