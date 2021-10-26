package cmd

import (
	"fmt"

	"github.com/ArtalkJS/ArtalkGo/http"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"server"},
	Short:   "HTTP 服务",
	Long:    Banner,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Banner)
		fmt.Print("-------------------------------\n\n")
		http.Run()
	},
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(serveCmd)

	flagPV(serveCmd, "host", "", "0.0.0.0", "监听 IP")
	flagPV(serveCmd, "port", "", 23366, "监听端口")
}
