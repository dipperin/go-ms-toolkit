package gin_healthcheck

import (
	"github.com/dipperin/go-ms-toolkit/heath-check/config"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestRegisterHealthCheck(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	RegisterHealthCheck(engine)
	go func() {
		_ = engine.Run(":9991")
	}()

	time.Sleep(30 * time.Millisecond)

	client := resty.New()

	resp, err := client.R().Get("http://127.0.0.1:9991" + config.HealthCheckRoute)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	resp1, err1 := client.R().Post("http://127.0.0.1:9991" + config.HealthCheckRoute)
	assert.NoError(t, err1)
	assert.Equal(t, http.StatusOK, resp1.StatusCode())

	resp2, err2 := client.R().Head("http://127.0.0.1:9991" + config.HealthCheckRoute)
	assert.NoError(t, err2)
	assert.Equal(t, http.StatusOK, resp2.StatusCode())
}
