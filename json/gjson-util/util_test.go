package gjson_util

import (
	"github.com/dipperin/go-ms-toolkit/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInt(t *testing.T) {
	_, err := json.GetInt(`{}`, "a")
	assert.Error(t, err)
	_, err = json.GetInt(`{"a": "sdf"}`, "a")
	assert.Error(t, err)
	_, err = json.GetInt(`{"a": "123"}`, "a")
	assert.Error(t, err)

	a, err := json.GetInt(`{"a": 123}`, "a")
	assert.NoError(t, err)
	assert.Equal(t, int64(123), a)
}

func TestGetString(t *testing.T) {
	_, err := json.GetString(`{}`, "a")
	assert.Error(t, err)
	_, err = json.GetString(`{"a":123}`, "a")
	assert.Error(t, err)

	a, err := json.GetString(`{"a":"2123"}`, "a")
	assert.NoError(t, err)
	assert.Equal(t, "2123", a)
}

func TestGetStringArr(t *testing.T) {
	_, err := json.GetStringArr(`{}`, "a")
	assert.Error(t, err)
	_, err = json.GetStringArr(`{"a":123}`, "a")
	assert.Error(t, err)
	_, err = json.GetStringArr(`{"a":["123",123]}`, "a")
	assert.Error(t, err)

	a, err := json.GetStringArr(`{"a":["123","da"]}`, "a")
	assert.NoError(t, err)
	assert.Equal(t, []string{"123", "da"}, a)
}

func TestGetFloat(t *testing.T) {
	_, err := json.GetFloat(`{}`, "a")
	assert.Error(t, err)
	_, err = json.GetFloat(`{"a": "sdf"}`, "a")
	assert.Error(t, err)
	_, err = json.GetFloat(`{"a": "123.3"}`, "a")
	assert.Error(t, err)

	a, err := json.GetFloat(`{"a": 123.123}`, "a")
	assert.NoError(t, err)
	assert.Equal(t, float64(123.123), a)
}
