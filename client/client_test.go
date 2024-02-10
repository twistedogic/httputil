package client

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockHandler struct {
	header       map[string][]string
	code, called int
	content      []byte
}

func setupHandler(t *testing.T, code int, content []byte, header map[string][]string) *mockHandler {
	t.Helper()
	return &mockHandler{
		code: code, header: header, content: content,
	}
}

func (m *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.called += 1
	for k, values := range m.header {
		for _, v := range values {
			w.Header().Set(k, v)
		}
	}
	w.WriteHeader(m.code)
	w.Write(m.content)
	return
}

func setupRequest(t *testing.T, method, target string, body []byte, header map[string][]string) *http.Request {
	req, err := http.NewRequest(method, target, bytes.NewBuffer(body))
	require.NoError(t, err, "create request")
	for k, values := range header {
		for _, v := range values {
			req.Header.Set(k, v)
		}
	}
	return req
}

func getResponse(t *testing.T, rt http.RoundTripper, req *http.Request) *http.Response {
	client := &http.Client{Transport: rt}
	res, err := client.Do(req)
	require.NoError(t, err, "make request")
	return res
}
