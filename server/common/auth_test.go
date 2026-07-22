package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginGetUserToken(t *testing.T) {
	const key = "test-secret"
	user := entity.User{}
	user.ID = 42
	token, err := LoginGetUserToken(user, key, 60)
	require.NoError(t, err)

	claims, err := parseJWTToken(token, key)
	require.NoError(t, err)
	assert.Equal(t, uint(42), claims.UserID)
	assert.WithinDuration(t, time.Now(), claims.IssuedAt.Time, time.Second)
	assert.WithinDuration(t, time.Now().Add(time.Minute), claims.ExpiresAt.Time, time.Second)
}

func TestParseJWTTokenAcceptsLegacyClaims(t *testing.T) {
	const key = "test-secret"
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"user_id":42,"exp":4102444800,"iat":1600000000}`))
	unsigned := header + "." + payload
	mac := hmac.New(sha256.New, []byte(key))
	_, _ = mac.Write([]byte(unsigned))
	token := unsigned + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	claims, err := parseJWTToken(token, key)
	require.NoError(t, err)
	assert.Equal(t, uint(42), claims.UserID)
	assert.Equal(t, int64(1600000000), claims.IssuedAt.Unix())
}

func TestParseJWTTokenRejectsUnexpectedAlgorithm(t *testing.T) {
	const key = "test-secret"
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS384, jwt.MapClaims{
		"user_id": 42,
		"exp":     time.Now().Add(time.Minute).Unix(),
	}).SignedString([]byte(key))
	require.NoError(t, err)

	_, err = parseJWTToken(token, key)
	require.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("signing method %s is invalid", jwt.SigningMethodHS384.Alg()))
}

func TestParseJWTTokenRejectsExpiredToken(t *testing.T) {
	const key = "test-secret"
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": 42,
		"exp":     time.Now().Add(-time.Minute).Unix(),
	}).SignedString([]byte(key))
	require.NoError(t, err)

	_, err = parseJWTToken(token, key)
	require.Error(t, err)
	assert.ErrorIs(t, err, jwt.ErrTokenExpired)
}
