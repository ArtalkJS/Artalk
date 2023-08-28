package cmd

import (
	"fmt"

	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/spf13/cobra"
)

func NewVersionCommand(app *core.App) *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Output version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Artalk (" + Version + ")")
		},
	}

	return versionCmd
}
