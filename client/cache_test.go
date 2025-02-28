package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/require"
	"github.com/twistedogic/httputil/client/cache"
)

func Test_WithCache(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		c, err := cache.NewLRU(10)
		require.NoError(t, err, "setup cache")
		clk := clock.NewMock()
		clk.Set(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
		rt := cacheTransport{
			cache: responseCache{c: c, requestToKey: DefaultKeyFunc}, clk: clk,
			rt:  http.DefaultTransport,
			ttl: time.Hour,
		}
		h := setupHandler(t, 200, []byte("ok"), nil)
		ts := httptest.NewServer(h)
		defer ts.Close()
		for i := 0; i < 10; i++ {
			req := setupRequest(t, "GET", ts.URL, nil, nil)
			res := getResponse(t, rt, req)
			require.Equal(t, 200, res.StatusCode, "response code")
			if i != 0 {
				require.Equal(
					t, "true", res.Header.Get(fromCacheHeader), "response from cache",
				)
			}
		}
		require.Equal(t, 1, h.called, "should make call once only")
	})
}
