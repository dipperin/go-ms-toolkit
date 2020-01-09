package gjson_util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInt(t *testing.T) {
	_, err := GetInt(`{}`, "a")
	assert.Error(t, err)
	_, err = GetInt(`{"a": "sdf"}`, "a")
	assert.Error(t, err)
	_, err = GetInt(`{"a": "123"}`, "a")
	assert.Error(t, err)

	a, err := GetInt(`{"a": 123}`, "a")
	assert.NoError(t, err)
	assert.Equal(t, int64(123), a)
}

func TestGetString(t *testing.T) {
	_, err := GetString(`{}`, "a")
	assert.Error(t, err)
	_, err = GetString(`{"a":123}`, "a")
	assert.Error(t, err)

	a, err := GetString(`{"a":"2123"}`, "a")
	assert.NoError(t, err)
	assert.Equal(t, "2123", a)
}

func TestGetStringArr(t *testing.T) {
	_, err := GetStringArr(`{}`, "a")
	assert.Error(t, err)
	_, err = GetStringArr(`{"a":123}`, "a")
	assert.Error(t, err)
	_, err = GetStringArr(`{"a":["123",123]}`, "a")
	assert.Error(t, err)

	a, err := GetStringArr(`{"a":["123","da"]}`, "a")
	assert.NoError(t, err)
	assert.Equal(t, []string{"123", "da"}, a)
}

func TestGetFloat(t *testing.T) {
	_, err := GetFloat(`{}`, "a")
	assert.Error(t, err)
	_, err = GetFloat(`{"a": "sdf"}`, "a")
	assert.Error(t, err)
	_, err = GetFloat(`{"a": "123.3"}`, "a")
	assert.Error(t, err)

	a, err := GetFloat(`{"a": 123.123}`, "a")
	assert.NoError(t, err)
	assert.Equal(t, float64(123.123), a)
}

func TestGetArray(t *testing.T) {
	_, err := GetArray(`{}`, "a")
	assert.Error(t, err)
	_, err = GetArray(`{"a":123}`, "a")
	assert.Error(t, err)
	_, err = GetArray(`{"a":["123","321"]}`, "a")
	assert.NoError(t, err)
	_, err = GetArray(`{"grade":[]}`, "scores")
	assert.Error(t, err)
}
