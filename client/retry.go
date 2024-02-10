package client

import (
	"net/http"

	rhttp "github.com/hashicorp/go-retryablehttp"
)

func WithRetryClient(client *rhttp.Client) RoundTripWrapper {
	return func(rt http.RoundTripper) http.RoundTripper {
		client.HTTPClient.Transport = rt
		return &rhttp.RoundTripper{Client: client}
	}
}
func WithRetry(retries int) RoundTripWrapper {
	client := rhttp.NewClient()
	client.RetryMax = retries
	return WithRetryClient(client)
}
