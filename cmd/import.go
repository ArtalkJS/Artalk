package cmd

import (
	"strings"

	"github.com/ArtalkJS/ArtalkGo/lib/importer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:     "import",
	Version: rootCmd.Version,
	Aliases: []string{},
	Short:   "数据导入",
	Long:    rootCmd.Long,
	Run: func(cmd *cobra.Command, args []string) {
		dataType, _ := cmd.Flags().GetString("type")
		if strings.TrimSpace(dataType) == "" {
			logrus.Fatal("请指定参数 `--type`，帮助文档加 `-h`")
		}

		isSupport := false
		for _, t := range importer.GetSupportTypes() {
			if strings.EqualFold(t, dataType) {
				isSupport = true
			}
		}

		if !isSupport {
			logrus.Fatal("不支持该数据类型")
		}

		importer.RunByName(dataType, args)
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	flagPV(importCmd, "type", "", "", "数据类型 ["+strings.Join(importer.GetSupportTypes(), ",")+"]")
	flagPV(importCmd, "port", "", 23366, "监听端口")
}
