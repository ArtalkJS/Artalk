package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

func NewUpgradeCommand() *cobra.Command {
	upgradeCmd := &cobra.Command{
		Use:     "upgrade",
		Aliases: []string{"update"},
		Short:   "Upgrade to the latest version",
		Long:    "Upgrade Artalk to the latest version, \n update source is GitHub Releases, \n update need to restart Artalk to take effect.",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			log.Info(i18n.T("Checking for updates") + "...")

			latest, found, err := selfupdate.DetectLatest("ArtalkJS/Artalk")
			if err != nil {
				log.Fatal("Error occurred while detecting version: ", err)
			}

			ignoreVersionCheck, _ := cmd.Flags().GetBool("force")
			if !ignoreVersionCheck {
				v := semver.MustParse(strings.TrimPrefix(config.Version, "v"))
				if !found || latest.Version.LTE(v) {
					log.Info(i18n.T("Current version is the latest") + " (v" + v.String() + ")")
					return
				}
			}

			log.Info(i18n.T("New version available") + ": v" + latest.Version.String())
			log.Info(i18n.T("Downloading") + "...")

			exe, err := os.Executable()
			if err != nil {
				log.Fatal("Could not locate executable path ", err)
			}

			if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
				log.Fatal(i18n.T("Update failed")+" ", err)
			}

			log.Info(i18n.T("Update complete"))
			fmt.Printf("\n"+
				"-------------------------------\n"+
				"    v%s  Release Note\n"+
				"-------------------------------\n\n"+
				"%s\n",
				latest.Version.String(), latest.ReleaseNotes)
		},
		Annotations: map[string]string{
			BootModeKey: MODE_MINI_BOOT,
		},
	}

	flagPV(upgradeCmd, "force", "f", false, "Force upgrade ignore version comparison.")

	return upgradeCmd
}
