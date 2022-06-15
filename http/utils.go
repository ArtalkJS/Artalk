package http

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

type Map = map[string]interface{}

func ParamsDecode(c echo.Context, destParams interface{}) (isContinue bool, resp error) {
	params := make(map[string]interface{})

	refVal := reflect.ValueOf(destParams)
	for i := 0; i < refVal.Elem().Type().NumField(); i++ {
		field := refVal.Elem().Type().Field(i)
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
		requiredField := paramName != "" && paramTagP == "required"
		if requiredField {
			if strings.TrimSpace(paramVal) == "" {
				return false, RespError(c, "Param `"+paramName+"` is required")
			}
		}

		typeString := field.Type.Kind() == reflect.String
		typeInt := field.Type.Kind() == reflect.Int
		typeUint := field.Type.Kind() == reflect.Uint
		typeBool := field.Type.Kind() == reflect.Bool

		// convert type
		if typeString {
			params[paramName] = paramVal
		} else if typeBool {
			params[paramName] = (paramVal == "1" || paramVal == "true")
		} else if typeInt || typeUint {
			u64, err := strconv.ParseInt(paramVal, 10, 32)
			if requiredField && (err != nil || u64 == 0) {
				return false, RespError(c, "Param `"+paramName+"` is required")
			}

			if typeUint {
				params[paramName] = uint(u64)
			} else {
				params[paramName] = int(u64)
			}
		}
		// } else if field.Type.Kind() == reflect.Array {
		// 	params[paramName] = c.QueryParams()[paramName]
		// }
	}

	err := mapstructure.Decode(params, destParams)
	if err != nil {
		logrus.Error("Params decode error: ", err)
		return false, RespError(c, "Params decode error")
	}
	return true, nil
}

func CheckIsAllowed(c echo.Context, name string, email string, page model.Page, siteName string) (bool, error) {
	isAdminUser := model.IsAdminUserByNameEmail(name, email)

	// 如果用户是管理员，或者当前页只能管理员评论
	if isAdminUser || page.AdminOnly {
		if !CheckIsAdminReq(c) {
			return false, RespError(c, "需要验证管理员身份", Map{"need_login": true})
		}
	}

	return true, nil
}
