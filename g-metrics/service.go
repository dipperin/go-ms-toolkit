package g_metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// This method should be placed at the forefront to ensure that it can be call before other services are registered.
func NewPrometheusMetricsServer(port int) *PrometheusMetricsServer {
	pms := &PrometheusMetricsServer{
		port: port,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%v", port),
			Handler: promhttp.Handler(),
		},
	}
	if port != 0 {
		EnableMeter()
	}
	return pms
}

type PrometheusMetricsServer struct {
	port   int
	server *http.Server
}

func (p *PrometheusMetricsServer) Start() error {
	if p.port == 0 {
		fmt.Println("port is 0, do not start prometheus metrics server")
		return nil
	}

	fmt.Println("start prometheus metrics", "addr", p.server.Addr)
	go func() {
		if err := p.server.ListenAndServe(); err != nil {
			fmt.Println("pMetrics serve failed: " + err.Error())
		}
	}()
	return nil
}

func (p *PrometheusMetricsServer) Stop() {}
