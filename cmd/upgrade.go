package cmd

import (
	"fmt"
	"os"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Aliases: []string{"update"},
	Short:   "升级 ArtalkGo 到最新版本",
	Long:    "将 ArtalkGo 升级到最新版本，\n更新源为 GitHub Releases，\n更新需要重启 ArtalkGo 才能生效。",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("从 GitHub Release 检索更新中...")

		latest, found, err := selfupdate.DetectLatest("ArtalkJS/ArtalkGo")
		if err != nil {
			logrus.Fatal("Error occurred while detecting version: ", err)
		}

		v := semver.MustParse(lib.Version[1:])
		if !found || latest.Version.LTE(v) {
			logrus.Println("当前已是最新版本 v" + v.String() + " 无需升级")
			return
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
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
