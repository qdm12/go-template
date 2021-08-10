package metrics

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/qdm12/go-template/internal/metrics/mock_metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_metricsMiddleware(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	const responseStatus = http.StatusBadRequest
	responseBody := []byte{1, 2, 3, 4}
	req := httptest.NewRequest(http.MethodGet, "https://test.com", nil)
	const routePattern = "/test"
	chiCtx := chi.NewRouteContext()
	chiCtx.RoutePatterns = []string{routePattern}
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx)
	req = req.Clone(ctx)

	childHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, req.Method, r.Method)
			assert.Equal(t, req.URL, r.URL)
			w.WriteHeader(responseStatus)
			_, _ = w.Write(responseBody)
		},
	)

	timeNowIndex := 0
	timeNow := func() time.Time {
		unix := int64(4156132 + timeNowIndex)
		timeNowIndex++
		return time.Unix(unix, 0)
	}

	metrics := mock_metrics.NewMockInterface(ctrl)
	metrics.EXPECT().InflightRequestsGaugeAdd(1)
	metrics.EXPECT().InflightRequestsGaugeAdd(-1)
	metrics.EXPECT().RequestCountInc(routePattern, responseStatus)
	metrics.EXPECT().ResponseBytesCountAdd(routePattern, responseStatus, len(responseBody))
	metrics.EXPECT().ResponseTimeHistogramObserve(routePattern, responseStatus, time.Second)

	handler := &metricsHandler{
		childHandler: childHandler,
		metrics:      metrics,
		timeNow:      timeNow,
	}

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	result := w.Result()
	b, err := ioutil.ReadAll(result.Body)
	require.NoError(t, err)
	defer result.Body.Close()
	assert.Equal(t, responseStatus, result.StatusCode)
	assert.Equal(t, responseBody, b)
}
