package http

import (
	"fmt"
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

func CheckIfAllowed(c echo.Context, name string, email string, page model.Page, siteName string) (bool, error) {
	isAdminUser := model.IsAdminUserByNameEmail(name, email)

	// 如果用户是管理员，或者当前页只能管理员评论
	if isAdminUser || page.AdminOnly {
		if !CheckIsAdminReq(c) {
			return false, RespError(c, "需要验证管理员身份", Map{"need_login": true})
		}
	}

	return true, nil
}

func CheckReferer(c echo.Context, site model.Site) (bool, error) {
	isAdminReq := CheckIsAdminReq(c)
	if isAdminReq || site.IsEmpty() {
		return true, nil
	}

	// 可信域名配置
	confTrustedDomains := config.Instance.TrustedDomains

	// 请求 Referer 合法性判断
	if strings.TrimSpace(site.Urls) == "" && len(confTrustedDomains) == 0 {
		return true, nil // 若 url 字段为空，则取消控制
	}

	// 可信域名出现通配符关闭 Referer 控制
	if lib.ContainsStr(confTrustedDomains, "*") {
		return true, nil
	}

	allowUrls := site.ToCooked().Urls
	if len(confTrustedDomains) != 0 {
		allowUrls = append(allowUrls, confTrustedDomains...)
	}

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
			break
		}
	}
	if !allow {
		return false, RespError(c, "非法请求：Referer 不被允许")
	}

	return true, nil
}

func CheckSite(c echo.Context, siteName *string, destID *uint, destSiteAll *bool) (bool, error) {
	// 启用源 SiteAll
	if destSiteAll != nil {
		// 传入站点名参数 == "__ATK_SITE_ALL" 时取消站点隔离
		if *siteName == lib.ATK_SITE_ALL {
			if !CheckIsAdminReq(c) {
				return false, RespError(c, "仅管理员查询允许取消站点隔离")
			}
			*destSiteAll = true
			return true, nil
		} else {
			*destSiteAll = false
		}
	}

	if *siteName == "" {
		// 传入值为空，使用默认 site
		siteDefault := strings.TrimSpace(config.Instance.SiteDefault)
		if siteDefault != "" {
			// 没有则创建
			model.FindCreateSite(siteDefault)
		}
		*siteName = siteDefault // 更新源 name

		return true, nil
	}

	site := model.FindSite(*siteName)
	if site.IsEmpty() {
		return false, RespError(c, fmt.Sprintf("未找到站点：`%s`，请控制台创建站点", *siteName), Map{
			"err_no_site": true,
		})
	}

	// 检测 Referer 合法性
	if isOK, resp := CheckReferer(c, site); !isOK {
		return false, resp
	}

	*destID = site.ID // 更新源 id

	return true, nil
}
