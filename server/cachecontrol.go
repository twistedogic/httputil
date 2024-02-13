package server

import (
	"net/http"
)

const (
	cacheControlKey = "Cache-Control"

	etagKey        = "Etag"
	ifNoneMatchKey = "If-None-Match"

	mustValidateValue = "must-revalidate"
	noCacheValue      = "no-cache"
	noStoreValue      = "no-store"
	maxAgeValue       = "max-age"
)

func Etag(h http.Header) string {
	if val := h.Get(ifNoneMatchKey); val != "" {
		return val
	}
	return h.Get(etagKey)
}
func SetEtag(h http.Header, etag string) { h.Add(etagKey, etag) }

func ForceValidate(h http.Header) {
	h.Del(cacheControlKey)
	h.Add(cacheControlKey, maxAgeValue+"=0")
	h.Add(cacheControlKey, noCacheValue)
	h.Add(cacheControlKey, noStoreValue)
	h.Add(cacheControlKey, mustValidateValue)
}

func WithEtagCache(etag string) HandlerWrapper {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet, http.MethodHead:
				if rh := r.Header; etag == Etag(rh) {
					SetEtag(w.Header(), etag)
					w.WriteHeader(http.StatusNotModified)
					return
				}
			}
			SetEtag(w.Header(), etag)
			h.ServeHTTP(w, r)
		})
	}
}

func WithNoCache() HandlerWrapper {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ForceValidate(w.Header())
			h.ServeHTTP(w, r)
		})
	}
}
