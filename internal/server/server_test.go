package server

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/go-template/internal/config"
	"github.com/qdm12/go-template/internal/metrics/mock_metrics"
	"github.com/qdm12/go-template/internal/models"
	"github.com/qdm12/go-template/internal/processor/mock_processor"
	"github.com/qdm12/golibs/logging/mock_logging"
	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	config := config.HTTP{
		Address: "test",
	}
	logger := mock_logging.NewMockLogger(ctrl)
	metrics := mock_metrics.NewMockInterface(ctrl)
	buildInformation := models.BuildInformation{}
	proc := mock_processor.NewMockInterface(ctrl)

	server := New(config, proc, logger, metrics, buildInformation)
	assert.Equal(t, config.Address, server.address)
	assert.Equal(t, logger, server.logger)
	assert.NotNil(t, server.handler)
}

func Test_server_Run(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	const address = "127.0.0.1:9000"

	logger := mock_logging.NewMockLogger(ctrl)
	logger.EXPECT().Info("listening on " + address)

	server := &Server{
		address: address,
		handler: nil,
		logger:  logger,
	}

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error)

	go func() {
		errCh <- server.Run(ctx)
	}()

	cancel()
	err := <-errCh
	assert.NoError(t, err)
}
