package consul_discovery

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"net"
	"os"
	"time"
)

const (
	GRPC = "GRPC"
	HTTP = "HTTP"
)

// NewConsulRegister create a new consul register
func NewConsulRegister(serverType, serverName, runVersion string) *ConsulRegister {
	return &ConsulRegister{
		Address:                        "127.0.0.1:8500",
		ServerName:                     serverName,
		Tag:                            []string{},
		Port:                           3000,
		DeregisterCriticalServiceAfter: time.Duration(1) * time.Minute,
		Interval:                       time.Duration(10) * time.Second,
		ip:                             localIP(),
		serverType:                     serverType,
		K8sHostName:                    os.Getenv("HOSTNAME"),
		ProductName:                    os.Getenv("product"),
		RunVersion:                     runVersion,
	}
}

// ConsulRegister consul service register
type ConsulRegister struct {
	serverType                     string
	Address                        string
	ServerName                     string
	Tag                            []string
	Port                           int
	DeregisterCriticalServiceAfter time.Duration
	Interval                       time.Duration
	ip                             string
	K8sHostName                    string
	ProductName                    string
	RunVersion                     string
}

// Register register service
func (r *ConsulRegister) Register() error {
	config := api.DefaultConfig()
	config.Address = r.Address
	client, err := api.NewClient(config)

	if err != nil {
		return err
	}

	agent := client.Agent()

	if err := agent.ServiceRegister(r.newReg()); err != nil {
		return err
	}

	return nil
}

func (r *ConsulRegister) genServerID() string {
	if r.K8sHostName != "" {
		return fmt.Sprintf("%v-%v-%v", r.K8sHostName, r.ip, r.Port)
	}

	return fmt.Sprintf("%v-%v-%v", r.ServerName, r.ip, r.Port)
}

func (r *ConsulRegister) newReg() *api.AgentServiceRegistration {
	switch r.serverType {
	case GRPC:
		return r.newGrpcRegistration()
	case HTTP:
		return r.newHttpRegistration()
	}

	panic("only support grpc or http")
}

func (r *ConsulRegister) newHttpRegistration() *api.AgentServiceRegistration {

	// TODO impl http
	return &api.AgentServiceRegistration{}
}

func (r *ConsulRegister) newGrpcRegistration() *api.AgentServiceRegistration {

	return &api.AgentServiceRegistration{
		ID:      r.genServerID(),                                                       // 服务节点的名称
		Name:    fmt.Sprintf("qy.%v.%v.%v", r.ProductName, r.RunVersion, r.ServerName), // 服务名称
		Tags:    r.Tag,                                                                 // tag，可以为空
		Port:    r.Port,                                                                // 服务端口
		Address: r.ip,                                                                  // 服务 IP
		Check: &api.AgentServiceCheck{ // 健康检查
			Interval:                       r.Interval.String(),                                 // 健康检查间隔
			GRPC:                           fmt.Sprintf("%v:%v/%v", r.ip, r.Port, r.ServerName), // grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
			DeregisterCriticalServiceAfter: r.DeregisterCriticalServiceAfter.String(),           // 注销时间，相当于过期时间
		},
	}
}

func localIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
