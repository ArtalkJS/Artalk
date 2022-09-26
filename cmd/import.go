package cmd

import (
	"errors"
	"os"

	"github.com/ArtalkJS/ArtalkGo/lib/artransfer"
	"github.com/ArtalkJS/ArtalkGo/lib/core"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:     "import <数据行囊文件路径>",
	Aliases: []string{},
	Short:   "数据迁移 - 迁入",
	Long: "\n# 数据迁移 - 迁入\n\n  从其他评论系统迁移数据到 Artalk\n" + `
- 导入前需要使用转换工具 Artransfer 将其他评论数据转为 Artrans 格式
- 文档：https://artalk.js.org/guide/transfer.html
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		core.LoadCore(cfgFile, workDir) // 装载核心

		parcelFile := args[0]
		if _, err := os.Stat(parcelFile); errors.Is(err, os.ErrNotExist) {
			logrus.Fatal("`数据行囊` 文件不存在，请检查路径是否正确")
		}

		payload := args[1:]
		payload = append(payload, "json_file:"+parcelFile)

		// 导入 Artrans
		artransfer.RunImportArtrans(payload)

		logrus.Info("导入结束")
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	flagPV(importCmd, "assumeyes", "y", false, "Automatically answer yes for all questions.")
}
