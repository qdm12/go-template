package cors

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_setCrossOriginHeaders(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		requestHeaders          http.Header
		responseHeaders         http.Header
		allowedOrigins          map[string]struct{}
		allowedHeaders          []string
		expectedResponseHeaders http.Header
		err                     error
	}{
		"no origin": {
			requestHeaders:          http.Header{},
			responseHeaders:         http.Header{},
			expectedResponseHeaders: http.Header{},
			err:                     errNoOriginHeader,
		},
		"origin not allowed": {
			requestHeaders: http.Header{
				"Origin": []string{"not-allowed"},
			},
			responseHeaders:         http.Header{},
			allowedOrigins:          map[string]struct{}{},
			expectedResponseHeaders: http.Header{},
			err:                     errors.New("origin not allowed: not-allowed"),
		},
		"origin allowed": {
			requestHeaders: http.Header{
				"Origin": []string{"allowed"},
			},
			responseHeaders: http.Header{
				"key": []string{"value"},
			},
			allowedOrigins: map[string]struct{}{"allowed": {}},
			allowedHeaders: []string{"Authorization", "HeaderKey"},
			expectedResponseHeaders: http.Header{
				"key":                          []string{"value"},
				"Access-Control-Allow-Origin":  []string{"allowed"},
				"Access-Control-Max-Age":       []string{"14400"},
				"Access-Control-Allow-Headers": []string{"Authorization", "HeaderKey"},
			},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			requestHeaders := testCase.requestHeaders.Clone()
			responseHeaders := testCase.responseHeaders

			err := setCrossOriginHeaders(requestHeaders, responseHeaders,
				testCase.allowedOrigins, testCase.allowedHeaders)

			if testCase.err != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.requestHeaders, requestHeaders)
			assert.Equal(t, testCase.expectedResponseHeaders, responseHeaders)
		})
	}
}

func Test_AllowCORSMethods(t *testing.T) {
	t.Parallel()

	optionsRequestWithHeader := func(header http.Header) (req *http.Request) {
		req = httptest.NewRequest(http.MethodOptions, "http://test.com", nil)
		req.Header = header
		return req
	}

	testCases := map[string]struct {
		request    *http.Request
		methods    []string
		respHeader http.Header
		status     int
	}{
		"no request method header": {
			request: optionsRequestWithHeader(http.Header{}),
			methods: []string{"POST", "PUT"},
			respHeader: http.Header{
				"Access-Control-Allow-Methods": []string{"POST, PUT, OPTIONS"},
			},
			status: http.StatusMethodNotAllowed,
		},
		"not accepted request method header": {
			request: optionsRequestWithHeader(http.Header{
				"Access-Control-Request-Method": []string{"DELETE"},
			}),
			methods: []string{"POST", "PUT"},
			respHeader: http.Header{
				"Access-Control-Allow-Methods": []string{"POST, PUT, OPTIONS"},
			},
			status: http.StatusMethodNotAllowed,
		},
		"accepted request method header": {
			request: optionsRequestWithHeader(http.Header{
				"Access-Control-Request-Method": []string{"POST"},
			}),
			methods: []string{"POST", "PUT"},
			respHeader: http.Header{
				"Access-Control-Allow-Methods": []string{"POST, PUT, OPTIONS"},
			},
			status: http.StatusOK,
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			writer := httptest.NewRecorder()

			AllowCORSMethods(testCase.request, writer, testCase.methods...)

			result := writer.Result()
			result.Body.Close() // for the linter
			assert.Equal(t, testCase.status, result.StatusCode)
			assert.Equal(t, testCase.respHeader, result.Header)
		})
	}
}
