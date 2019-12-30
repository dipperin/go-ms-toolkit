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
func RegisterHealthCheck(engine *gin.Engine) {
	engine.POST(config.HealthCheckRoute, pong)
	engine.GET(config.HealthCheckRoute, pong)
	engine.HEAD(config.HealthCheckRoute, pong)
}
