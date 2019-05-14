package qyenv

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUseDocker(t *testing.T) {
	assert.Equal(t, 0, GetUseDocker())

	flag.String("docker_env", "0", "程序运行环境配置：0非docker，1docker非生产，2docker生产")
	err := flag.Set("docker_env", "1")
	assert.NoError(t, err)
	assert.Equal(t, 1, GetUseDocker())

	err1 := flag.Set("docker_env", "2")
	assert.NoError(t, err1)

	assert.Equal(t, 2, GetUseDocker())
}
