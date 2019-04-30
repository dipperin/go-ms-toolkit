package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToCamelCase(t *testing.T) {
	assert.Equal(t, "", ToCamelCase(""))
	assert.Equal(t, "ABC", ToCamelCase("_a_b_c"))
	assert.Equal(t, "Play", ToCamelCase("play"))
	assert.Equal(t, "Play", ToCamelCase("Play"))
}
