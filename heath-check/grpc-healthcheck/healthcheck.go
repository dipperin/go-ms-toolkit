package grpc_healthcheck

import (
	"github.com/dipperin/go-ms-toolkit/heath-check/config"
	"github.com/dipperin/go-ms-toolkit/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net/http"
)

// 这个是在开启grpc health check， 本身grpc 就支持健康检查协议， 如果k8s能直接调同的话 use grpc_health_probe
// 那就只用调用这个方法就好了
func RegisterHealthCheck(server *grpc.Server) {
	// register health check for grpc
	grpc_health_v1.RegisterHealthServer(server, &healthCheckGRPCServerImpl{})
}

var (
	healthClient grpc_health_v1.HealthClient
)

// 这是一个wrap 方法， 开启一个gin，可以同http的方法去调用grpc的健康检查，注意： 必须启动完grpc server 后才能调用
// 先决条件：
// 1. 开启了grpc health check, 已经调用了 RegisterHealthCheck
// 2. grpc server 已经Serve起来了
// 如果不遵循先决条件，可能会panic
func RegisterHealthCheckWithGin(grpcPort string, engine *gin.Engine) {
	conn, err := grpc.Dial("localhost:"+grpcPort, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	healthClient = grpc_health_v1.NewHealthClient(conn)

	// 先校验一次是否是通的
	initResp, err := healthClient.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		log.QyLogger.Error("grpc check health gin wrap init failed", zap.Error(err))
		panic(err)
	}

	if initResp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		log.QyLogger.Warn("grpc health check failed, server no serving")
		panic("grpc health check failed, server no serving")
	}

	var pong gin.HandlerFunc = func(c *gin.Context) {
		resp, err := healthClient.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
		if err != nil {
			log.QyLogger.Error("grpc check health gin wrap failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"err_msg": err.Error()})
			return
		}

		if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
			log.QyLogger.Error("grpc health check failed, server no serving")
			c.JSON(http.StatusInternalServerError, gin.H{"err_msg": "grpc health check failed, server no serving"})
			return
		}

		c.JSON(http.StatusOK, "OK")
	}

	engine.POST(config.HealthCheckRoute, pong)
	engine.GET(config.HealthCheckRoute, pong)
	engine.HEAD(config.HealthCheckRoute, pong)
}

// impl grpc_health_v1.HealthServer
// 如果以后还有多的需求，可以改写这个server
type healthCheckGRPCServerImpl struct {
}

func (s *healthCheckGRPCServerImpl) Check(context.Context, *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	log.QyLogger.Debug("check rpc server health.....")
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *healthCheckGRPCServerImpl) Watch(*grpc_health_v1.HealthCheckRequest, grpc_health_v1.Health_WatchServer) error {
	log.QyLogger.Warn("no support health check --- watch")
	return nil
}
