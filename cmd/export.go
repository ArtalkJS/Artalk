package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/ArtalkJS/Artalk/internal/artransfer"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:     "export",
	Aliases: []string{},
	Short:   "Artransfer - Export",
	Long:    "\n# Artransfer - Export\n\n  See the documentation to learn more: https://artalk.js.org/guide/transfer.html",
	Run: func(cmd *cobra.Command, args []string) {
		core.LoadCore(cfgFile, workDir)

		jsonStr, err := artransfer.ExportArtransString()
		if err != nil {
			logrus.Fatal(err)
		}

		if len(args) < 1 || args[0] == "" {
			// write to stdout
			fmt.Println(jsonStr)
		} else {
			filename := args[0]

			// make sure is abs path
			filename, err := filepath.Abs(filename)
			if err != nil {
				logrus.Fatal(err)
			}

			// check dir
			stat, err := os.Stat(filename)
			if err == nil {
				if stat.IsDir() {
					filename = path.Join(filename, "backup-"+time.Now().Format("20060102-150405")+".artrans")
				}
			}

			// mkdir -p
			if err := utils.EnsureDir(filepath.Dir(filename)); err != nil {
				logrus.Fatal(err)
			}

			// touch
			f, err := os.Create(filename)
			if err != nil {
				logrus.Fatal(err)
			}

			// >
			_, err2 := f.WriteString(jsonStr)
			if err2 != nil {
				logrus.Fatal(err2)
			}

			logrus.Info(i18n.T("Export complete") + ": " + filename)
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
