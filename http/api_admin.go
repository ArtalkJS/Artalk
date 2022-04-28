package http

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func AdminOnly(c echo.Context) (isOK bool, resp error) {
	if !CheckIsAdminReq(c) {
		return false, RespError(c, "需要验证管理员身份", Map{"need_login": true})
	}

	return true, nil
}

func LoginGetUserToken(user model.User) string {
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

func GetJwtInstanceByReq(c echo.Context) *jwt.Token {
	token := c.QueryParam("token")
	if token == "" {
		token = c.FormValue("token")
	}
	if token == "" {
		token = c.Request().Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
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

func GetUserByJwt(jwt *jwt.Token) model.User {
	if jwt == nil {
		return model.User{}
	}

	claims := jwtCustomClaims{}
	tmp, _ := json.Marshal(jwt.Claims)
	_ = json.Unmarshal(tmp, &claims)

	user := model.FindUserByID(claims.UserID)

	return user
}

func CheckIsAdminReq(c echo.Context) bool {
	jwt := GetJwtInstanceByReq(c)
	if jwt == nil {
		return false
	}

	user := GetUserByJwt(jwt)
	return user.IsAdmin
}

func GetUserByReq(c echo.Context) model.User {
	jwt := GetJwtInstanceByReq(c)
	user := GetUserByJwt(jwt)

	return user
}
