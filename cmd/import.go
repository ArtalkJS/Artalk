package cmd

import (
	"github.com/ArtalkJS/ArtalkGo/lib/importer"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:     "import",
	Aliases: []string{},
	Short:   "数据搬家 - 导入",
	Long: "\n# 数据搬家 - 导入\n\n  从其他评论系统迁移数据到 ArtalkGo\n" + `
- 例如：artalk-go import typecho [参数]
- 文档：https://artalk.js.org/guide/transfer.html
`,
}

func init() {
	rootCmd.AddCommand(importCmd)

	for _, item := range importer.Supports {
		imp := importer.GetImporterInfo(item)
		subCmd := &cobra.Command{
			Use:   imp.Name + " [参数]",
			Short: imp.Desc,
			Run: func(cmd *cobra.Command, args []string) {
				importer.RunByName(imp.Name, args)
			},
		}

		importCmd.AddCommand(subCmd)
	}
}
