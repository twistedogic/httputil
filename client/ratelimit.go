package client

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type ratelimitTransport struct {
	rt      http.RoundTripper
	limiter *rate.Limiter
}

func (r ratelimitTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := r.limiter.Wait(req.Context()); err != nil {
		return nil, err
	}
	return r.rt.RoundTrip(req)
}

func WithRateLimit(size int, interval time.Duration) RoundTripWrapper {
	limiter := rate.NewLimiter(rate.Every(interval), size)
	return func(rt http.RoundTripper) http.RoundTripper {
		return ratelimitTransport{rt: rt, limiter: limiter}
	}
}
