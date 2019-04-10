package g_metrics

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestNewPrometheusMetricsServer(t *testing.T) {
	//port must be 0 for this error test
	svr := NewPrometheusMetricsServer(0)
	assert.NotNil(t, svr)

	assert.Nil(t, svr.Start())

	svr = NewPrometheusMetricsServer(35231)
	svr.Start()
	time.Sleep(100 * time.Millisecond)
	svr.Start()
	svr.Stop()
}

func TestCreateCounter(t *testing.T) {
	//test enable switch close
	CreateCounter("fail", "fail", nil)

	EnableMeter()
	//create counter
	CreateCounter("count", "count", nil)

	//create counter vector
	CreateCounter("countVector", "countVector", []string{"countvVectorS"})
}

func TestCreateGauge(t *testing.T) {
	EnableMeter()
	//create gauge
	CreateGauge("gauge", "gauge", nil)

	//create gauge vector
	CreateGauge("gaugeVector", "gaugeVector", []string{"gaugeVectorS"})
}

type testMeter struct {
	name string
	help string
	label []string
}

var meters []testMeter

func init() {
	meters = makeMeters()
}

func makeMeters() []testMeter {
	meters := []testMeter{
		{"cnt", "cnt", nil},
		{"cntV", "cntV", []string{"cv"}},
		{"gg", "gg", nil},
		{"ggV", "ggV", []string{"gv"}},
	}
	EnableMeter()
	CreateCounter(meters[0].name, meters[0].help, meters[0].label)
	CreateCounter(meters[1].name, meters[1].help, meters[1].label)
	CreateGauge(meters[2].name, meters[2].help, meters[2].label)
	CreateGauge(meters[3].name, meters[3].help, meters[3].label)

	return meters
}

func TestAdd(t *testing.T) {
	//test non-register
	Add("notExist", "", 1)

	for _, m := range meters {
		if len(m.label) > 0 {
			//for vector, test no label
			Add(m.name, "", 1)
			//with lable
			Add(m.name, m.label[0], 1)
		} else {
			Add(m.name, "", 1)
		}
	}
}

func TestSet(t *testing.T) {
	//test non-register
	Set("notExist", "", 1)

	for _, m := range meters {
		if len(m.label) > 0 {
			//for vector, test no label
			Set(m.name, "", 1)
			//with lable
			Set(m.name, m.label[0], 1)
		} else {
			Set(m.name, "", 1)
		}
	}
}

func TestSub(t *testing.T) {
	//test non-register
	Sub("notExist", "", 1)

	for _, m := range meters {
		if len(m.label) > 0 {
			//for vector, test no label
			Sub(m.name, "", 1)
			//with lable
			Sub(m.name, m.label[0], 1)
		} else {
			Sub(m.name, "", 1)
		}
	}
}

func TestDisable(t *testing.T) {
	enable = false
	CreateCounter("", "", nil)
	CreateGauge("", "", nil)
	Set("", "", 1)
	Add("", "", 1)
	Sub("", "", 1)
}