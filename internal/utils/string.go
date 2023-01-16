package utils

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"
)

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
//
//	(bool) true => (string) "true"
//	(int) 0 => (string) "0"
func ToString(val interface{}) string {
	return fmt.Sprintf("%v", val)
}

// 将 JSON "数组中的"对象的 Values 全部转成 String 类型
// @note Array style is not the same as JSON Array, it uses the ToString() function.
//
//	[{"a":233}, {"b":true}, {"c":"233"}]
//	=> [{"a":"233"}, {"b":"true"}, {"c":"233"}]
//
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
