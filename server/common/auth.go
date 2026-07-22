package common

import (
	"fmt"
	"strings"
	"time"

	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type jwtCustomClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func LoginGetUserToken(user entity.User, key string, ttl int) (string, error) {
	// Set custom claims
	claims := &jwtCustomClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                       // 签发时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(ttl))), // 过期时间
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return t, nil
}

var ErrTokenNotProvided = fmt.Errorf("token not provided")
var ErrTokenUserNotFound = fmt.Errorf("user not found")
var ErrTokenInvalidFromDate = fmt.Errorf("token is invalid starting from a certain date")

func GetTokenByReq(c *fiber.Ctx) string {
	token := c.Query("token")
	if token == "" {
		token = c.FormValue("token")
	}
	if token == "" {
		token = c.Get(fiber.HeaderAuthorization)
		token = strings.TrimPrefix(token, "Bearer ")
	}
	return token
}

func GetJwtDataByReq(app *core.App, c *fiber.Ctx) (jwtCustomClaims, error) {
	token := GetTokenByReq(c)
	if token == "" {
		return jwtCustomClaims{}, ErrTokenNotProvided
	}
	return parseJWTToken(token, app.Conf().AppKey)
}

func parseJWTToken(token, key string) (jwtCustomClaims, error) {
	claims := jwtCustomClaims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}

		return []byte(key), nil // 密钥
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return jwtCustomClaims{}, err
	}

	return claims, nil
}

func GetUserByReq(app *core.App, c *fiber.Ctx) (entity.User, error) {
	claims, err := GetJwtDataByReq(app, c)
	if err != nil {
		return entity.User{}, err
	}

	user := app.Dao().FindUserByID(claims.UserID)
	if user.IsEmpty() {
		return entity.User{}, ErrTokenUserNotFound
	}

	// check tokenValidFrom
	issuedAt := int64(0)
	if claims.IssuedAt != nil {
		issuedAt = claims.IssuedAt.Unix()
	}
	if user.TokenValidFrom.Valid && issuedAt < user.TokenValidFrom.Time.Unix() {
		return entity.User{}, ErrTokenInvalidFromDate
	}

	return user, nil
}
