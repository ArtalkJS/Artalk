package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/lib/core"
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Aliases: []string{"update"},
	Short:   "升级到最新版",
	Long:    "将 ArtalkGo 升级到最新版本，\n更新源为 GitHub Releases，\n更新需要重启 ArtalkGo 才能生效。",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// loadCore() // 装载核心
		core.LoadConfOnly(cfgFile, workDir)

		logrus.Info("从 GitHub Release 检索更新中...")

		latest, found, err := selfupdate.DetectLatest("ArtalkJS/ArtalkGo")
		if err != nil {
			logrus.Fatal("Error occurred while detecting version: ", err)
		}

		ignoreVersionCheck, _ := cmd.Flags().GetBool("force")
		if !ignoreVersionCheck {
			v := semver.MustParse(strings.TrimPrefix(lib.Version, "v"))
			if !found || latest.Version.LTE(v) {
				logrus.Println("当前已是最新版本 v" + v.String() + " 无需升级")
				return
			}
		}

		logrus.Info("发现新版本: v" + latest.Version.String())
		logrus.Info("正在下载更新中...")

		exe, err := os.Executable()
		if err != nil {
			logrus.Fatal("Could not locate executable path ", err)
		}

		if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
			logrus.Fatal("更新失败 ", err)
		}

		logrus.Println("更新完毕")
		fmt.Println("\n-------------------------------\n    v" +
			latest.Version.String() +
			"  Release Note\n" +
			"-------------------------------\n\n" +
			latest.ReleaseNotes)
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)

	flagPV(upgradeCmd, "force", "f", false, "Force upgrade ignore version comparison.")
}
