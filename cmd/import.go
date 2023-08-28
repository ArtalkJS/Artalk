package cmd

import (
	"errors"
	"os"

	"github.com/ArtalkJS/Artalk/internal/artransfer"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/spf13/cobra"
)

func NewImportCommand(app *core.App) *cobra.Command {
	importCmd := &cobra.Command{
		Use:     "import <FILENAME>",
		Aliases: []string{},
		Short:   "Artransfer import",
		Long:    "\n# Artransfer - Import\n\n  See the documentation to learn more: https://artalk.js.org/guide/transfer.html",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			parcelFile := args[0]
			if _, err := os.Stat(parcelFile); errors.Is(err, os.ErrNotExist) {
				log.Fatal(i18n.T("{{name}} not found", map[string]interface{}{"name": i18n.T("File")}))
			}

			payload := args[1:]
			payload = append(payload, "json_file:"+parcelFile)

			// import Artrans
			artransfer.RunImportArtrans(app.Dao(), artransfer.ArrToImportParams(payload))

			log.Info(i18n.T("Import complete"))
		},
	}

	flagPV(importCmd, "assumeyes", "y", false, "Automatically answer yes for all questions.")

	return importCmd
}
