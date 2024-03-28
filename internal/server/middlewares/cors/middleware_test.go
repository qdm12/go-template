package cors

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_corsHandler(t *testing.T) {
	t.Parallel()
	allowedOrigins := []string{"http://test"}
	allowedHeaders := []string{"Authorization"}
	middleware := New(allowedOrigins, allowedHeaders)

	childHandler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
	handler := middleware(childHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	ctx := context.Background()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	request.Header.Set("origin", "http://test")

	client := server.Client()
	response, err := client.Do(request)
	require.NoError(t, err)
	_ = response.Body.Close()

	response.Header.Del("Date")

	expectedResponseHeader := http.Header{
		"Access-Control-Allow-Origin":  []string{"http://test"},
		"Access-Control-Max-Age":       []string{"14400"},
		"Content-Length":               []string{"0"},
		"Access-Control-Allow-Headers": []string{"Authorization"},
	}

	assert.Equal(t, expectedResponseHeader, response.Header)
}
