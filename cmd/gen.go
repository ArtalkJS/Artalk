package cmd

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/spf13/cobra"
)

func NewGenCommand(app *core.App) *cobra.Command {
	genCmd := &cobra.Command{
		Use:   "gen <TYPE> <DEST>",
		Short: "A collection of several useful generators",
		Long:  "Generate some content\ne.g. `artalk gen config ./artalk.yml`",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			// change working directory
			// if workDir != "" {
			// 	if err := os.Chdir(workDir); err != nil {
			// 		log.Fatal("change working directory error ", err)
			// 	}
			// }

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
		Annotations: map[string]string{
			BootModeKey: MODE_MINI_BOOT,
		},
	}

	flagPV(genCmd, "force", "f", false, "Force overwrite an existing file")

	return genCmd
}
