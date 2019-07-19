package consul_discovery

import (
	"fmt"
	"github.com/dipperin/go-ms-toolkit/log"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc/resolver"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	ConsulResolver = "consul"
)

func InitConsulGRPCResolver() {
	resolver.Register(newBuilder())
}

func newBuilder() resolver.Builder {
	return &consulBuilder{}
}

type consulBuilder struct {
}

func (b *consulBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	name, version := splitVersion(target.Endpoint, "@")
	r := newConsulResolver(target.Authority, name, version, cc)
	return r, nil
}

func (b *consulBuilder) Scheme() string {
	return ConsulResolver
}

// splitVersion returns the values from strings.SplitN(s, sep, 2).
// If sep is not found, it returns (s, "") instead.
func splitVersion(s, sep string) (string, string) {
	spl := strings.SplitN(s, sep, 2)
	if len(spl) < 2 {
		return s, ""
	}
	return spl[0], spl[1]
}

func newConsulResolver(address, serverName, runVersion string, cc resolver.ClientConn) *consulResolver {
	r := &consulResolver{
		cc:            cc,
		address:       address,
		serverName:    serverName,
		serverVersion: runVersion,
		lastIndex:     0,
		rn:            make(chan struct{}),
		stopC:         make(chan struct{}),
	}

	// 初始化consul client
	r.initConsulClient()

	r.wg.Add(1)
	go r.watcher()

	return r
}

type consulResolver struct {
	// client conn 这里主要是去回调client，更新stat（Resolved addresses for the target）
	cc resolver.ClientConn
	// consul的地址端口, 默认应该是127.0.0.1:8500, 在k8s里应该是consul-master-app:8500
	address string
	// 服务名
	serverName string
	// server version
	serverVersion string
	// consul同步点，这个调用将一直阻塞，直到有新的更新
	lastIndex uint64
	// consul client
	client *api.Client
	// rn channel is used by ResolveNow() to force an immediate resolution of the target.
	rn    chan struct{}
	stopC chan struct{}
	wg    sync.WaitGroup
}

func (r *consulResolver) ResolveNow(resolver.ResolveNowOption) {
	select {
	case r.rn <- struct{}{}:
	default:
	}
}

func (r *consulResolver) Close() {
	if r.stopC == nil {
		return
	}

	close(r.stopC)
	r.wg.Wait()
}

func (r *consulResolver) initConsulClient() {
	if r.client != nil {
		return
	}

	config := api.DefaultConfig()
	config.Address = r.address

	client, err := api.NewClient(config)
	if err != nil {
		log.QyLogger.Error("new consul client failed", zap.Error(err))
		panic(err)
	}

	r.client = client
}

func (r *consulResolver) watcher() {
	defer r.wg.Done()
	for {
		select {
		case <-r.stopC:
			return
		case <-r.rn:
			r.pullServer()
		case <-time.After(1 * time.Second):
			r.pullServer()
		}
	}
}

func (r *consulResolver) pullServer() {
	sn := fmt.Sprintf("qy.%v.%v.%v", os.Getenv("product"), r.serverVersion, r.serverName)

	services, metainfo, err := r.client.Health().Service(sn, r.serverVersion, true,
		&api.QueryOptions{
			WaitIndex: r.lastIndex, // 同步点，这个调用将一直阻塞，直到有新的更新
		})

	if err != nil {
		log.QyLogger.Error("error retrieving instances from consul", zap.Error(err))
		return
	}

	r.lastIndex = metainfo.LastIndex
	var newAddrs []resolver.Address
	for _, service := range services {
		addr := fmt.Sprintf("%v:%v", service.Service.Address, service.Service.Port)
		newAddrs = append(newAddrs, resolver.Address{Addr: addr})
	}

	log.QyLogger.Info(fmt.Sprintf("grpc resolver newAddrs: %v", newAddrs))

	r.cc.UpdateState(resolver.State{Addresses: newAddrs})
}
