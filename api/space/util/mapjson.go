package util

import (
	"encoding/json"
	"errors"
)
var _ = errors.New("")
func MapCsvJson(field []string, data []interface {}) map[string] interface{} {
	mapping := make(map[string] interface {})
	for i := 0; i < len(field); i++ {
		mapping[field[i]] = data[i]
	}
	return mapping
}

func FindPositionArray(strFind string, listStr []string) int {
	for i, str := range listStr {
		if str == strFind {
			return i
		}
	}
	return -1
}

func JsonStringify(objson map[string] interface{}) (string, error) {
	str, err := json.Marshal(objson)
	return string(str), err
}

func JsonParseInterface(str string) (map[string] interface {}, error) {
	var dat map[string] interface{}
	byt := []byte(str) 
	if err := json.Unmarshal(byt, &dat); err != nil {
		return nil, err
	}
	return dat, nil
}