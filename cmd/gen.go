package cmd

import (
	"os"

	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen <TYPE> <DEST>",
	Short: "A collection of several useful generators",
	Long:  "Generate some content\ne.g. `artalk gen config ./artalk.yml`",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		// change working directory
		if workDir != "" {
			if err := os.Chdir(workDir); err != nil {
				logrus.Fatal("change working directory error ", err)
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
