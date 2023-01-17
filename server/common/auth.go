package common

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/query"
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

func LoginGetUserToken(user entity.User) string {
	// Set custom claims
	claims := &jwtCustomClaims{
		UserID:  user.ID,
		Name:    user.Name,
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(config.Instance.LoginTimeout)).Unix(), // 过期时间
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(config.Instance.AppKey))
	if err != nil {
		return ""
	}

	return t
}

func GetJwtStrByReqCookie(c *fiber.Ctx) string {
	if !config.Instance.Cookie.Enabled {
		return ""
	}
	cookie := c.Cookies(config.COOKIE_KEY_ATK_AUTH)
	return cookie
}

func GetJwtInstanceByReq(c *fiber.Ctx) *jwt.Token {
	token := c.Query("token")
	if token == "" {
		token = c.FormValue("token")
	}
	if token == "" {
		token = c.Get(fiber.HeaderAuthorization)
		token = strings.TrimPrefix(token, "Bearer ")
	}
	if token == "" {
		token = GetJwtStrByReqCookie(c)
	}
	if token == "" {
		return nil
	}

	jwt, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != "HS256" {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}

		return []byte(config.Instance.AppKey), nil // 密钥
	})
	if err != nil {
		return nil
	}

	return jwt
}

func GetUserByJwt(jwt *jwt.Token) entity.User {
	if jwt == nil {
		return entity.User{}
	}

	claims := jwtCustomClaims{}
	tmp, _ := json.Marshal(jwt.Claims)
	_ = json.Unmarshal(tmp, &claims)

	user := query.FindUserByID(claims.UserID)

	return user
}

func GetUserByReq(c *fiber.Ctx) entity.User {
	jwt := GetJwtInstanceByReq(c)
	user := GetUserByJwt(jwt)

	return user
}
