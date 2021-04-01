package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func newCounterVec(name, help string, labelNames []string, register bool) (c *prometheus.CounterVec, err error) {
	c = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: promNamespace,
		Subsystem: promSubsystem,
		Name:      name,
		Help:      help,
	}, labelNames)
	if register {
		err = prometheus.Register(c)
	}
	return c, err
}

func newGauge(name, help string, register bool) (g prometheus.Gauge, err error) {
	g = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Subsystem: promSubsystem,
		Name:      name,
		Help:      help,
	})
	if register {
		err = prometheus.Register(g)
	}
	return g, err
}

func newHistogramVec(name, help string, buckets []float64, labelNames []string, register bool) (
	histogram *prometheus.HistogramVec, err error) {
	histogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: promNamespace,
		Subsystem: promSubsystem,
		Name:      name,
		Help:      help,
		Buckets:   buckets,
	}, labelNames)
	if register {
		err = prometheus.Register(histogram)
	}
	return histogram, err
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
