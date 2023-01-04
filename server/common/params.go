package common

import (
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Map = map[string]interface{}

func ParamsDecode(c *fiber.Ctx, destParams interface{}) (isContinue bool, resp error) {
	reqMethod := c.Method()

	var errParse error
	switch reqMethod {
	case "GET":
		errParse = c.QueryParser(destParams)
	case "POST":
		errParse = c.BodyParser(destParams)
	}
	if errParse != nil {
		return false, errParse
	}

	refVal := reflect.ValueOf(destParams)
	for i := 0; i < refVal.Elem().Type().NumField(); i++ {
		k := refVal.Elem().Type().Field(i)
		// v := refVal.Elem().Field(i)
		// fieldName := k.Name

		validateTag := k.Tag.Get("validate")
		requiredField := (validateTag == "required")

		if !requiredField {
			continue // 仅 required 检测时才继续执行之后的代码
		}

		paramName := k.Tag.Get("query")
		if paramName == "" {
			paramName = k.Tag.Get("form")
		}
		if paramName == "" {
			continue
		}

		// get param value
		paramVal := func() string {
			switch reqMethod {
			case "GET":
				return c.Query(paramName)
			case "POST":
				return c.FormValue(paramName)
			}
			return ""
		}()

		// check required param
		if requiredField && strings.TrimSpace(paramVal) == "" {
			return false, RespError(c, "Param `"+paramName+"` is required")
		}

		// 类型转换交给 fiber 内置 BodyParser 来做，这里不再实现
		// convert type
		// kind := k.Type.Kind()
		// if kind == reflect.String {
		// 	v.SetString(paramVal)
		// } else if kind == reflect.Bool {
		// 	v.SetBool((paramVal == "1" || paramVal == "true"))
		// } else if (kind == reflect.Int) || (kind == reflect.Uint) {
		// 	u64, err := strconv.ParseInt(paramVal, 10, 32)
		// 	if requiredField && (err != nil || u64 == 0) {
		// 		return false, RespError(c, "Param `"+paramName+"` is required")
		// 	}
		// 	if kind == reflect.Uint {
		// 		v.SetUint(uint64(u64))
		// 	} else {
		// 		v.SetInt(u64)
		// 	}
		// }
		// // } else if kind == reflect.Array {
		// // }
	}

	return true, nil
}
