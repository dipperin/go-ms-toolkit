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

func RegisterHealthCheck(server *grpc.Server) {
	// register health check for grpc
	grpc_health_v1.RegisterHealthServer(server, &healthCheckGRPCServerImpl{})
}

var (
	healthClient grpc_health_v1.HealthClient
)

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

		c.JSON(http.StatusOK, "pong")
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
