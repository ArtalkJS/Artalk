package auth

import (
	"fmt"
	"time"

	"github.com/ArtalkJS/Artalk/internal/dao"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/utils"
)

func RegisterSocialUser(dao *dao.Dao, u SocialUser) (entity.AuthIdentity, error) {
	if u.Name == "" {
		return entity.AuthIdentity{}, fmt.Errorf("name is required")
	}
	if u.Email == "" {
		return entity.AuthIdentity{}, fmt.Errorf("email is required")
	}
	if !utils.ValidateEmail(u.Email) {
		return entity.AuthIdentity{}, fmt.Errorf("email is invalid")
	}

	// Create user if not exists
	user, err := dao.FindCreateUser(u.Name, u.Email, u.Link)
	if err != nil {
		return entity.AuthIdentity{}, err
	}

	// Store user auth identity in db
	now := time.Now()
	authIdentity := entity.AuthIdentity{
		Provider:    u.Provider,
		RemoteUID:   u.RemoteUID,
		UserID:      user.ID,
		Token:       u.AccessToken,
		ConfirmedAt: &now,
		ExpiresAt:   &u.ExpiresAt,
	}
	if err := dao.CreateAuthIdentity(&authIdentity); err != nil {
		return entity.AuthIdentity{}, err
	}

	return authIdentity, nil
}
