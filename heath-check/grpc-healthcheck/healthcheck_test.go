package grpc_healthcheck

import (
	"fmt"
	"github.com/dipperin/go-ms-toolkit/heath-check/config"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestRegisterHealthCheckWithGin(t *testing.T) {
	server := grpc.NewServer()
	RegisterHealthCheck(server)

	address, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", "8991"))
	assert.NoError(t, err)

	go func() {
		if err := server.Serve(address); err != nil {
			panic(err)
		}
	}()

	time.Sleep(30 * time.Millisecond)

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	RegisterHealthCheckWithGin("8991", engine)

	go func() {
		_ = engine.Run(":8992")
	}()

	time.Sleep(30 * time.Millisecond)

	client := resty.New()

	resp, err := client.R().Get("http://127.0.0.1:8992" + config.HealthCheckRoute)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	resp1, err1 := client.R().Post("http://127.0.0.1:8992" + config.HealthCheckRoute)
	assert.NoError(t, err1)
	assert.Equal(t, http.StatusOK, resp1.StatusCode())

	resp2, err2 := client.R().Head("http://127.0.0.1:8992" + config.HealthCheckRoute)
	assert.NoError(t, err2)
	assert.Equal(t, http.StatusOK, resp2.StatusCode())
}
