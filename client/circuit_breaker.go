package client

import (
	"fmt"
	"net/http"

	"github.com/sony/gobreaker"
)

type circuitBreaker struct {
	rt      http.RoundTripper
	breaker *gobreaker.CircuitBreaker
}

func (c circuitBreaker) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := c.breaker.Execute(func() (interface{}, error) {
		resp, err := c.rt.RoundTrip(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode >= 500 {
			return resp, fmt.Errorf("server returns %d", resp.StatusCode)
		}
		return resp, nil
	})
	if err != nil {
		return nil, err
	}
	r, ok := res.(*http.Response)
	if !ok {
		return nil, fmt.Errorf("response is not *http.Response")
	}
	return r, nil
}

func WithCircuitBreakerSetting(setting gobreaker.Settings) RoundTripWrapper {
	return func(rt http.RoundTripper) http.RoundTripper {
		return circuitBreaker{
			rt:      rt,
			breaker: gobreaker.NewCircuitBreaker(setting),
		}
	}
}

func WithCircuitBreaker(count int) RoundTripWrapper {
	if count <= 0 {
		count = 1
	}
	readyToTrip := func(counts gobreaker.Counts) bool {
		return counts.ConsecutiveFailures > uint32(count)
	}
	return WithCircuitBreakerSetting(gobreaker.Settings{ReadyToTrip: readyToTrip})
}
