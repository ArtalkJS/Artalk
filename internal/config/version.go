package config

import (
	"runtime/debug"

	"github.com/samber/lo"
)

// The version of Artalk
//
// Which is automatically set by the CI release workflow
const Version = "v2.9.1"

// The commit hash from which the binary was built (optional)
//
// This value is set by the build script:
// `go build -ldflags "-X 'github.com/artalkjs/artalk/v2/internal/config.CommitHash=$(git rev-parse --short HEAD)'"`
func CommitHash() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	setting, found := lo.Find(info.Settings, func(s debug.BuildSetting) bool { return s.Key == "vcs.revision" })
	if !found {
		return ""
	}
	if len(setting.Value) >= 7 {
		return setting.Value[:7]
	}
	return setting.Value
}

// Get the version string
// (format is "Version/CommitHash" or only "Version" if CommitHash is empty)
func VersionString() string {
	var str = Version
	if CommitHash() != "" {
		str += "/" + CommitHash()
	}
	return str
}
