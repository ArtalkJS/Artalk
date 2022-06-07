package http

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
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

func CheckReferer(c echo.Context, site model.Site) (bool, error) {
	isAdminReq := CheckIsAdminReq(c)
	if isAdminReq {
		return true, nil // 管理员直接允许
	}

	// 读取 Referer
	referer := c.Request().Referer()
	if referer == "" {
		return false, RespError(c, "需携带 Referer 访问，请检查前端 Referrer-Policy 设置")
	}

	pReferer, err := url.Parse(referer)
	if err != nil {
		return false, RespError(c, "Referer 不合法")
	}

	// 可信来源 URL
	allowReferrers := []string{}

	// 用户配置
	allowReferrers = append(allowReferrers, config.Instance.TrustedDomains...) // 允许配置文件域名
	allowReferrers = append(allowReferrers, site.ToCooked().Urls...)           // 允许数据库站点 URLs 中的域名
	if len(allowReferrers) == 0 {
		return true, nil // 若用户配置列表中无数据，则取消控制
	}
	if lib.ContainsStr(allowReferrers, "*") {
		return true, nil // 列表中出现通配符关闭 Referer 控制
	}

	// 系统配置：默认允许来自相同域名的请求
	allowReferrers = append(allowReferrers, c.Scheme()+"://"+c.Request().Host)

	allowReferrers = lib.RemoveDuplicates(allowReferrers) // 去重
	for _, a := range allowReferrers {
		a = strings.TrimSpace(a)
		if a == "" {
			continue
		}
		pAllow, err := url.Parse(a)
		if err != nil {
			continue
		}

		// 在可信来源列表中匹配 Referer 的 host 部分 (含端口) 则放行
		// @see https://web.dev/referrer-best-practices/
		// Referrer-Policy 不能设为 no-referer，
		// Chrome v85+ 默认为：strict-origin-when-cross-origin。
		// 前端页面 head 不配置 <meta name="referrer" content="no-referer" />，
		// 浏览器默认都会至少携带 Origin 数据 (不带 path，但包含端口)
		if pAllow.Host == pReferer.Host {
			return true, nil
		}
	}

	return false, RespError(c, "不允许的 Referer，请将其加入可信域名")
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
