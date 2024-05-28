package utils

import (
	"fmt"
	"net/url"
	"strings"
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
	split := strings.Split(s, sep)
	arr := []string{}
	for _, v := range split {
		arr = append(arr, strings.TrimSpace(v))
	}
	return RemoveBlankStrings(arr)
}

func RemoveBlankStrings(s []string) []string {
	r := []string{}
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

// 任何类型转 String
//
//	(bool) true => (string) "true"
//	(int) 0 => (string) "0"
func ToString(val interface{}) string {
	return fmt.Sprintf("%v", val)
}
