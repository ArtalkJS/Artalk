package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOktaConfiguration(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "artalk.yml")
	require.NoError(t, os.WriteFile(filename, []byte("locale: en\nauth:\n  okta:\n    issuer: https://file.example\n"), 0600))
	conf, err := NewFromFile(filename)
	require.NoError(t, err)
	require.False(t, conf.Auth.Okta.Enabled)
	require.Equal(t, "https://file.example", conf.Auth.Okta.Issuer)
	t.Setenv("ATK_AUTH_OKTA_ENABLED", "true")
	t.Setenv("ATK_AUTH_OKTA_CLIENT_ID", "client")
	t.Setenv("ATK_AUTH_OKTA_CLIENT_SECRET", "secret")
	t.Setenv("ATK_AUTH_OKTA_ISSUER", "https://tenant.okta.com/oauth2/default")
	conf, err = NewFromFile(filename)
	require.NoError(t, err)
	require.True(t, conf.Auth.Okta.Enabled)
	require.Equal(t, "client", conf.Auth.Okta.ClientID)
	require.Equal(t, "secret", conf.Auth.Okta.ClientSecret)
	require.Equal(t, "https://tenant.okta.com/oauth2/default", conf.Auth.Okta.Issuer)
	for _, field := range []string{"enabled", "client_id", "client_secret", "issuer"} {
		found := false
		for _, option := range OptionsMetaCache {
			if option.Path == "auth.okta."+field {
				require.True(t, option.AllowsSet)
				found = true
				break
			}
		}
		require.True(t, found, "missing settings metadata for %s", field)
	}
}
