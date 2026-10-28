package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGiteaSiteURLEnvironment(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "artalk.yml")
	require.NoError(t, os.WriteFile(filename, []byte("locale: en\nauth:\n  gitea:\n    site_url: https://file.example\n"), 0600))
	t.Setenv("ATK_AUTH_GITEA_ENABLED", "true")
	t.Setenv("ATK_AUTH_GITEA_SITE_URL", "http://localhost:3000/git/")
	conf, err := NewFromFile(filename)
	require.NoError(t, err)
	require.True(t, conf.Auth.Gitea.Enabled)
	require.Equal(t, "http://localhost:3000/git/", conf.Auth.Gitea.SiteURL)
	for _, option := range OptionsMetaCache {
		if option.Path == "auth.gitea.site_url" {
			require.True(t, option.AllowsSet)
			return
		}
	}
	t.Fatal("site_url is missing from settings metadata")
}
