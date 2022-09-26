package cmd

import (
	"fmt"

	"github.com/ArtalkJS/ArtalkGo/http"
	"github.com/ArtalkJS/ArtalkGo/lib/core"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"serve"},
	Short:   "启动服务器",
	Long:    Banner,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		core.LoadCore(cfgFile, workDir) // 装载核心

		fmt.Println(Banner)
		fmt.Print("-------------------------------\n\n")
		http.Run()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	flagPV(serverCmd, "host", "", "0.0.0.0", "监听 IP")
	flagPV(serverCmd, "port", "", 23366, "监听端口")
}
