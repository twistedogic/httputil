package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockHandler struct {
	statusCode int
	content    []byte
}

func (m mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(m.statusCode)
	w.Write(m.content)
}

func Test_WithNoCache(t *testing.T) {
	h := mockHandler{statusCode: http.StatusOK, content: []byte("ok")}
	handler := WithNoCache()(h)
	t.Run("set no cache", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		handler.ServeHTTP(rec, req)
		values := rec.Result().Header.Values(cacheControlKey)
		require.ElementsMatch(
			t,
			[]string{"max-age=0", noCacheValue, noStoreValue, mustValidateValue},
			values,
			"cache control directives",
		)
	})
}

func Test_WithEtagCache(t *testing.T) {
	content := []byte("ok")
	tag := "tag"
	h := mockHandler{statusCode: http.StatusOK, content: content}
	handler := WithEtagCache(tag)(h)
	t.Run("set etag", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		handler.ServeHTTP(rec, req)
		require.Equal(t, tag, rec.Result().Header.Get(etagKey), "set etag")
		require.Equal(t, 200, rec.Result().StatusCode, "set etag")
		require.Equal(t, content, rec.Body.Bytes(), "set etag with content")
	})
	t.Run("invalid etag", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add(ifNoneMatchKey, "different tag")
		handler.ServeHTTP(rec, req)
		require.Equal(t, tag, rec.Result().Header.Get(etagKey), "set etag")
		require.Equal(t, 200, rec.Result().StatusCode, "revalidate etag")
		require.Equal(t, content, rec.Body.Bytes(), "revalidate content")
	})
	t.Run("valid etag", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add(ifNoneMatchKey, tag)
		handler.ServeHTTP(rec, req)
		require.Equal(t, 304, rec.Result().StatusCode, "validate etag")
		require.Equal(t, tag, rec.Result().Header.Get(etagKey), "set etag")
		require.Nil(t, rec.Body.Bytes(), "reuse cache")
	})
}
