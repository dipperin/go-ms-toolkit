package gjson_util

import (
	gjson_util "gitlab2018.com/arch/rule-engine-go/common/util/gjson-util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInt(t *testing.T) {
	_, err := gjson_util.GetInt(`{}`, "a")
	assert.Error(t, err)
	_, err = gjson_util.GetInt(`{"a": "sdf"}`, "a")
	assert.Error(t, err)
	_, err = gjson_util.GetInt(`{"a": "123"}`, "a")
	assert.Error(t, err)

	a, err := gjson_util.GetInt(`{"a": 123}`, "a")
	assert.NoError(t, err)
	assert.Equal(t, int64(123), a)
}

func TestGetString(t *testing.T) {
	_, err := gjson_util.GetString(`{}`, "a")
	assert.Error(t, err)
	_, err = gjson_util.GetString(`{"a":123}`, "a")
	assert.Error(t, err)

	a, err := gjson_util.GetString(`{"a":"2123"}`, "a")
	assert.NoError(t, err)
	assert.Equal(t, "2123", a)
}

func TestGetStringArr(t *testing.T) {
	_, err := gjson_util.GetStringArr(`{}`, "a")
	assert.Error(t, err)
	_, err = gjson_util.GetStringArr(`{"a":123}`, "a")
	assert.Error(t, err)
	_, err = gjson_util.GetStringArr(`{"a":["123",123]}`, "a")
	assert.Error(t, err)

	a, err := gjson_util.GetStringArr(`{"a":["123","da"]}`, "a")
	assert.NoError(t, err)
	assert.Equal(t, []string{"123", "da"}, a)
}

func TestGetFloat(t *testing.T) {
	_, err := gjson_util.GetFloat(`{}`, "a")
	assert.Error(t, err)
	_, err = gjson_util.GetFloat(`{"a": "sdf"}`, "a")
	assert.Error(t, err)
	_, err = gjson_util.GetFloat(`{"a": "123.3"}`, "a")
	assert.Error(t, err)

	a, err := gjson_util.GetFloat(`{"a": 123.123}`, "a")
	assert.NoError(t, err)
	assert.Equal(t, float64(123.123), a)
}
