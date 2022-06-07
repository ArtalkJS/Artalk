package lib

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/jeremywohl/flatten"
	"github.com/microcosm-cc/bluemonday"
	"github.com/tidwall/gjson"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	"github.com/asaskevich/govalidator"
)

// EnsureDir ensures that a target directory exists (like `mkdir -p`),
func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0700)
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func StructToMap(s interface{}) map[string]interface{} {
	b, _ := json.Marshal(s)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	return m
}

func StructToFlatDotMap(s interface{}) map[string]interface{} {
	m := StructToMap(s)
	mainFlat, err := flatten.Flatten(m, "", flatten.DotStyle)
	if err != nil {
		return map[string]interface{}{}
	}
	return mainFlat
}

func Marked(markdownStr string) (string, error) {
	bmPolicy := bluemonday.UGCPolicy()
	bmPolicy.RequireNoReferrerOnLinks(true)
	bmPolicy.AllowAttrs("width", "height", "align", "atk-emoticon").OnElements("img")
	bmPolicy.AllowAttrs("style", "class", "align").OnElements("span", "p", "div", "a")

	// https://github.com/yuin/goldmark#security
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(markdownStr), &buf); err != nil {
		return "", err
	}

	return bmPolicy.SanitizeReader(&buf).String(), nil
}

func AddQueryToURL(urlStr string, queryMap map[string]string) string {
	u, _ := url.Parse(urlStr)

	q, _ := url.ParseQuery(u.RawQuery)
	for k, v := range queryMap {
		q.Add(k, v)
	}

	u.RawQuery = q.Encode()
	return u.String()
}

// ContainsStr returns true if an str is present in a iteratee.
func ContainsStr(s []string, v string) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

// RemoveDuplicates removes the duplicates strings from a Slice
func RemoveDuplicates(arr []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range arr {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func SplitAndTrimSpace(s string, sep string) []string {
	splitted := strings.Split(s, sep)
	arr := []string{}
	for _, v := range splitted {
		arr = append(arr, strings.TrimSpace(v))
	}
	return RemoveBlankStrings(arr)
}

func RemoveBlankStrings(s []string) []string {
	var r []string
	for _, str := range s {
		if strings.TrimSpace(str) != "" {
			r = append(r, str)
		}
	}
	return r
}

func TruncateString(str string, length int) string {
	if length <= 0 {
		return ""
	}

	// This code cannot support Chinese
	// orgLen := len(str)
	// if orgLen <= length {
	//     return str
	// }
	// return str[:length]

	// Support Chinese
	truncated := ""
	count := 0
	for _, char := range str {
		truncated += string(char)
		count++
		if count >= length {
			break
		}
	}
	return truncated
}

//#region JSON Any To String (for Transfer)
//******************************************

// 任何类型转 String
//	(bool) true => (string) "true"
//	(int) 0 => (string) "0"
func ToString(val interface{}) string {
	return fmt.Sprintf("%v", val)
}

// 将 JSON "数组中的"对象的 Values 全部转成 String 类型
// @note Array style is not the same as JSON Array, it uses the ToString() function.
//	[{"a":233}, {"b":true}, {"c":"233"}]
//	=> [{"a":"233"}, {"b":"true"}, {"c":"233"}]
// @relevant ToString()
func JsonObjInArrAnyStr(jsonStr string) string {
	var dest []map[string]string
	for _, item := range gjson.Parse(jsonStr).Array() {
		dItem := map[string]string{}
		item.ForEach(func(key, value gjson.Result) bool {
			dItem[key.String()] = value.String()
			return true
		})
		dest = append(dest, dItem)
	}
	j, _ := json.Marshal(dest)
	return string(j)
}

//#endregion

//#region Validators
//******************************************

func ValidateEmail(email string) bool {
	return govalidator.IsEmail(email)
}

func ValidateURL(url string) bool {
	return govalidator.IsRequestURL(url)
}

//#endregion
