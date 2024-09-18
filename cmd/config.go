package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/artalkjs/artalk/v2/internal/log"
	"github.com/spf13/cobra"
)

func NewConfigCommand() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Output Config Information",
		Run: func(cmd *cobra.Command, args []string) {
			// Get config filename from cmd flags
			filename, _ := cmd.Flags().GetString("config")

			// Get new config instance
			config, err := getConfig(filename)
			if err != nil {
				log.Fatal("Config fail: ", err)
			}

			// Output JSON of config
			buf, _ := json.MarshalIndent(config, "", "    ")
			fmt.Println(string(buf))
		},
		Annotations: map[string]string{
			BootModeKey: MODE_MINI_BOOT,
		},
	}

	return configCmd
}
