package decodejson

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	contenttype "github.com/qdm12/REPONAME_GITHUB/internal/server/contenttypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_DecodeBody(t *testing.T) {
	t.Parallel()

	type exampleStruct struct {
		A int `json:"a"`
	}

	testCases := map[string]struct {
		maxBytes     int64
		requestBody  string
		v            interface{}
		expectedV    interface{}
		ok           bool
		status       int
		responseBody string
	}{
		"success": {
			maxBytes:    1024,
			requestBody: `{"a":1}`,
			v:           &exampleStruct{},
			expectedV:   &exampleStruct{A: 1},
			ok:          true,
			status:      http.StatusOK,
		},
		"max size": {
			maxBytes:    2,
			requestBody: `{"a":1}`,
			v:           &exampleStruct{},
			ok:          false,
			status:      http.StatusRequestEntityTooLarge,
			responseBody: `{"error":"request body is too large"}
`,
		},
		"unknown field": {
			maxBytes:    1024,
			requestBody: `{"a":1,"b":2}`,
			v:           &exampleStruct{},
			ok:          false,
			status:      http.StatusBadRequest,
			responseBody: `{"error":"request body contains unknown field \"b\""}
`,
		},
		"extra after JSON": {
			maxBytes:    1024,
			requestBody: `{"a":1}\n`,
			v:           &exampleStruct{},
			ok:          false,
			status:      http.StatusBadRequest,
			responseBody: `{"error":"request body must only contain a single JSON object"}
`,
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			requestBody := ioutil.NopCloser(
				strings.NewReader(testCase.requestBody),
			)

			ok := DecodeBody(
				w, testCase.maxBytes, requestBody, testCase.v, contenttype.JSON)

			require.Equal(t, testCase.ok, ok)
			bytes, err := ioutil.ReadAll(w.Body)
			require.NoError(t, err)
			responseBody := string(bytes)
			assert.Equal(t, testCase.status, w.Code)
			assert.Equal(t, testCase.responseBody, responseBody)
		})
	}
}

func Test_extractFromJSONErr(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		err       error
		errString string
		errCode   int
	}{
		"syntax error": {
			err:       &json.SyntaxError{Offset: 1},
			errString: "request body contains badly-formed JSON (at position 1)",
			errCode:   http.StatusBadRequest,
		},
		"unexpected EOF": {
			err:       io.ErrUnexpectedEOF,
			errString: "request body contains badly-formed JSON",
			errCode:   http.StatusBadRequest,
		},
		"unmarshal type error": {
			err:       &json.UnmarshalTypeError{},
			errString: `request body contains an invalid value for the "" field (at position 0)`,
			errCode:   http.StatusBadRequest,
		},
		"unknown field": {
			err:       errors.New("json: unknown field bla"),
			errString: `request body contains unknown field bla`,
			errCode:   http.StatusBadRequest,
		},
		"EOF": {
			err:       io.EOF,
			errString: `request body cannot be empty`,
			errCode:   http.StatusBadRequest,
		},
		"request too large": {
			err:       errors.New("http: request body too large"),
			errString: `request body is too large`,
			errCode:   http.StatusRequestEntityTooLarge,
		},
		"default error": {
			err:       errors.New("some error"),
			errString: `some error`,
			errCode:   http.StatusInternalServerError,
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			errString, errCode := extractFromJSONErr(testCase.err)
			assert.Equal(t, testCase.errString, errString)
			assert.Equal(t, testCase.errCode, errCode)
		})
	}
}
