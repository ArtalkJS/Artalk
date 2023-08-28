package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/spf13/cobra"
)

func NewConfigCommand(app *core.App) *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Output Config Information",
		Run: func(cmd *cobra.Command, args []string) {
			buf, _ := json.MarshalIndent(app.Conf(), "", "    ")
			fmt.Println(string(buf))
		},
	}

	return configCmd
}
