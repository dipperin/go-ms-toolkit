package qyenv

import (
	"flag"
	"os"
)

func GetUseDocker() int {
	// 也可以在环境变量中设置
	dEnv := os.Getenv("docker_env")
	if dEnv != "" {
		switch dEnv {
		case "1":
			return 1
		case "2":
			return 2
		default:
			return 0
		}
	}

	f := flag.Lookup("docker_env")
	if f == nil || f.Value.String() == "0" {
		// 非docker
		return 0
	} else if f.Value.String() == "2" {
		// 生产
		return 2
	} else {
		// 默认返回开发和测试的配置
		return 1
	}
}

// 0 非docker环境 1 docker中非生产环境 2 docker中生产环境
func GetDockerEnv() string {
	return os.Getenv("docker_env")
}

// db名称配置 dev test preprod prod
func GetDBEnv() string {
	return os.Getenv("db_env")
}

// 程序执行环境配置 dev test preprod prod
func GetRunEnv() string {
	return os.Getenv("run_env")
}

// determine if it is a unit test environment
func IsUnitTestEnv() bool {
	return flag.Lookup("test.v") != nil
}

// 获取产品名，用于：1. 区分数据库名称
func GetProductName() string {
	return os.Getenv("product")
}