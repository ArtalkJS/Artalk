package artransfer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ArtalkJS/Artalk/internal/utils"
)

func checkIfJsonArr(str string) bool {
	x := bytes.TrimSpace([]byte(str))
	return len(x) > 0 && x[0] == '['
}

func checkIfJsonObj(str string) bool {
	x := bytes.TrimSpace([]byte(str))
	return len(x) > 0 && x[0] == '{'
}

func tryConvertLineJsonToArr(str string) (string, error) {
	// 尝试将一行一行的 Obj 转成 Arr
	arrTmp := []map[string]interface{}{}
	for _, line := range strings.Split(strings.TrimSpace(str), "\n") {
		var tmp map[string]interface{}
		err := json.Unmarshal([]byte(line), &tmp)
		if err != nil {
			return "", err
		}
		arrTmp = append(arrTmp, tmp)
	}
	r, err := json.Marshal(arrTmp)
	if err != nil {
		return "", err
	}
	return string(r), nil
}

// Json Decode (FAS: Fields All String Type)
// 解析 json 为字段全部是 string 类型的 struct
func jsonDecodeFAS(str string, fasStructure interface{}) error {
	if !checkIfJsonArr(str) {
		var err error
		str, err = tryConvertLineJsonToArr(str)
		if err != nil {
			return fmt.Errorf("JSON of array type is required: %w", err)
		}
	}

	err := json.Unmarshal([]byte(utils.JsonObjInArrAnyStr(str)), fasStructure) // lib.ToString()
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	return nil
}
