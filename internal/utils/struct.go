package utils

import (
	"encoding/json"

	"github.com/jeremywohl/flatten"
	"github.com/vmihailenco/msgpack"
)

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

func CopyStruct(src *map[string]interface{}, dest *map[string]interface{}) error {
	b, err := msgpack.Marshal(src)
	if err != nil {
		return err
	}
	err = msgpack.Unmarshal(b, dest)
	if err != nil {
		return err
	}
	return nil
}
