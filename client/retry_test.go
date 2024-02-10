package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_WithRetry(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		rt := WithRetry(1)(http.DefaultTransport)
		h := setupHandler(t, 500, []byte("not ok"), nil)
		ts := httptest.NewServer(h)
		defer ts.Close()
		req := setupRequest(t, "GET", ts.URL, nil, nil)
		client := &http.Client{Transport: rt}
		client.Do(req)
		require.Equal(t, 2, h.called, "should make call once only")
	})
}
