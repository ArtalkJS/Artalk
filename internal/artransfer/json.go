package artransfer

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

func isJsonArray(str string) bool {
	x := bytes.TrimSpace([]byte(str))
	return len(x) > 0 && x[0] == '['
}

// Convert JSON "Array of" Object's Values to String type
//
// Example:
//
//	[{"a":233}, {"b":true}, {"c":"233"}]
//	 => [{"a":"233"}, {"b":"true"}, {"c":"233"}]
func convertJSONValuesToString(jsonStr string) string {
	dest := []map[string]string{}
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

// Json Decode (FAS: Fields All String Type)
// Parse json to struct with all fields are string type
func jsonDecodeFAS(str string, fasStructure any) error {
	if !isJsonArray(str) {
		return fmt.Errorf("JSON of array type is required")
	}

	err := json.Unmarshal([]byte(convertJSONValuesToString(str)), fasStructure)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	return nil
}
