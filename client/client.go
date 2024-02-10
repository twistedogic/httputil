package client

import (
	"net/http"
)

type RoundTripWrapper func(http.RoundTripper) http.RoundTripper

func Compose(rt http.RoundTripper, wrappers ...RoundTripWrapper) http.RoundTripper {
	out := rt
	for _, w := range wrappers {
		out = w(out)
	}
	return out
}
