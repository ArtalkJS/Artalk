package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Aliases: []string{"update"},
	Short:   "Upgrade to the latest version",
	Long:    "Upgrade Artalk to the latest version, \n update source is GitHub Releases, \n update need to restart Artalk to take effect.",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// loadCore()
		core.LoadConfOnly(cfgFile, workDir)

		logrus.Info(i18n.T("Checking for updates") + "...")

		latest, found, err := selfupdate.DetectLatest("ArtalkJS/Artalk")
		if err != nil {
			logrus.Fatal("Error occurred while detecting version: ", err)
		}

		ignoreVersionCheck, _ := cmd.Flags().GetBool("force")
		if !ignoreVersionCheck {
			v := semver.MustParse(strings.TrimPrefix(config.Version, "v"))
			if !found || latest.Version.LTE(v) {
				logrus.Println(i18n.T("Current version is the latest") + " (v" + v.String() + ")")
				return
			}
		}

		logrus.Info(i18n.T("New version available") + ": v" + latest.Version.String())
		logrus.Info(i18n.T("Downloading") + "...")

		exe, err := os.Executable()
		if err != nil {
			logrus.Fatal("Could not locate executable path ", err)
		}

		if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
			logrus.Fatal(i18n.T("Update failed")+" ", err)
		}

		logrus.Println(i18n.T("Update complete"))
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
