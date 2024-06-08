package web

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheckEndPoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	HealthCheckEndpoint(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(io.Reader(res.Body))

	require.Nil(t, err)

	assert.Contains(t, string(data), "Keep calm I'm absolutely alive 🐛")
}
