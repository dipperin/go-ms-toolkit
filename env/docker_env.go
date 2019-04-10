package env

import "flag"

func GetUseDocker() int {
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