package json

import (
	"github.com/json-iterator/go"
)

// ParseJson parsing json strings
func ParseJson(data string, result interface{}) error {
	return ParseJsonFromBytes([]byte(data), result)
}

// StringifyJson json to string
func StringifyJson(data interface{}) string {
	return string(StringifyJsonToBytes(data))
}

// ParseJsonFromBytes parsing json bytes
func ParseJsonFromBytes(data []byte, result interface{}) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Unmarshal(data, result)
}

// json bytes to string
func StringifyJsonToBytes(data interface{}) []byte {
	b, _ := StringifyJsonToBytesWithErr(data)
	return b
}

func StringifyJsonToBytesWithErr(data interface{}) ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	b, err := json.Marshal(&data)
	return b, err
}
