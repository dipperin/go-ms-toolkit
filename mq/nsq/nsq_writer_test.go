package nsq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNsqWriter(t *testing.T) {
	t.Skip("ignored")
	writer := NewNsqWriter([]string{"127.0.0.1:4150", "127.0.0.1:4152"})
	err := writer.newProducers().PublishString("topic1", "something")
	assert.NoError(t, err)
	err = writer.Publish("topic2", []string{"sssss", "heiheihei"})
	assert.NoError(t, err)
}
