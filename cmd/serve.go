package cmd

import (
	"fmt"

	"github.com/ArtalkJS/Artalk-API-Go/http"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:     "serve",
	Version: rootCmd.Version,
	Aliases: []string{"server"},
	Short:   "启动 HTTP 服务",
	Long:    rootCmd.Long,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Banner)
		fmt.Println("----------------------")
		http.Run()
	},
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
