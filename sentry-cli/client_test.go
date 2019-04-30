package sentry_cli

import (
	"errors"
	"testing"
)

func TestClient(t *testing.T) {
	//os.Setenv("sentry_dsn", "http://e428e37420954758ad932a9af3fda39f:e6a916f8a3e64060a5e95627b538e544@114.119.116.157:9394/3")
	Client().CaptureMessage("asdasdsad", nil)
	Client().CaptureError(errors.New("xx"), nil)
	Client().Wait()
}
