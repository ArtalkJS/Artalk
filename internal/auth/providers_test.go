package auth

import (
	"fmt"
	"strings"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/stretchr/testify/require"
)

func TestEnabledProvidersHaveLoginMetadata(t *testing.T) {
	names := []string{"github", "gitlab", "google", "mastodon", "twitter", "facebook", "discord", "steam", "apple", "microsoft", "wechat", "tiktok", "slack", "line", "patreon", "auth0"}
	yaml := "auth:\n  callback: https://comments.example.org/api/v2/auth/{provider}/callback\n"
	for _, name := range names {
		yaml += fmt.Sprintf("  %s:\n    enabled: true\n", name)
	}
	conf, err := config.NewFromYAML("", []byte(yaml))
	require.NoError(t, err)
	providers := GetProviders(conf)
	require.Len(t, providers, len(names))
	for i, info := range GetProviderInfo(conf, providers) {
		want := names[i]
		if want == "microsoft" {
			want = "microsoftonline"
			require.Equal(t, "Microsoft", info.Label)
		}
		require.Equal(t, want, info.Name)
		require.Equal(t, "/api/v2/auth/"+want, info.Path)
		require.True(t, strings.HasPrefix(info.Icon, "data:image/svg+xml;base64,"))
	}
}
