package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_WithRateLimit(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		rt := WithRateLimit(10, time.Second)(http.DefaultTransport)
		h := setupHandler(t, 200, []byte("ok"), nil)
		ts := httptest.NewServer(h)
		defer ts.Close()
		req := setupRequest(t, "GET", ts.URL, nil, nil)
		client := &http.Client{Transport: rt}
		for i := 0; i < 10; i++ {
			client.Do(req)
		}
		require.Equal(t, 10, h.called, "should make call once only")
	})
	t.Run("rate 0", func(t *testing.T) {
		rt := WithRateLimit(0, time.Second)(http.DefaultTransport)
		h := setupHandler(t, 200, []byte("ok"), nil)
		ts := httptest.NewServer(h)
		defer ts.Close()
		req := setupRequest(t, "GET", ts.URL, nil, nil)
		client := &http.Client{Transport: rt}
		for i := 0; i < 10; i++ {
			_, err := client.Do(req)
			require.Error(t, err, "expect rate limit")
		}
		require.Equal(t, 0, h.called, "should make call once only")
	})
}
