package server

import (
	"net/http"
)

type HandlerWrapper func(http.Handler) http.Handler

func Compose(h http.Handler, wrappers ...HandlerWrapper) http.Handler {
	out := h
	for _, w := range wrappers {
		out = w(h)
	}
	return out
}
