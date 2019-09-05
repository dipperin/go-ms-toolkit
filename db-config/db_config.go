package db_config

import (
	"fmt"
	"github.com/dipperin/go-ms-toolkit/qyenv"
	"io/ioutil"
	"strings"
)

type DbConfig struct {
	Username     string
	Password     string
	Host         string
	Port         string
	DbName       string
	MaxIdleConns int
	MaxOpenConns int
	DbCharset    string
}

// 只能通过这种方式获取配置对象
func NewDbConfig() *DbConfig {
	pwd := "1234"
	data, err := ioutil.ReadFile("/usr/local/.db/mysql.pas")
	if err != nil {
		fmt.Println("读取mysql密码文件出错:" + err.Error())
	} else {
		pwd = string(data)
		pwd = strings.TrimSpace(pwd)
	}
	appConf := GetAppConfig()
	conf := &DbConfig{
		Username:     appConf.MysqlUname,
		Password:     pwd,
		Host:         appConf.MysqlHost,
		Port:         appConf.MysqlPort,
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		DbCharset:    "utf8",
	}

	name, err := ioutil.ReadFile("/usr/local/.db/mysql.uname")
	if err == nil && len(name) > 0 {
		conf.Username = strings.TrimSpace(string(name))
	}

	return conf
}

var appConfig *AppConfig

func GetAppConfig() *AppConfig {
	if appConfig != nil {
		return appConfig
	}
	useDocker := qyenv.GetUseDocker()
	if useDocker == 1 {
		fmt.Println("采用的是dev容器的配置")
		appConfig = GetDevDockerConf()
	} else if useDocker == 2 {
		fmt.Println("采用的是prod生产的配置")
		appConfig = GetProdDockerConf()
	} else {
		fmt.Println("采用的是非docker环境的配置")
		appConfig = GetAppDefaultConf()
	}
	return appConfig
}

// 应用的配置
type AppConfig struct {
	RedisUrl   string
	MysqlHost  string
	MysqlPort  string
	MysqlUname string
	// http服务的端口（在docker中都是3000）
	HttpServerPort string
	MongoHost      string
	MongoUname     string
}

// 获取正常环境下的配置
func GetAppDefaultConf() *AppConfig {
	return &AppConfig{
		RedisUrl:   "127.0.0.1:6379",
		MysqlHost:  "localhost",
		MysqlPort:  "3306",
		MysqlUname: "root",
		MongoHost:  "localhost:27017",
		MongoUname: "admin",
	}
}

// 获取容器中的配置
func GetDevDockerConf() *AppConfig {
	return &AppConfig{
		RedisUrl:       "redis-master:6379",
		MysqlHost:      "qy-mysql",
		MysqlPort:      "3306",
		MysqlUname:     "qy",
		HttpServerPort: "3000",
		MongoHost:      "",
		MongoUname:     "",
	}
}

// 获取容器中的配置
func GetProdDockerConf() *AppConfig {
	uname := "qy"
	if data, err := ioutil.ReadFile("/usr/local/.db/mysql.uname"); err != nil {
		fmt.Println("读取mysql用户名文件出错:" + err.Error() + "。使用默认用户名。")
	} else {
		uname = strings.TrimSpace(string(data))
	}
	return &AppConfig{
		RedisUrl:       "redis-master:6379",
		MysqlHost:      "qy-mysql",
		MysqlPort:      "3306",
		MysqlUname:     uname,
		HttpServerPort: "3000",
		MongoHost:      "",
		MongoUname:     "",
	}
}
