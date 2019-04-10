package g_metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var metrics map[string]interface{}
var enable bool = false

func init() {
	metrics = make(map[string]interface{})
}

func CreateCounter(name string, help string, label []string) {
	if !enable {
		return
	}
	if label == nil {
		counter := prometheus.NewCounter(prometheus.CounterOpts{
			Name: name,
			Help: help,
		})
		metrics[name] = counter
		prometheus.MustRegister(counter)
	} else {
		counter := prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: name,
			Help: help,
		}, label)
		metrics[name] = counter
		prometheus.MustRegister(counter)
	}
}

func CreateGauge(name string, help string, label []string) {
	if !enable {
		return
	}
	if label == nil {
		gauge := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: name,
			Help: help,
		})
		metrics[name] = gauge
		prometheus.MustRegister(gauge)
	} else {
		gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: name,
			Help: help,
		}, label)
		metrics[name] = gauge
		prometheus.MustRegister(gauge)
	}
}

func EnableMeter() {
	enable = true
}

func Set(name string, label string, value float64) {
	if !enable {
		return
	}
	if metrics[name] == nil {
		return
	}

	switch g := metrics[name].(type) {
	case *prometheus.GaugeVec:
		if label == "" {
			return
		}
		g.WithLabelValues(label).Set(value)
	case prometheus.Gauge:
		g.Set(value)
	default:
	}
}

func Add(name string, label string, value float64) {
	if !enable {
		return
	}
	if metrics[name] == nil {
		return
	}

	switch meter := metrics[name].(type) {
	case prometheus.Gauge:
		meter.Add(value)
	case prometheus.Counter:
		meter.Add(value)
	case *prometheus.GaugeVec:
		if label == "" {
			return
		}
		meter.WithLabelValues(label).Add(value)
	case *prometheus.CounterVec:
		if label == "" {
			return
		}
		meter.WithLabelValues(label).Add(value)
	}
}

func Sub(name string, label string, value float64) {
	if !enable {
		return
	}
	if metrics[name] == nil {
		return
	}

	switch meter := metrics[name].(type) {
	case prometheus.Gauge:
		meter.Sub(value)

	case *prometheus.GaugeVec:
		if label == "" {
			return
		}
		meter.WithLabelValues(label).Sub(value)

	default:
	}
}
