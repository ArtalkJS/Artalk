package common

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/gofiber/fiber/v2"
)

type Map = map[string]interface{}

func ParamsDecode(c *fiber.Ctx, destParams interface{}) (isContinue bool, resp error) {
	reqMethod := c.Method()

	// validate required
	if isContinue, resp := ValidateRequired(c, destParams); !isContinue {
		return false, resp
	}

	// parse query
	if err := c.QueryParser(destParams); err != nil {
		return false, err
	}
	// parse path params
	if err := c.ParamsParser(destParams); err != nil {
		return false, err
	}
	// parse body
	if reqMethod == "POST" || reqMethod == "PUT" {
		if err := c.BodyParser(destParams); err != nil {
			return false, err
		}
	}

	return true, nil
}

func ValidateRequired(c *fiber.Ctx, destParams interface{}) (isContinue bool, resp error) {
	reqMethod := c.Method()
	reqPath := c.Path()

	refVal := reflect.ValueOf(destParams)
	for i := 0; i < refVal.Elem().Type().NumField(); i++ {
		k := refVal.Elem().Type().Field(i)

		validateTag := k.Tag.Get("validate")
		requiredField := (validateTag == "required")

		if !requiredField {
			continue // only check required field
		}

		// get param key
		paramKey := ""
		tagNames := []string{"query", "json", "form"}
		for _, tagName := range tagNames {
			paramKey = k.Tag.Get(tagName)
			if paramKey != "" {
				break
			}
		}
		if paramKey == "" {
			log.Errorf("[Validator] %s field '%s' real param key not found", reqPath, k.Name)
			continue
		}

		// get param value
		paramVal := func() string {
			switch {
			case reqMethod == "GET" || reqMethod == "DELETE":
				return c.Query(paramKey)
			case reqMethod == "POST" || reqMethod == "PUT":
				var destMap map[string]interface{}
				err := c.BodyParser(&destMap)
				if err != nil {
					log.Error("[Validator] "+reqPath+" field '"+k.Name+"' ", err)
					return ""
				}
				if v, ok := destMap[paramKey]; ok {
					return fmt.Sprintf("%v", v)
				}
				return ""
			}
			return ""
		}()

		// check required param
		if strings.TrimSpace(paramVal) == "" {
			return false, RespError(c, 400, "Param `"+paramKey+"` is required")
		}
	}

	return true, nil
}
