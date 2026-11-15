package auth

import (
	"fmt"
	"strings"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/markbates/goth"
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

func TestSocialProfileFallbacks(t *testing.T) {
	for _, tc := range []struct {
		name                          string
		profile                       goth.User
		wantEmail, wantName, wantLink string
	}{
		{"GitHub blog", goth.User{Provider: "github", UserID: "42", Name: "Reader", Email: "reader@example.org", RawData: map[string]any{"blog": "https://reader.example.org", "html_url": "https://github.com/reader"}}, "reader@example.org", "Reader", "https://reader.example.org"},
		{"GitHub profile", goth.User{Provider: "github", UserID: "42", Email: "reader@example.org", RawData: map[string]any{"html_url": "https://github.com/reader"}}, "reader@example.org", "reader", "https://github.com/reader"},
		{"Steam without email", goth.User{Provider: "steam", UserID: "42"}, "42@steam.com", "42", ""},
		{"Provider without profile details", goth.User{Provider: "discord", UserID: "42"}, "42@discord.com", "42", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			user, err := GetSocialUser(tc.profile, &config.AuthConf{})
			require.NoError(t, err)
			require.Equal(t, tc.wantEmail, user.Email)
			require.Equal(t, tc.wantName, user.Name)
			require.Equal(t, tc.wantLink, user.Link)
			require.Equal(t, tc.profile.UserID, user.RemoteUID)
		})
	}
}
