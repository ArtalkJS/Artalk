package cmd

import (
	"errors"
	"os"

	"github.com/ArtalkJS/Artalk/internal/artransfer"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:     "import <FILENAME>",
	Aliases: []string{},
	Short:   "Artransfer - Import",
	Long:    "\n# Artransfer - Import\n\n  See the documentation to learn more: https://artalk.js.org/guide/transfer.html",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		core.LoadCore(cfgFile, workDir) // load core

		parcelFile := args[0]
		if _, err := os.Stat(parcelFile); errors.Is(err, os.ErrNotExist) {
			logrus.Fatal(i18n.T("{{name}} not found", map[string]interface{}{"name": i18n.T("File")}))
		}

		payload := args[1:]
		payload = append(payload, "json_file:"+parcelFile)

		// import Artrans
		artransfer.RunImportArtrans(payload)

		logrus.Info(i18n.T("Import complete"))
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	flagPV(importCmd, "assumeyes", "y", false, "Automatically answer yes for all questions.")
}
