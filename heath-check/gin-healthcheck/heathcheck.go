package gin_healthcheck

import (
	"github.com/dipperin/go-ms-toolkit/heath-check/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

var pong gin.HandlerFunc = func(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}

// server has gin Engine
// 如果本身server是gin的话， 只需要把 engine丢进来， 注意不要带入gin middleware（类似鉴权的）
func RegisterHealthCheck(engine *gin.Engine) {
	engine.POST(config.HealthCheckRoute, pong)
	engine.GET(config.HealthCheckRoute, pong)
	engine.HEAD(config.HealthCheckRoute, pong)
}
