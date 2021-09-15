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
	"gorm.io/gorm/clause"
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
	params := make(map[string]interface{})

	refVal := reflect.ValueOf(paramsStruct)
	for i := 0; i < refVal.Type().NumField(); i++ {
		field := refVal.Type().Field(i)
		//fieldName := field.Name
		paramName := field.Tag.Get("mapstructure")
		paramTagP := field.Tag.Get("param")
		paramMethod := strings.ToUpper(field.Tag.Get("method"))

		// get param value
		paramVal := func() string {
			if paramMethod == "" {
				if c.Request().Method == "GET" {
					return c.QueryParam(paramName)
				} else if c.Request().Method == "POST" {
					return c.FormValue(paramName)
				}
			}

			if paramMethod == "GET" {
				return c.QueryParam(paramName)
			} else if paramMethod == "POST" {
				return c.FormValue(paramName)
			}
			return ""
		}()

		// check required param
		if paramName != "" && paramTagP == "required" {
			if strings.TrimSpace(paramVal) == "" {
				return false, RespError(c, "Param `"+paramName+"` is required.")
			}
		}

		// convert type
		if field.Type.Kind() == reflect.String {
			params[paramName] = paramVal
		} else if field.Type.Kind() == reflect.Int {
			u64, _ := strconv.ParseInt(paramVal, 10, 32)
			params[paramName] = int(u64)
		} else if field.Type.Kind() == reflect.Uint {
			u64, _ := strconv.ParseUint(paramVal, 10, 32)
			params[paramName] = uint(u64)
		}
		// } else if field.Type.Kind() == reflect.Array {
		// 	params[paramName] = c.QueryParams()[paramName]
		// }
	}

	err := mapstructure.Decode(params, destParams)
	if err != nil {
		logrus.Error("Params decode error: ", err)
		return false, RespError(c, "Params decode error.")
	}
	return true, nil
}

func CheckIfAllowed(c echo.Context, name string, email string, page model.Page) (bool, error) {
	// 如果用户是管理员，或者当前页只能管理员评论
	if IsAdminUser(name, email) || page.AdminOnly {
		if !CheckIsAdminReq(c) {
			return false, RespError(c, "需要验证管理员身份", Map{"need_login": true})
		}
	}

	return true, nil
}

func FindComment(id uint) model.Comment {
	var comment model.Comment
	lib.DB.Preload(clause.Associations).First(&comment, id)
	return comment
}

// 查找用户（返回：精确查找 AND）
func FindUser(name string, email string) model.User {
	var user model.User // 注：user 查找是 AND
	lib.DB.Where("name = ? AND email = ?", name, email).First(&user)
	return user
}

func IsAdminUser(name string, email string) bool {
	var user model.User // 还是用 AND 吧，OR 太混乱了
	lib.DB.Where("(name = ? AND email = ?) AND type = ?", name, email, model.UserAdmin).First(&user)
	return !user.IsEmpty()
}

func UpdateComment(comment *model.Comment) error {
	err := lib.DB.Save(comment).Error
	if err != nil {
		logrus.Error("Update Comment error: ", err)
	}
	return err
}

func FindCreatePage(pageKey string) model.Page {
	page := FindPage(pageKey)
	if page.IsEmpty() {
		page = NewPage(pageKey)
	}
	return page
}

func FindCreateUser(name string, email string) model.User {
	user := FindUser(name, email)
	if user.IsEmpty() {
		user = NewUser(name, email) // save a new user
	}
	return user
}

func NewUser(name string, email string) model.User {
	user := model.User{
		Name:  name,
		Email: email,
	}

	err := lib.DB.Create(&user).Error
	if err != nil {
		logrus.Error("Save User error: ", err)
	}

	return user
}

func UpdateUser(user *model.User) error {
	err := lib.DB.Save(user).Error
	if err != nil {
		logrus.Error("Update User error: ", err)
	}

	return err
}

func FindPage(key string) model.Page {
	var page model.Page
	lib.DB.Where(&model.Page{Key: key}).First(&page)
	return page
}

func FindPageByID(id uint) model.Page {
	var page model.Page
	lib.DB.Where("id = ?", id).First(&page)
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

func UpdatePage(page *model.Page) error {
	err := lib.DB.Save(page).Error
	if err != nil {
		logrus.Error("Update Page error: ", err)
	}
	return err
}

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ValidateURL(urlStr string) bool {
	_, err := url.ParseRequestURI(urlStr)
	return err == nil
}
