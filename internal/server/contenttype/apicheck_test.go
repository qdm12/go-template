package contenttype

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_APICheck(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		requestHeader       http.Header
		requestContentType  string
		responseContentType string
		err                 error
	}{
		"empty header": {
			requestHeader:       http.Header{},
			requestContentType:  "application/json",
			responseContentType: "application/json",
		},
		"header with valid accept and content type": {
			requestHeader: http.Header{
				"Accept":       []string{"application/json"},
				"Content-Type": []string{"application/json"},
			},
			requestContentType:  "application/json",
			responseContentType: "application/json",
		},
		"header with html accept": {
			requestHeader: http.Header{
				"Accept": []string{"text/html"},
			},
			requestContentType:  "application/json",
			responseContentType: "application/json",
		},
		"header with invalid accept": {
			requestHeader: http.Header{
				"Accept":       []string{"invalid, invalid2"},
				"Content-Type": []string{"application/json"},
			},
			responseContentType: "application/json",
			err:                 errors.New("no response content type supported: invalid, invalid2"),
		},
		"header with one valid accept of many": {
			requestHeader: http.Header{
				"Accept":       []string{"invalid,application/json ,invalid"},
				"Content-Type": []string{"application/json"},
			},
			requestContentType:  "application/json",
			responseContentType: "application/json",
		},
		"header with invalid content type": {
			requestHeader: http.Header{
				"Content-Type": []string{"invalid"},
			},
			responseContentType: "application/json",
			err:                 errors.New(`content type is not supported: "invalid"`),
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			requestContentType, responseContentType, err := APICheck(testCase.requestHeader)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.requestContentType, requestContentType)
			assert.Equal(t, testCase.responseContentType, responseContentType)
		})
	}
}
