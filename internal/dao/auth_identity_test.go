package dao_test

import (
	"errors"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/test"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestSocialIdentityLifecycleStaysProviderScoped(t *testing.T) {
	app, err := test.NewTestApp()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	d := app.Dao()
	user := d.FindUserByID(1001)
	first := entity.AuthIdentity{UserID: user.ID, Provider: "first", RemoteUID: "remote", Token: "shared-value"}
	second := entity.AuthIdentity{UserID: 1002, Provider: "second", RemoteUID: "remote", Token: first.Token}
	require.NoError(t, d.CreateAuthIdentity(&first))
	require.NoError(t, d.CreateAuthIdentity(&second))
	require.Equal(t, first.ID, d.FindAuthIdentityByToken(first.Provider, first.Token).ID)
	require.Equal(t, second.ID, d.FindAuthIdentityByUserID(second.Provider, second.UserID).ID)
	first.Token = "rotated-token"
	require.NoError(t, d.UpdateAuthIdentity(&first))
	require.True(t, d.FindAuthIdentityByToken(first.Provider, second.Token).IsEmpty())
	require.Equal(t, first.ID, d.FindAuthIdentityByToken(first.Provider, first.Token).ID)
	require.NoError(t, d.DB().Callback().Update().Before("gorm:update").Register("test:identity_failure", func(tx *gorm.DB) { tx.AddError(errors.New("identity write failure")) }))
	first.Token = "unsaved-token"
	require.Error(t, d.UpdateAuthIdentity(&first))
	require.NoError(t, d.DB().Callback().Update().Remove("test:identity_failure"))
	require.Equal(t, "rotated-token", d.FindAuthIdentityByUserID(first.Provider, first.UserID).Token)
	require.NoError(t, d.DelUser(&user))
	require.True(t, d.FindAuthIdentityByUserID(first.Provider, user.ID).IsEmpty())
	require.Equal(t, second.ID, d.FindAuthIdentityByToken(second.Provider, second.Token).ID)
}
