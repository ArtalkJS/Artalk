package auth_test

import (
	"errors"
	"testing"
	"time"

	"github.com/artalkjs/artalk/v2/internal/auth"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/test"
	"github.com/markbates/goth"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestSocialRegistrationRejectsIncompleteProfilesAndStorageFailures(t *testing.T) {
	app, err := test.NewTestApp()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	for _, profile := range []goth.User{
		{Email: "reader@example.org"},
		{Name: "Reader"},
		{Name: "Reader", Email: "not-an-email"},
	} {
		identity, err := auth.RegisterSocialUser(app.Dao(), auth.SocialUser{User: profile})
		require.Error(t, err)
		require.True(t, identity.IsEmpty())
	}
	profile := auth.SocialUser{User: goth.User{
		Provider: "github", UserID: "remote-reader", Name: "Reader", Email: "reader@example.org",
		AccessToken: "test-token", ExpiresAt: time.Now().UTC().Add(time.Hour).Truncate(time.Second),
	}, RemoteUID: "remote-reader"}
	for _, table := range []string{"atk_users", "atk_auth_identities"} {
		t.Run(table, func(t *testing.T) {
			const callback = "test:social_registration_failure"
			require.NoError(t, app.Dao().DB().Callback().Create().Before("gorm:create").Register(callback, func(tx *gorm.DB) {
				if tx.Statement.Table == table {
					tx.AddError(errors.New("registration storage unavailable"))
				}
			}))
			t.Cleanup(func() { require.NoError(t, app.Dao().DB().Callback().Create().Remove(callback)) })
			identity, err := auth.RegisterSocialUser(app.Dao(), profile)
			require.ErrorContains(t, err, "registration storage unavailable")
			require.True(t, identity.IsEmpty())
			var count int64
			require.NoError(t, app.Dao().DB().Model(&entity.AuthIdentity{}).Where("remote_uid = ?", profile.RemoteUID).Count(&count).Error)
			require.Zero(t, count)
		})
	}
	identity, err := auth.RegisterSocialUser(app.Dao(), profile)
	require.NoError(t, err)
	require.Equal(t, profile.ExpiresAt, *identity.ExpiresAt)
	require.Equal(t, profile.AccessToken, app.Dao().FindAuthIdentityByRemoteUID(profile.Provider, profile.RemoteUID).Token)
	require.False(t, app.Dao().FindUserByID(identity.UserID).IsAdmin)
}
