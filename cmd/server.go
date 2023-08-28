package cmd

import (
	"fmt"

	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/ArtalkJS/Artalk/server"

	"github.com/spf13/cobra"
)

func NewServeCommand(app *core.App) *cobra.Command {
	var serverCmd = &cobra.Command{
		Use:     "server",
		Aliases: []string{"serve"},
		Short:   "Start the server",
		Long:    Banner,
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Banner)
			fmt.Print("-------------------------------\n\n")

			// init fiber app
			_, err := server.Serve(app)
			if err != nil {
				log.Fatalln(err)
			}
		},
	}

	flagPV(serverCmd, "host", "", "0.0.0.0", "Listening IP")
	flagPV(serverCmd, "port", "", 23366, "Listening port")

	return serverCmd
}
