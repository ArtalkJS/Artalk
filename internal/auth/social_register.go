package auth

import (
	"fmt"
	"time"

	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/utils"
)

func RegisterSocialUser(dao *dao.Dao, u SocialUser) (entity.AuthIdentity, error) {
	if u.Name == "" {
		return entity.AuthIdentity{}, fmt.Errorf("cannot fetch user name from social identity provider")
	}
	if u.Email == "" {
		return entity.AuthIdentity{}, fmt.Errorf("cannot fetch user email from social identity provider")
	}
	if !utils.ValidateEmail(u.Email) {
		return entity.AuthIdentity{}, fmt.Errorf("email is invalid which fetched from social identity provider")
	}

	// Create user if not exists
	user, err := dao.FindCreateUser(u.Name, u.Email, u.Link)
	if err != nil {
		return entity.AuthIdentity{}, err
	}

	var expiresAt *time.Time
	if !u.ExpiresAt.IsZero() {
		expiresAt = &u.ExpiresAt
	}

	// Store user auth identity in db
	now := time.Now()
	authIdentity := entity.AuthIdentity{
		Provider:    u.Provider,
		RemoteUID:   u.RemoteUID,
		UserID:      user.ID,
		Token:       u.AccessToken,
		ConfirmedAt: &now,
		ExpiresAt:   expiresAt,
	}
	if err := dao.CreateAuthIdentity(&authIdentity); err != nil {
		return entity.AuthIdentity{}, err
	}

	return authIdentity, nil
}
