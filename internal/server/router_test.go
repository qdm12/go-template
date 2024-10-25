package server

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"net/url"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/go-template/internal/config/settings"
	"github.com/qdm12/go-template/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ptrTo[T any](x T) *T { return &x }

func Test_Router(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		config         settings.HTTP
		makeLogger     func(ctrl *gomock.Controller) *MockLogger
		makeMetrics    func(ctrl *gomock.Controller) *MockMetrics
		makeProcessor  func(ctrl *gomock.Controller) *MockProcessor
		buildInfo      models.BuildInformation
		path           string
		method         string
		requestBody    string
		requestHeader  http.Header
		expectedStatus int
		expectedHeader http.Header
		expectedBody   string
	}{
		"get_build": {
			config: settings.HTTP{
				RootURL:     ptrTo("/rooturl"),
				LogRequests: ptrTo(true),
			},
			makeLogger: func(ctrl *gomock.Controller) *MockLogger {
				logger := NewMockLogger(ctrl)
				logger.EXPECT().Infof("HTTP request: %d %s %s %s %s %dms",
					http.StatusOK, http.MethodGet, "/rooturl/api/v1/build",
					netip.AddrFrom4([4]byte{127, 0, 0, 1}), "73B",
					gomock.AssignableToTypeOf(time.Second))
				return logger
			},
			makeMetrics: func(ctrl *gomock.Controller) *MockMetrics {
				metrics := NewMockMetrics(ctrl)
				metrics.EXPECT().InflightRequestsGaugeAdd(1)
				metrics.EXPECT().InflightRequestsGaugeAdd(-1)
				metrics.EXPECT().RequestCountInc("/rooturl/api/v1/build", http.StatusOK)
				metrics.EXPECT().ResponseBytesCountAdd("/rooturl/api/v1/build",
					http.StatusOK, 73)
				metrics.EXPECT().ResponseTimeHistogramObserve("/rooturl/api/v1/build",
					http.StatusOK, gomock.AssignableToTypeOf(time.Second))
				return metrics
			},
			makeProcessor: func(_ *gomock.Controller) *MockProcessor {
				return nil
			},
			buildInfo: models.BuildInformation{
				Version:   "1.2.3",
				Commit:    "abcdef",
				BuildDate: "2023-05-26T00:00:00Z",
			},
			path:           "/rooturl/api/v1/build",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedHeader: http.Header{
				"Content-Length": []string{"73"},
				"Content-Type":   []string{"application/json"},
			},
			expectedBody: `{"version":"1.2.3","commit":"abcdef",` +
				`"buildDate":"2023-05-26T00:00:00Z"}` + "\n",
		},
		"create_user": {
			config: settings.HTTP{
				RootURL:     ptrTo("/"),
				LogRequests: ptrTo(true),
			},
			makeLogger: func(ctrl *gomock.Controller) *MockLogger {
				logger := NewMockLogger(ctrl)
				logger.EXPECT().Infof("HTTP request: %d %s %s %s %s %dms",
					http.StatusCreated, http.MethodPost, "/api/v1/users",
					netip.AddrFrom4([4]byte{127, 0, 0, 1}), "0B",
					gomock.AssignableToTypeOf(time.Second))
				return logger
			},
			makeMetrics: func(ctrl *gomock.Controller) *MockMetrics {
				metrics := NewMockMetrics(ctrl)
				metrics.EXPECT().InflightRequestsGaugeAdd(1)
				metrics.EXPECT().InflightRequestsGaugeAdd(-1)
				metrics.EXPECT().RequestCountInc("/api/v1/users", http.StatusCreated)
				metrics.EXPECT().ResponseBytesCountAdd("/api/v1/users",
					http.StatusCreated, 0)
				metrics.EXPECT().ResponseTimeHistogramObserve("/api/v1/users",
					http.StatusCreated, gomock.AssignableToTypeOf(time.Second))
				return metrics
			},
			makeProcessor: func(ctrl *gomock.Controller) *MockProcessor {
				processor := NewMockProcessor(ctrl)
				expectedUser := models.User{
					ID:       1,
					Account:  "admin",
					Username: "qdm12",
				}
				processor.EXPECT().CreateUser(gomock.Any(), expectedUser).
					Return(nil)
				return processor
			},
			path:           "/api/v1/users",
			method:         http.MethodPost,
			requestBody:    `{"id": 1, "account": "admin", "username": "qdm12"}`,
			expectedStatus: http.StatusCreated,
			expectedHeader: http.Header{
				"Content-Length": []string{"0"},
				"Content-Type":   []string{"application/json"},
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			logger := testCase.makeLogger(ctrl)
			metrics := testCase.makeMetrics(ctrl)
			processor := testCase.makeProcessor(ctrl)
			router := NewRouter(testCase.config, logger, metrics,
				testCase.buildInfo, processor)

			server := httptest.NewServer(router)
			t.Cleanup(server.Close)
			client := server.Client()

			ctx := context.Background()
			testDeadline, ok := t.Deadline()
			if ok {
				var cancel context.CancelFunc
				ctx, cancel = context.WithDeadline(context.Background(), testDeadline)
				defer cancel()
			}

			body := bytes.NewBufferString(testCase.requestBody)
			url, err := url.Parse(server.URL)
			require.NoError(t, err)
			url.Path = testCase.path
			request, err := http.NewRequestWithContext(ctx,
				testCase.method, url.String(), body)
			require.NoError(t, err)
			for key, values := range testCase.requestHeader {
				for _, value := range values {
					request.Header.Add(key, value)
				}
			}

			response, err := client.Do(request)
			require.NoError(t, err)
			t.Cleanup(func() {
				_ = response.Body.Close()
			})

			for expectedKey, expectedValues := range testCase.expectedHeader {
				values := response.Header.Values(expectedKey)
				assert.Equal(t, expectedValues, values)
			}
			assert.Equal(t, testCase.expectedStatus, response.StatusCode)
			responseData, err := io.ReadAll(response.Body)
			require.NoError(t, err)
			assert.Equal(t, testCase.expectedBody, string(responseData))

			err = response.Body.Close()
			require.NoError(t, err)
		})
	}
}
