package cmd

import (
	"os"

	"github.com/ArtalkJS/ArtalkGo/lib/core"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen <类型> <目标路径>",
	Short: "生成一些内容",
	Long:  "生成一些内容\n例如：artalk-go gen config ./artalk-go.yml",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		// 工作目录
		if workDir != "" {
			if err := os.Chdir(workDir); err != nil {
				logrus.Fatal("工作目录切换错误 ", err)
			}
		}

		var (
			specificPath string
			isForce      bool
		)
		if len(args) > 1 {
			specificPath = args[1]
		}
		isForce, _ = cmd.Flags().GetBool("force")

		core.Gen(args[0], specificPath, isForce)
	},
}

func init() {
	rootCmd.AddCommand(genCmd)

	flagPV(genCmd, "force", "f", false, "Force overwrite an existing file")
}
