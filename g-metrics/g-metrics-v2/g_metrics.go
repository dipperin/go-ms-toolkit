package g_metrics_v2

import "github.com/prometheus/client_golang/prometheus"

// Prom struct info
type Prom struct {
	counter prometheus.Counter
	state prometheus.Gauge

	timerVec   *prometheus.HistogramVec
	counterVec *prometheus.CounterVec
	stateVec   *prometheus.GaugeVec
}

// New creates a Prom instance.
func New() *Prom {
	return &Prom{}
}

// WithTimerVec with summary timerVec
func (p *Prom) WithTimerVec(name string, labels []string) *Prom {
	if p == nil || p.timerVec != nil {
		return p
	}
	p.timerVec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: name,
			Help: name,
		}, labels)
	prometheus.MustRegister(p.timerVec)
	return p
}

// WithCounterVec sets counterVec.
func (p *Prom) WithCounterVec(name string, labels []string) *Prom {
	if p == nil || p.counterVec != nil {
		return p
	}
	p.counterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name,
			Help: name,
		}, labels)
	prometheus.MustRegister(p.counterVec)
	return p
}

// WithStateVec sets stateVec.
func (p *Prom) WithStateVec(name string, labels []string) *Prom {
	if p == nil || p.stateVec != nil {
		return p
	}
	p.stateVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: name,
		}, labels)
	prometheus.MustRegister(p.stateVec)
	return p
}

// WithCounter sets counter.
func (p *Prom) WithCounter(name string) *Prom {
	if p == nil || p.counter != nil {
		return p
	}
	p.counter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: name,
			Help: name,
		})
	prometheus.MustRegister(p.counter)
	return p
}

// WithState sets state.
func (p *Prom) WithState(name string) *Prom {
	if p == nil || p.state != nil {
		return p
	}
	p.state = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: name,
			Help: name,
		})
	prometheus.MustRegister(p.state)
	return p
}

// Incr increments one stat counterVec without sampling
func (p *Prom) Incr() {
	if p.counter != nil {
		p.counter.Inc()
	}
	if p.state != nil {
		p.state.Inc()
	}
}

// Decr decrements one stat counterVec without sampling
func (p *Prom) Decr() {
	if p.state != nil {
		p.state.Dec()
	}
}

// State set stateVec
func (p *Prom) State(v int64) {
	if p.state != nil {
		p.state.Set(float64(v))
	}
}

// Add add count    v must > 0
func (p *Prom) Add(v int64) {
	if p.counter != nil {
		p.counter.Add(float64(v))
	}
	if p.state != nil {
		p.state.Add(float64(v))
	}
}

// IncrVec increments one stat counterVec without sampling
func (p *Prom) IncrVec(name string, extra ...string) {
	label := append([]string{name}, extra...)
	if p.counterVec != nil {
		p.counterVec.WithLabelValues(label...).Inc()
	}
	if p.stateVec != nil {
		p.stateVec.WithLabelValues(label...).Inc()
	}
}

// DecrVec decrements one stat counterVec without sampling
func (p *Prom) DecrVec(name string, extra ...string) {
	if p.stateVec != nil {
		label := append([]string{name}, extra...)
		p.stateVec.WithLabelValues(label...).Dec()
	}
}

// StateVec set stateVec
func (p *Prom) StateVec(name string, v int64, extra ...string) {
	if p.stateVec != nil {
		label := append([]string{name}, extra...)
		p.stateVec.WithLabelValues(label...).Set(float64(v))
	}
}

// AddVec add count    v must > 0
func (p *Prom) AddVec(name string, v int64, extra ...string) {
	label := append([]string{name}, extra...)
	if p.counterVec != nil {
		p.counterVec.WithLabelValues(label...).Add(float64(v))
	}
	if p.stateVec != nil {
		p.stateVec.WithLabelValues(label...).Add(float64(v))
	}
}
