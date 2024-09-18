package cmd

import (
	"fmt"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/spf13/cobra"
)

func NewVersionCommand() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Output version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Artalk (" + config.VersionString() + ")")
		},
		Annotations: map[string]string{
			BootModeKey: MODE_MINI_BOOT,
		},
	}

	return versionCmd
}
