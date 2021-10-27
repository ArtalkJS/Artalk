package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/lib/artransfer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:     "export",
	Aliases: []string{},
	Short:   "数据搬家 - 迁出",
	Long: "\n# 数据搬家 - 迁出\n\n  将所有数据从 ArtalkGo 导出，用作备份，或迁移至其他地方\n  打包所有数据并导出成 “ArtalkGo 数据行囊 (Artrans)”，为数据搬家做准备\n" + `
- 重新导入 ArtalkGo，可执行: artalk-go import <数据行囊文件路径>
- 文档：https://artalk.js.org/guide/transfer.html
`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonStr, err := artransfer.ExportArtransString()
		if err != nil {
			logrus.Fatal(err)
		}

		outputFile := cmd.Flag("output").Value.String()
		if outputFile == "" {
			// write to stdout
			fmt.Println(jsonStr)
		} else {
			// mkdir -p
			if err := lib.EnsureDir(filepath.Dir(outputFile)); err != nil {
				logrus.Fatal(err)
			}

			// touch
			f, err := os.Create(outputFile)
			if err != nil {
				logrus.Fatal(err)
			}

			// >
			_, err2 := f.WriteString(jsonStr)
			if err2 != nil {
				logrus.Fatal(err2)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	flagP(exportCmd, "output", "o", "", "Write to file instead of stdout")
}
