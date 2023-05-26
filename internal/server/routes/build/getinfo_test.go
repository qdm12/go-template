package build

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/qdm12/go-template/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_handler_getBuild(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		handler        *handler
		buildRequest   func() *http.Request
		expectedStatus int
		expectedBody   string
	}{
		"no_response_content_type_supported": {
			handler: &handler{},
			buildRequest: func() *http.Request {
				request := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
				request.Header.Set("Accept", "blah")
				return request
			},
			expectedStatus: http.StatusNotAcceptable,
			expectedBody:   `{"error":"no response content type supported: blah"}` + "\n",
		},
		"success": {
			handler: &handler{
				build: models.BuildInformation{
					Version:   "1.2.3",
					Commit:    "abcdef",
					BuildDate: "2023-05-26T00:00:00Z",
				},
			},
			buildRequest: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "http://test.com", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"version":"1.2.3","commit":"abcdef","buildDate":"2023-05-26T00:00:00Z"}` + "\n",
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			writer := httptest.NewRecorder()
			request := testCase.buildRequest()

			testCase.handler.getBuild(writer, request)

			result := writer.Result()
			body, err := io.ReadAll(result.Body)
			closeErr := result.Body.Close()
			require.NoError(t, err)
			assert.NoError(t, closeErr)

			assert.Equal(t, testCase.expectedStatus, result.StatusCode)
			assert.Equal(t, testCase.expectedBody, string(body))
		})
	}
}
