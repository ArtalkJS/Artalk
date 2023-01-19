package utils

import (
	"fmt"
	"regexp"
)

// 解析 Mustaches 语法
// 替换 {{ key }} 为 val
func RenderMustaches(data string, dict map[string]interface{}, valueGetter ...func(k string, v interface{}) string) string {
	r := regexp.MustCompile(`{{\s*(.*?)\s*}}`)

	return r.ReplaceAllStringFunc(data, func(m string) string {
		key := r.FindStringSubmatch(m)[1]
		if val, isExist := dict[key]; isExist {
			if len(valueGetter) > 0 {
				return valueGetter[0](key, val)
			} else {
				return fmt.Sprintf("%v", val)
			}
		}

		return m
	})
}
