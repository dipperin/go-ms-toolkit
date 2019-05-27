package g_metrics_v2

import (
	"testing"
)

func TestProm_WithCounter(t *testing.T) {
	c := New().WithCounter("xx_count")
	c.Add(2)
}
