package json

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testModel struct {
	Text string `json:"text"`
}

func TestParseJson(t *testing.T) {
	jsonStr := `{"text":"hahaha"}`

	var model testModel
	assert.NoError(t, ParseJson(jsonStr, &model))
	assert.Equal(t, "hahaha", model.Text)
}

func TestStringifyJson(t *testing.T) {
	var model = testModel{Text: "123"}

	str := StringifyJson(model)
	assert.Equal(t, `{"text":"123"}`, str)
}
