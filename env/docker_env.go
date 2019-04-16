package env

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
