package env_provider

import (
	"strconv"

	"github.com/knadh/koanf/maps"
)

// Unflatten unflattens a map[string]interface{} to a nested map.
func Unflatten(mp map[string]interface{}) map[string]interface{} {
	// The unflatten function in koanf@v1.5.0/maps is only unflatten
	// to map[string]interface{} type without handling array.
	mp = maps.Unflatten(mp, ".")

	// The fixUnflattenResult function is used to fix the unflatten result,
	// convert the map to array if all keys are number.
	result, ok := NumberKeysMapToSlice(mp).(map[string]interface{})
	if !ok {
		return map[string]interface{}{}
	}
	return result
}

// NumberKeysMapToSlice converts a map with all keys as numbers to a slice recursively.
//
// Example:
//
//	{ "a": { "0": 1, "1": 2 } } => { "a": [1, 2] }
func NumberKeysMapToSlice(input interface{}) interface{} {
	allNumberKeys := func(m map[string]interface{}) (ok bool, maxNum int) {
		maxNum = -1
		for k := range m {
			if n, err := strconv.Atoi(k); err == nil {
				maxNum = max(maxNum, n)
			} else {
				return false, maxNum
			}
		}
		return true, maxNum
	}
	convertToSlice := func(m map[string]interface{}, size int) []interface{} {
		arr := make([]interface{}, size)
		for k, v := range m {
			i, _ := strconv.Atoi(k)
			arr[i] = NumberKeysMapToSlice(v) // Recursive call to fix nested map
		}
		return arr
	}
	handleSlice := func(arr []interface{}) []interface{} {
		for i, v := range arr {
			arr[i] = NumberKeysMapToSlice(v) // Recursive call to fix nested map or slice
		}
		return arr
	}
	handleMap := func(m map[string]interface{}) interface{} {
		if ok, maxNum := allNumberKeys(m); ok {
			return handleSlice(convertToSlice(m, maxNum+1))
		}
		for k, v := range m {
			m[k] = NumberKeysMapToSlice(v) // Recursive call to fix nested map
		}
		return m
	}

	switch v := input.(type) {
	case map[string]interface{}:
		return handleMap(v)
	case []interface{}:
		return handleSlice(v)
	default:
		return input
	}
}
