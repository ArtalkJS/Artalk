package http

import (
	"net/mail"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

func LoginGetUserToken(user model.User) string {
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

func ParamsDecode(c echo.Context, paramsStruct interface{}, destParams interface{}) (isContinue bool, resp error) {
	refVal := reflect.ValueOf(paramsStruct)
	for i := 0; i < refVal.Type().NumField(); i++ {
		field := refVal.Type().Field(i)
		//fieldName := field.Name
		paramTagM := field.Tag.Get("mapstructure")
		paramTagP := field.Tag.Get("param")
		//fmt.Println(field, paramTagM, paramTagP)

		if paramTagM != "" && paramTagP == "required" {
			if strings.TrimSpace(c.QueryParam(paramTagM)) == "" {
				return false, RespError(c, "Param `"+paramTagM+"` is required.")
			}
		}
	}

	// get the first
	params := make(map[string]interface{})
	for k, p := range c.QueryParams() {
		params[k] = p[0]
	}

	// convet type
	for i := 0; i < refVal.Type().NumField(); i++ {
		field := refVal.Type().Field(i)
		paramName := field.Tag.Get("mapstructure")

		if field.Type.Kind() == reflect.Int {
			u64, _ := strconv.ParseInt(c.QueryParam(paramName), 10, 32)
			params[paramName] = int(u64)
		}

		if field.Type.Kind() == reflect.Uint {
			u64, _ := strconv.ParseUint(c.QueryParam(paramName), 10, 32)
			params[paramName] = uint(u64)
		}

		if field.Type.Kind() == reflect.Array {
			params[paramName] = c.QueryParams()[paramName]
		}
	}

	err := mapstructure.Decode(params, destParams)
	if err != nil {
		logrus.Error("Params decode error: ", err)
		return false, RespError(c, "Params decode error.")
	}
	return true, nil
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
	return true, nil
}

func FindComment(id uint) model.Comment {
	var comment model.Comment
	lib.DB.First(&comment, id)
	return comment
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

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ValidateURL(urlStr string) bool {
	_, err := url.ParseRequestURI(urlStr)
	return err == nil
}
