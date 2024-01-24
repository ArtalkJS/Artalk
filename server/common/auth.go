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
	UserID  uint   `json:"user_id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
	jwt.StandardClaims
}

func LoginGetUserToken(user entity.User, key string, ttl int) string {
	// Set custom claims
	claims := &jwtCustomClaims{
		UserID:  user.ID,
		Name:    user.Name,
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(ttl)).Unix(), // 过期时间
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(key))
	if err != nil {
		return ""
	}

	return t
}

func GetJwtInstanceByReq(app *core.App, c *fiber.Ctx) *jwt.Token {
	token := c.Query("token")
	if token == "" {
		token = c.FormValue("token")
	}
	if token == "" {
		token = c.Get(fiber.HeaderAuthorization)
		token = strings.TrimPrefix(token, "Bearer ")
	}
	if token == "" {
		return nil
	}

	jwt, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != "HS256" {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}

		return []byte(app.Conf().AppKey), nil // 密钥
	})
	if err != nil {
		return nil
	}

	return jwt
}

func GetUserByJwt(app *core.App, jwt *jwt.Token) entity.User {
	if jwt == nil {
		return entity.User{}
	}

	claims := jwtCustomClaims{}
	tmp, _ := json.Marshal(jwt.Claims)
	_ = json.Unmarshal(tmp, &claims)

	user := app.Dao().FindUserByID(claims.UserID)

	return user
}

func GetUserByReq(app *core.App, c *fiber.Ctx) entity.User {
	jwt := GetJwtInstanceByReq(app, c)
	user := GetUserByJwt(app, jwt)

	return user
}
