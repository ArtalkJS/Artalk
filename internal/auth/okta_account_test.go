package auth_test

import (
	"testing"

	"github.com/artalkjs/artalk/v2/internal/auth"
	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/test"
	"github.com/markbates/goth"
	"github.com/stretchr/testify/require"
)

func TestOktaAccountsStaySeparate(t *testing.T) {
	app, err := test.NewTestApp()
	require.NoError(t, err)
	defer app.Cleanup()
	existing, err := app.Dao().NewUser("Alice", "alice@example.test", "")
	require.NoError(t, err)
	existing.IsAdmin = true
	require.NoError(t, app.Dao().DB().Save(&existing).Error)
	conf := &config.AuthConf{}
	profile := goth.User{Provider: "okta", UserID: "same-subject", Name: existing.Name, Email: existing.Email, RawData: map[string]interface{}{"email_verified": false}}
	previousID := existing.ID
	for _, issuer := range []string{"https://tenant.okta.com", "https://tenant.okta.com/oauth2/custom", "https://other.okta.com"} {
		conf.Okta.Issuer = issuer
		user, err := auth.GetSocialUser(profile, conf)
		require.NoError(t, err)
		require.True(t, app.Dao().FindAuthIdentityByRemoteUID("okta", user.RemoteUID).IsEmpty())
		identity, err := auth.RegisterSocialUser(app.Dao(), user)
		require.NoError(t, err)
		require.NotEqual(t, existing.ID, identity.UserID)
		require.NotEqual(t, previousID, identity.UserID)
		require.False(t, app.Dao().FindUserByID(identity.UserID).IsAdmin)
		// A later login uses the same identity even if the profile changes.
		profile.Name = "Renamed Alice"
		profile.Email = "new@example.test"
		again, err := auth.GetSocialUser(profile, conf)
		require.NoError(t, err)
		require.Equal(t, identity.ID, app.Dao().FindAuthIdentityByRemoteUID("okta", again.RemoteUID).ID)
		profile.Name = existing.Name
		profile.Email = existing.Email
		previousID = identity.UserID
	}
}
