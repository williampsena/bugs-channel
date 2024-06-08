package web

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	t.Setenv("WEB_RATE_LIMIT", "100")

	svr := buildTestServer(t)

	defer svr.Close()

	requestURL := fmt.Sprintf("%v/health", svr.URL)

	res, err := http.Get(requestURL)

	require.Nil(t, err)

	assert.Equal(t, 200, res.StatusCode)
}

func TestServerRateLimit(t *testing.T) {
	t.Setenv("WEB_RATE_LIMIT", "1")

	svr := buildTestServer(t)

	defer svr.Close()

	doRequest := func(index int) *http.Response {
		requestURL := fmt.Sprintf("%v/health?%v", svr.URL, index)
		res, err := http.Get(requestURL)

		require.Nil(t, err)

		return res
	}

	res := doRequest(1)

	assert.Equal(t, 200, res.StatusCode)

	res = doRequest(-1)

	assert.Equal(t, 429, res.StatusCode)
}

func buildTestServer(t *testing.T) *httptest.Server {
	router, err := buildRouter()

	require.Nil(t, err)

	return httptest.NewServer(router)

}
