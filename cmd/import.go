package cmd

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/ArtalkJS/ArtalkGo/lib/artransfer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:     "import <数据行囊文件路径>",
	Aliases: []string{},
	Short:   "数据搬家 - 迁入",
	Long: "\n# 数据搬家 - 迁入\n\n  从其他评论系统迁移数据到 ArtalkGo\n" + `
- 例如：artalk-go import typecho [参数]
- 文档：https://artalk.js.org/guide/transfer.html
`,
	Run: func(cmd *cobra.Command, args []string) {
		parcelFile := args[0]
		if _, err := os.Stat(parcelFile); errors.Is(err, os.ErrNotExist) {
			logrus.Fatal("`数据行囊` 文件不存在，请检查路径是否正确")
		}

		buf, err := ioutil.ReadFile(parcelFile)
		if err != nil {
			logrus.Fatal("数据行囊文件打开失败：", err)
		}

		content := string(buf)
		basic := artransfer.GetBasicParamsFrom(args[1:])
		basic.UrlResolver = false // 默认关闭 URL 解析器：因为 pageKey 是完整，且站点隔离开
		artransfer.ImportArtransByStr(basic, content)

		logrus.Info("导入结束")
	},
	Args:       cobra.MinimumNArgs(1),
	ArgAliases: artransfer.GetSupportNames(),
}

func init() {
	rootCmd.AddCommand(importCmd)

	for _, item := range artransfer.Supports {
		imp := artransfer.GetImporterInfo(item)
		subCmd := &cobra.Command{
			Use:   imp.Name + " [参数]",
			Short: imp.Desc,
			Run: func(cmd *cobra.Command, args []string) {
				artransfer.Assumeyes, _ = cmd.Flags().GetBool("assumeyes")
				artransfer.RunByName(imp.Name, args)
			},
		}

		importCmd.AddCommand(subCmd)
	}

	flagPV(importCmd, "assumeyes", "y", false, "Automatically answer yes for all questions.")
}
