package metrics

import "time"

type Metrics interface {
	RequestCountInc(routePattern string, statusCode int)
	ResponseBytesCountAdd(routePattern string, statusCode int, bytesWritten int)
	InflightRequestsGaugeAdd(addition int)
	ResponseTimeHistogramObserve(routePattern string, statusCode int, duration time.Duration)
}
