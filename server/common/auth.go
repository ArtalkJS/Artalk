package common

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type jwtCustomClaims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func LoginGetUserToken(user entity.User, key string, ttl int) (string, error) {
	// Set custom claims
	claims := &jwtCustomClaims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),                                       // 签发时间
			ExpiresAt: time.Now().Add(time.Second * time.Duration(ttl)).Unix(), // 过期时间
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

	jwt, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != "HS256" {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}

		return []byte(app.Conf().AppKey), nil // 密钥
	})
	if err != nil {
		return jwtCustomClaims{}, err
	}

	claims := jwtCustomClaims{}
	tmp, _ := json.Marshal(jwt.Claims)
	_ = json.Unmarshal(tmp, &claims)

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
	if user.TokenValidFrom.Valid && user.TokenValidFrom.Time.After(time.Unix(claims.IssuedAt, 0)) {
		return entity.User{}, ErrTokenInvalidFromDate
	}

	return user, nil
}
