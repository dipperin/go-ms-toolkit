package gjson_util

import (
	"fmt"
	"github.com/tidwall/gjson"
)

var (
	ErrInvalidValueTypeInJson = "invalid value type in json"
)

// 只有一定为数字时才能获取到， "123"也是非法的
func GetInt(data, path string) (int64, error) {
	r := gjson.Get(data, path)
	if r.Type != gjson.Number {
		return 0, fmt.Errorf("%v, path: %v need: number", ErrInvalidValueTypeInJson, path)
	}
	return r.Int(), nil
}

func GetFloat(data, path string) (float64, error) {
	r := gjson.Get(data, path)
	if r.Type != gjson.Number {
		return 0, fmt.Errorf("%v, path: %v need: number", ErrInvalidValueTypeInJson, path)
	}
	return r.Float(), nil
}

func GetString(data, path string) (string, error) {
	r := gjson.Get(data, path)
	if r.Type != gjson.String {
		return "", fmt.Errorf("%v, path: %v need: string", ErrInvalidValueTypeInJson, path)
	}
	return r.String(), nil
}

func GetStringArr(data, path string) (result []string, err error) {
	r := gjson.Get(data, path)
	if r.Type != gjson.JSON || !r.IsArray() {
		return nil, fmt.Errorf("%v, path: %v need: []string", ErrInvalidValueTypeInJson, path)
	}
	for _, elem := range r.Array() {
		if elem.Type != gjson.String {
			return nil, fmt.Errorf("%v, path: %v need: string in []string", ErrInvalidValueTypeInJson, path)
		}
		result = append(result, elem.String())
	}
	return
}

func GetBool(data, path string) (bool, error) {
	r := gjson.Get(data, path)
	if r.Type == gjson.True {
		return true, nil
	} else if r.Type == gjson.False {
		return false, nil
	} else {
		return false, fmt.Errorf("%v, path: %v need: bool", ErrInvalidValueTypeInJson, path)
	}
}

func GetArray(data, path string) ([]gjson.Result, error) {
	r := gjson.Get(data, path)
	if !r.IsArray() {
		return nil, fmt.Errorf("%v, path: %v type: %d need: array", ErrInvalidValueTypeInJson, path, r.Type)
	}
	return r.Array(), nil
}
