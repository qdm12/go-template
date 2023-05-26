// Package metrics contains a metrics interface with methods to modify the
// metrics for Prometheus.
package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	requestsCounter       *prometheus.CounterVec
	responseBytesCounter  *prometheus.CounterVec
	inFlighRequestsGauge  prometheus.Gauge
	responseTimeHistogram *prometheus.HistogramVec
}

func New(register bool) (m *Metrics, err error) {
	requestsCounter, err := newCounterVec(
		"requests",
		"Counter for the number of requests by handler and HTTP status",
		[]string{"handler", "status"}, register)
	if err != nil {
		return nil, err
	}
	responseBytesCounter, err := newCounterVec(
		"response_bytes",
		"Counter for the number of bytes written in the response by handler and HTTP status",
		[]string{"handler", "status"}, register)
	if err != nil {
		return nil, err
	}
	inFlighRequestsGauge, err := newGauge(
		"requests_inflight",
		"Gauge for the current number of inflight requests by handler and HTTP status",
		register)
	if err != nil {
		return nil, err
	}
	responseTimeHistogram, err := newResponseTimeHistogramVec(register)
	if err != nil {
		return nil, err
	}

	return &Metrics{
		requestsCounter:       requestsCounter,
		responseBytesCounter:  responseBytesCounter,
		inFlighRequestsGauge:  inFlighRequestsGauge,
		responseTimeHistogram: responseTimeHistogram,
	}, nil
}

func (m *Metrics) RequestCountInc(routePattern string, statusCode int) {
	m.requestsCounter.WithLabelValues(routePattern, http.StatusText(statusCode)).Inc()
}

func (m *Metrics) ResponseBytesCountAdd(routePattern string, statusCode int, bytesWritten int) {
	m.responseBytesCounter.WithLabelValues(routePattern, http.StatusText(statusCode)).Add(float64(bytesWritten))
}

func (m *Metrics) InflightRequestsGaugeAdd(addition int) {
	m.inFlighRequestsGauge.Add(float64(addition))
}

func (m *Metrics) ResponseTimeHistogramObserve(routePattern string, statusCode int, duration time.Duration) {
	m.responseTimeHistogram.WithLabelValues(routePattern, http.StatusText(statusCode)).Observe(duration.Seconds())
}
