package client

import (
	"bufio"
	"bytes"
	"net/http"
	"time"

	"github.com/benbjohnson/clock"

	"github.com/twistedogic/httputil/client/cache"
)

const (
	fromCacheHeader = "X-HTTPCache"
)

func addCacheHeader(res *http.Response) {
	res.Header.Set(fromCacheHeader, "true")
}

func requestToKey(req *http.Request) string {
	if req == nil {
		return ""
	}
	return req.Method + "|" + req.URL.String()
}

type responseCache struct {
	c cache.Cache
}

func (r responseCache) Cache(
	req *http.Request, res *http.Response,
	expireAt time.Time,
) {
	key := requestToKey(req)
	buf := &bytes.Buffer{}
	if err := res.Write(buf); err == nil {
		r.c.Set(req.Context(), key, cache.Item{ExpireAt: expireAt, Content: buf.Bytes()})
	}
}

func (r responseCache) Fetch(req *http.Request) (*http.Response, bool) {
	key := requestToKey(req)
	item, exist := r.c.Get(req.Context(), key)
	if !exist {
		return nil, exist
	}
	buf := bytes.NewBuffer(item.Content)
	res, err := http.ReadResponse(bufio.NewReader(buf), req)
	if err != nil {
		return nil, false
	}
	addCacheHeader(res)
	return res, true
}

type cacheTransport struct {
	clk   clock.Clock
	rt    http.RoundTripper
	cache responseCache
	ttl   time.Duration
}

func (c cacheTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if res, ok := c.cache.Fetch(req); ok {
		return res, nil
	}
	res, err := c.rt.RoundTrip(req)
	if err == nil {
		c.cache.Cache(req, res, c.clk.Now().Add(c.ttl))
	}
	return res, err
}

func WithCache(c cache.Cache, ttl time.Duration) RoundTripWrapper {
	return func(rt http.RoundTripper) http.RoundTripper {
		return cacheTransport{
			clk:   clock.New(),
			rt:    rt,
			cache: responseCache{c: c},
			ttl:   ttl,
		}
	}
}
