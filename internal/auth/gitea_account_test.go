package auth_test

import (
	"testing"

	"github.com/artalkjs/artalk/v2/internal/auth"
	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/test"
	"github.com/markbates/goth"
	"github.com/stretchr/testify/require"
)

func TestSelfHostedGiteaDoesNotClaimExistingAccount(t *testing.T) {
	app, err := test.NewTestApp()
	require.NoError(t, err)
	defer app.Cleanup()
	existing, err := app.Dao().NewUser("Alice", "alice@example.test", "")
	require.NoError(t, err)
	existing.IsAdmin = true
	require.NoError(t, app.Dao().DB().Save(&existing).Error)
	profile := goth.User{Provider: "gitea", UserID: "42", Name: existing.Name, Email: existing.Email}
	conf := &config.Config{}
	legacy, err := auth.GetSocialUser(profile, conf.Auth.Gitea.SiteURL)
	require.NoError(t, err)
	oldIdentity, err := auth.RegisterSocialUser(app.Dao(), legacy)
	require.NoError(t, err)
	require.Equal(t, existing.ID, oldIdentity.UserID)
	previousID := existing.ID
	for _, site := range []string{"https://git.example/a", "https://git.example/b"} {
		conf.Auth.Gitea.SiteURL = site
		user, err := auth.GetSocialUser(profile, conf.Auth.Gitea.SiteURL)
		require.NoError(t, err)
		require.True(t, app.Dao().FindAuthIdentityByRemoteUID("gitea", user.RemoteUID).IsEmpty())
		identity, err := auth.RegisterSocialUser(app.Dao(), user)
		require.NoError(t, err)
		require.NotEqual(t, existing.ID, identity.UserID)
		require.NotEqual(t, previousID, identity.UserID)
		require.False(t, app.Dao().FindUserByID(identity.UserID).IsAdmin)
		require.Equal(t, identity.ID, app.Dao().FindAuthIdentityByRemoteUID("gitea", user.RemoteUID).ID)
		previousID = identity.UserID
	}
	require.Equal(t, oldIdentity.ID, app.Dao().FindAuthIdentityByRemoteUID("gitea", "42").ID)
}
