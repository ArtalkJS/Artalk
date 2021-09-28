package http

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

func LoginGetUserToken(user model.User) string {
	// Set custom claims
	claims := &jwtCustomClaims{
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

func ParamsDecode(c echo.Context, paramsStruct interface{}, destParams interface{}) (isContinue bool, resp error) {
	params := make(map[string]interface{})

	refVal := reflect.ValueOf(paramsStruct)
	for i := 0; i < refVal.Type().NumField(); i++ {
		field := refVal.Type().Field(i)
		//fieldName := field.Name
		paramName := field.Tag.Get("mapstructure")
		paramTagP := field.Tag.Get("param")
		paramMethod := strings.ToUpper(field.Tag.Get("method"))

		if paramName == "" {
			continue
		}

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
		} else if field.Type.Kind() == reflect.Bool {
			params[paramName] = (paramVal == "1")
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

func CheckIfAllowed(c echo.Context, name string, email string, page model.Page, siteName string) (bool, error) {
	isAdminUser := model.IsAdminUser(name, email)

	// 如果用户是管理员，或者当前页只能管理员评论
	if isAdminUser || page.AdminOnly {
		if !CheckIsAdminReq(c) {
			return false, RespError(c, "需要验证管理员身份", Map{"need_login": true})
		}
	}

	return true, nil
}

func CheckReferer(c echo.Context, siteName string) (bool, error) {
	isAdminReq := CheckIsAdminReq(c)
	if isAdminReq || siteName == "" {
		return true, nil
	}

	// 请求 Referer 合法性判断
	site := model.FindSite(siteName)
	allowUrls := strings.Split(site.Url, ",")
	referer := c.Request().Referer()
	if referer == "" {
		return true, nil
	}

	pr, err := url.Parse(referer)
	if err != nil {
		return true, nil
	}

	allow := false
	for _, u := range allowUrls {
		u = strings.TrimSpace(u)
		if u == "" {
			continue
		}
		pu, err := url.Parse(u)
		if err != nil {
			continue
		}
		if pu.Hostname() == pr.Hostname() {
			allow = true
		}
	}
	if !allow {
		return false, RespError(c, "非法请求：Referer 不被允许")
	}

	return true, nil
}

func CheckSite(c echo.Context, siteName string, destID *uint) (bool, error) {
	if siteName == "" {
		return true, nil
	}

	site := model.FindSite(siteName)
	if site.IsEmpty() {
		return false, RespError(c, "Site 未找到")
	}

	// 检测 Referer 合法性
	if isOK, resp := CheckReferer(c, siteName); !isOK {
		return false, resp
	}

	*destID = site.ID

	return true, nil
}
