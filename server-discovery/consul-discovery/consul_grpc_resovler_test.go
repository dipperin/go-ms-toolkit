package consul_discovery

import (
	"fmt"
	"google.golang.org/grpc/resolver"
	"strings"
	"sync"
	"testing"
	"time"
)

func Test_consulResolver_watcher(t *testing.T) {
	d := &ddd{stopC: make(chan struct{}), rn: make(chan struct{})}

	d.wg.Add(1)
	go d.watcher()

	time.Sleep(200 * time.Microsecond)

	d.rn <- struct{}{}

	time.Sleep(50 * time.Microsecond)

	time.Sleep(1 * time.Second)

	close(d.stopC)

	d.wg.Wait()
}

type ddd struct {
	rn    chan struct{}
	stopC chan struct{}
	wg    sync.WaitGroup
}

func (r *ddd) watcher() {
	defer r.wg.Done()
	for {
		select {
		case <-r.stopC:
			println("aaaaaa close")
			return
		case <-r.rn:
			println("aaaaaa rn")
			//r.pullServer()
		case <-time.After(500 * time.Millisecond):
			println("aaaaaa time")
			//r.pullServer()
		}
	}
}

type tCc struct {
}

func (c *tCc) UpdateState(s resolver.State) {
	fmt.Println(s)
}

func (c *tCc) NewAddress(addresses []resolver.Address) {
	panic("implement me")
}

func (c *tCc) NewServiceConfig(serviceConfig string) {
	panic("implement me")
}

func Test_newConsulResolver(t *testing.T) {
	r := newConsulResolver("127.0.0.1:8500", "cd-demo", "v0.1", &tCc{})

	time.Sleep(25 *time.Second)

	r.Close()
}


// split2 returns the values from strings.SplitN(s, sep, 2).
// If sep is not found, it returns ("", "", false) instead.
func split2(s, sep string) (string, string, bool) {
	spl := strings.SplitN(s, sep, 2)
	if len(spl) < 2 {
		return "", "", false
	}
	return spl[0], spl[1], true
}

// parseTarget splits target into a struct containing scheme, authority and
// endpoint.
//
// If target is not a valid scheme://authority/endpoint, it returns {Endpoint:
// target}.
func parseTarget(target string) (ret resolver.Target) {
	var ok bool
	ret.Scheme, ret.Endpoint, ok = split2(target, "://")
	if !ok {
		return resolver.Target{Endpoint: target}
	}
	ret.Authority, ret.Endpoint, ok = split2(ret.Endpoint, "/")
	if !ok {
		return resolver.Target{Endpoint: target}
	}
	return ret
}

func TestParseTarget(t *testing.T) {
	fmt.Println(parseTarget("consul://127.0.0.1:8500/helloworld@v0.1"))
	fmt.Println(splitVersion("helloworld", "@"))
}