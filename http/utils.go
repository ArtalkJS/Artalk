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

// 检测 Origin 合法性
// 防止 CSRF 攻击
// @see https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html
func CheckOrigin(c echo.Context, site model.Site) (bool, error) {
	isAdminReq := CheckIsAdminReq(c)
	if isAdminReq {
		return true, nil // 管理员直接允许
	}

	// 可信来源 URL
	allowSrcURLs := []string{}

	// 用户配置
	allowSrcURLs = append(allowSrcURLs, config.Instance.TrustedDomains...) // 允许配置文件域名
	allowSrcURLs = append(allowSrcURLs, site.ToCooked().Urls...)           // 允许数据库站点 URLs 中的域名
	if len(allowSrcURLs) == 0 {
		return true, nil // 若用户配置列表中无数据，则取消控制
	}
	if lib.ContainsStr(allowSrcURLs, "*") {
		return true, nil // 列表中出现通配符关闭控制
	}

	host := c.Request().Host
	realHostUnderProxy := c.Request().Header.Get("X-Forwarded-Host")
	if realHostUnderProxy != "" {
		host = realHostUnderProxy
	}

	// 读取 Origin 数据
	// @note Origin 标头在前端 fetch POST 操作中总是携带的，
	// 		 即使配置 Referrer-Policy: no-referrer
	// @see https://stackoverflow.com/questions/42239643/when-do-browsers-send-the-origin-header-when-do-browsers-set-the-origin-to-null
	origin := c.Request().Header.Get(echo.HeaderOrigin)
	if origin == "" || origin == "null" {
		// 从 Referer 获取 Origin
		referer := c.Request().Referer()
		if referer == "" {
			return false, RespError(c, "无效请求，Origin 无法获取")
		}
		origin = referer
	}

	pOrigin, err := url.Parse(origin)
	if err != nil {
		return false, RespError(c, "Origin 不合法")
	}

	// 系统配置：默认允许来自相同域名的请求
	allowSrcURLs = append(allowSrcURLs, c.Scheme()+"://"+host)

	allowSrcURLs = lib.RemoveDuplicates(allowSrcURLs) // 去重
	for _, a := range allowSrcURLs {
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
		if pAllow.Host == pOrigin.Host {
			return true, nil
		}
	}

	return false, RespError(c, "非法请求，请检查可信域名配置")
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

	// 检测 Origin 合法性
	if isOK, resp := CheckOrigin(c, site); !isOK {
		return false, resp
	}

	*destID = site.ID // 更新源 id

	return true, nil
}
