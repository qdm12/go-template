package build

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/go-template/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_handler(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		makeLogger     func(ctrl *gomock.Controller) *MockLogger
		buildInfo      models.BuildInformation
		method         string
		requestHeader  http.Header
		expectedStatus int
		expectedHeader http.Header
		expectedBody   string
	}{
		"get_build": {
			makeLogger: func(_ *gomock.Controller) *MockLogger {
				return nil
			},
			buildInfo: models.BuildInformation{
				Version:   "1.2.3",
				Commit:    "abcdef",
				BuildDate: "2023-05-26T00:00:00Z",
			},
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedHeader: http.Header{
				"Content-Length": []string{"73"},
				"Content-Type":   []string{"application/json"},
			},
			expectedBody: `{"version":"1.2.3","commit":"abcdef",` +
				`"buildDate":"2023-05-26T00:00:00Z"}` + "\n",
		},
		"options": {
			makeLogger: func(_ *gomock.Controller) *MockLogger {
				return nil
			},
			buildInfo: models.BuildInformation{
				Version:   "1.2.3",
				Commit:    "abcdef",
				BuildDate: "2023-05-26T00:00:00Z",
			},
			method: http.MethodOptions,
			requestHeader: http.Header{
				"Access-Control-Request-Method": []string{http.MethodOptions},
			},
			expectedStatus: http.StatusOK,
			expectedHeader: http.Header{
				"Access-Control-Allow-Methods": []string{http.MethodGet + ", " + http.MethodOptions},
				"Content-Length":               []string{"0"},
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			logger := testCase.makeLogger(ctrl)
			handler := NewHandler(logger, testCase.buildInfo)

			server := httptest.NewServer(handler)
			t.Cleanup(server.Close)
			client := server.Client()

			ctx := context.Background()
			testDeadline, ok := t.Deadline()
			if ok {
				var cancel context.CancelFunc
				ctx, cancel = context.WithDeadline(context.Background(), testDeadline)
				defer cancel()
			}

			request, err := http.NewRequestWithContext(ctx,
				testCase.method, server.URL, nil)
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
