package httperr

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	contenttype "github.com/qdm12/REPONAME_GITHUB/internal/server/contenttypes"
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
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()

			Respond(w, testCase.status, testCase.errString, contenttype.JSON)

			response := w.Result()
			defer response.Body.Close()
			bytes, err := ioutil.ReadAll(response.Body)
			require.NoError(t, err)
			body := string(bytes)

			assert.Equal(t, testCase.status, response.StatusCode)
			assert.Equal(t, testCase.expectedBody, body)
		})
	}
}
