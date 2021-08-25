package http

import (
	"time"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func LoginUser(user model.User) string {
	// Set custom claims
	claims := &jwtCustomClaims{
		UserName:  user.Name,
		UserEmail: user.Email,
		UserType:  user.Type,
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

func CheckIsAdmin(c echo.Context) bool {
	jwt := c.Get("user").(*jwt.Token)
	claims := jwt.Claims.(*jwtCustomClaims)
	name := claims.UserName
	email := claims.UserEmail

	if claims.UserType != model.UserAdmin {
		return false
	}

	// check user from database
	user := FindUser(name, email)
	if user.IsEmpty() {
		return false
	}

	return user.Type == model.UserAdmin
}

func CheckIfAllowed(c echo.Context, user model.User, page model.Page) (bool, error) {
	return false, nil
}

func FindUser(name string, email string) model.User {
	var user model.User
	lib.DB.Where(&model.User{Name: name, Email: email}).First(&user)
	return user
}

func NewUser(name string, email string, link string) model.User {
	user := model.User{
		Name:  name,
		Email: email,
		Link:  link,
	}

	err := lib.DB.Create(&user).Error
	if err != nil {
		logrus.Error("Save User error: ", err)
	}

	return user
}

func UpdateUser(user *model.User) {
	err := lib.DB.Save(user).Error
	if err != nil {
		logrus.Error("Update User error: ", err)
	}
}

func FindPage(key string) model.Page {
	var page model.Page
	lib.DB.Where(&model.Page{Key: key}).First(&page)
	return page
}

func NewPage(key string) model.Page {
	page := model.Page{
		Key: key,
	}

	err := lib.DB.Create(&page).Error
	if err != nil {
		logrus.Error("Save Page error: ", err)
	}

	return page
}

func UpdatePage(page *model.Page) {
	err := lib.DB.Save(page).Error
	if err != nil {
		logrus.Error("Update Page error: ", err)
	}
}
