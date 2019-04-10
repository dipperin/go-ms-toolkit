package db_config

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAppConfig(t *testing.T) {
	assert.NotNil(t, GetAppConfig())

	appConfig = nil

	flag.String("docker_env", "0", "程序运行环境配置：0非docker，1docker非生产，2docker生产")
	err := flag.Set("docker_env", "1")
	assert.NoError(t, err)
	assert.NotNil(t, GetAppConfig())

	appConfig = nil

	err1 := flag.Set("docker_env", "2")
	assert.NoError(t, err1)
	assert.NotNil(t, GetAppConfig())

	assert.NotNil(t, GetAppConfig())
}

func TestGetAppDefaultConf(t *testing.T) {
	assert.NotNil(t, GetAppDefaultConf())
}

func TestGetDevDockerConf(t *testing.T) {
	assert.NotNil(t, GetDevDockerConf())
}

func TestGetProdDockerConf(t *testing.T) {
	assert.NotNil(t, GetProdDockerConf())
}

func TestNewDbConfig(t *testing.T) {
	cfg := NewDbConfig()
	println(cfg.Username)
	println(cfg.Password)
}
