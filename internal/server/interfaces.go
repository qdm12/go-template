package server

import (
	"context"
	"time"

	"github.com/qdm12/go-template/internal/models"
)

type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Error(s string)
}

type Metrics interface {
	RequestCountInc(routePattern string, statusCode int)
	ResponseBytesCountAdd(routePattern string, statusCode int, bytesWritten int)
	InflightRequestsGaugeAdd(addition int)
	ResponseTimeHistogramObserve(routePattern string, statusCode int, duration time.Duration)
}

type Processor interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUserByID(ctx context.Context, id uint64) (user models.User, err error)
}
