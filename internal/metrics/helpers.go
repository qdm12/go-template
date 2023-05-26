package metrics

import (
	"errors"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	ErrRegister = errors.New("registration error")
)

func newCounterVec(name, help string, labelNames []string, register bool) (c *prometheus.CounterVec, err error) {
	c = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: promNamespace,
		Subsystem: promSubsystem,
		Name:      name,
		Help:      help,
	}, labelNames)
	if register {
		if err := prometheus.Register(c); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrRegister, err)
		}
	}
	return c, nil
}

func newGauge(name, help string, register bool) (g prometheus.Gauge, err error) { //nolint:ireturn
	g = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Subsystem: promSubsystem,
		Name:      name,
		Help:      help,
	})
	if register {
		if err := prometheus.Register(g); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrRegister, err)
		}
	}
	return g, nil
}

func newHistogramVec(name, help string, buckets []float64, labelNames []string, register bool) (
	h *prometheus.HistogramVec, err error) {
	h = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: promNamespace,
		Subsystem: promSubsystem,
		Name:      name,
		Help:      help,
		Buckets:   buckets,
	}, labelNames)
	if register {
		if err := prometheus.Register(h); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrRegister, err)
		}
	}
	return h, nil
}

func newResponseTimeHistogramVec(register bool) (responseTimeHistogram *prometheus.HistogramVec, err error) {
	//nolint:gomnd
	buckets := []float64{
		float64(time.Millisecond),
		float64(10 * time.Millisecond),
		float64(50 * time.Millisecond),
		float64(100 * time.Millisecond),
		float64(150 * time.Millisecond),
		float64(200 * time.Millisecond),
		float64(500 * time.Millisecond),
		float64(750 * time.Millisecond),
		float64(time.Second),
		float64(2 * time.Second),
		float64(5 * time.Second),
		float64(10 * time.Second),
	}
	return newHistogramVec("response_time",
		"Histogram for the response times by handler and HTTP status",
		buckets,
		[]string{"handler", "status"},
		register)
}
