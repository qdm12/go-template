package httperr

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/qdm12/go-template/internal/server/contenttype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Respond(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		status       int
		errString    string
		expectedBody string
	}{
		"status without error string": {
			status: http.StatusBadRequest,
			expectedBody: `{"error":"Bad Request"}
`,
		},
		"status with error string": {
			status:    http.StatusBadRequest,
			errString: "bad parameter",
			expectedBody: `{"error":"bad parameter"}
`,
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()

			err := Respond(w, testCase.status, testCase.errString, contenttype.JSON)
			require.NoError(t, err)

			response := w.Result()
			defer response.Body.Close()
			bytes, err := io.ReadAll(response.Body)
			require.NoError(t, err)
			body := string(bytes)

			assert.Equal(t, testCase.status, response.StatusCode)
			assert.Equal(t, testCase.expectedBody, body)
		})
	}
}
