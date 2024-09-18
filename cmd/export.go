package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/artalkjs/artalk/v2/internal/artransfer"
	"github.com/artalkjs/artalk/v2/internal/i18n"
	"github.com/artalkjs/artalk/v2/internal/log"
	"github.com/artalkjs/artalk/v2/internal/utils"
	"github.com/spf13/cobra"
)

func NewExportCommand(app *ArtalkCmd) *cobra.Command {
	exportCmd := &cobra.Command{
		Use:     "export",
		Aliases: []string{},
		Short:   "Artransfer export",
		Long:    "\n# Artransfer - Export\n\n  See the documentation to learn more: https://artalk.js.org/guide/transfer.html",
		Run: func(cmd *cobra.Command, args []string) {
			jsonStr, err := artransfer.RunExportArtrans(app.Dao(), &artransfer.ExportParams{})
			if err != nil {
				log.Fatal(err)
			}

			if len(args) < 1 || args[0] == "" {
				// write to stdout
				fmt.Println(jsonStr)
			} else {
				filename := args[0]

				// make sure is abs path
				filename, err := filepath.Abs(filename)
				if err != nil {
					log.Fatal(err)
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
					log.Fatal(err)
				}

				// touch
				f, err := os.Create(filename)
				if err != nil {
					log.Fatal(err)
				}

				// >
				_, err2 := f.WriteString(jsonStr)
				if err2 != nil {
					log.Fatal(err2)
				}

				log.Info(i18n.T("Export complete") + ": " + filename)
			}
		},
	}

	return exportCmd
}
