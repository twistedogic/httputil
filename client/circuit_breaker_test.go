package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_WithCircuitBreaker(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		rt := WithCircuitBreaker(3)(http.DefaultTransport)
		h := setupHandler(t, 500, []byte("not ok"), nil)
		ts := httptest.NewServer(h)
		defer ts.Close()
		req := setupRequest(t, "GET", ts.URL, nil, nil)
		client := &http.Client{Transport: rt}
		for i := 0; i < 10; i++ {
			client.Do(req)
		}
		require.Equal(t, 4, h.called, "should make call once only")
	})
}
