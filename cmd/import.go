package cmd

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/artalkjs/artalk/v2/internal/artransfer"
	"github.com/artalkjs/artalk/v2/internal/i18n"
	"github.com/artalkjs/artalk/v2/internal/log"
	"github.com/spf13/cobra"
)

func NewImportCommand(app *ArtalkCmd) *cobra.Command {
	importCmd := &cobra.Command{
		Use:     "import <FILENAME>",
		Aliases: []string{},
		Short:   "Artransfer import",
		Long:    "\n# Artransfer - Import\n\n  See the documentation to learn more: https://artalk.js.org/guide/transfer.html",
		Run: func(cmd *cobra.Command, args []string) {
			// Prepare params
			params := &artransfer.ImportParams{}

			// Parse JSON parameters from flags
			if jsonParams, _ := cmd.Flags().GetString("parameters"); jsonParams != "" {
				if err := json.Unmarshal([]byte(jsonParams), params); err != nil {
					log.Fatal("Failed to parse JSON parameters: ", err)
				}
			}

			// If JSON file or JSON data is not provided in flags, try to get it from arguments
			if params.JsonFile == "" && params.JsonData == "" {
				if len(args) == 0 {
					log.Fatal(i18n.T("{{name}} is required", map[string]interface{}{"name": "FILENAME"}))
				}
				params.JsonFile = args[0]
			}

			// Parse flags to params
			if flagAssumeyes, err := cmd.Flags().GetBool("assumeyes"); err == nil {
				params.Assumeyes = flagAssumeyes
			}

			// Check if file exists if JsonFile is provided
			if params.JsonFile != "" {
				if _, err := os.Stat(params.JsonFile); errors.Is(err, os.ErrNotExist) {
					log.Fatal(i18n.T("{{name}} not found", map[string]interface{}{"name": i18n.T("File")}))
				}
			}

			// Run import
			artransfer.RunImportArtrans(app.Dao(), params)
		},
	}

	flagPV(importCmd, "assumeyes", "y", false, "Automatically answer yes for all questions.")
	flagPV(importCmd, "parameters", "p", "", "JSON format parameters for the import command.")

	return importCmd
}
