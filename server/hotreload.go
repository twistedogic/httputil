package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"golang.org/x/net/html"
)

func toScriptNode(content string) *html.Node {
	return &html.Node{
		Data: "script",
		Type: html.ElementNode,
		FirstChild: &html.Node{
			Data: content,
			Type: html.TextNode,
		},
	}
}

func bfsFirstElementNode(name string, root *html.Node) (*html.Node, error) {
	queue := []*html.Node{root}
	var current *html.Node
	for len(queue) != 0 {
		current, queue = queue[0], queue[1:]
		if current.Type == html.ElementNode && current.Data == name {
			return current, nil
		}
		if next := current.NextSibling; next != nil {
			queue = append(queue, next)
		}
		if child := current.FirstChild; child != nil {
			queue = append(queue, child)
		}
	}
	return nil, fmt.Errorf("no element node with %q found", name)
}

func injectClient(content string, r io.Reader, w io.Writer) error {
	n, err := html.Parse(r)
	if err != nil {
		return err
	}
	head, err := bfsFirstElementNode("head", n)
	if err != nil {
		return err
	}
	head.AppendChild(toScriptNode(content))
	return html.Render(w, n)
}

type responseWrapper struct {
	w   http.ResponseWriter
	buf *bytes.Buffer
}

func wrapResponseWriter(w http.ResponseWriter) responseWrapper {
	return responseWrapper{w: w, buf: &bytes.Buffer{}}
}

// implements http.ResponseWriter
func (r responseWrapper) Write(b []byte) (int, error) { return r.buf.Write(b) }
func (r responseWrapper) Header() http.Header         { return r.w.Header() }
func (r responseWrapper) WriteHeader(statusCode int)  { r.w.WriteHeader(statusCode) }

func (r responseWrapper) inject(content string) error {
	return injectClient(content, r.buf, r.w)
}

func WithHotReload(rootRoute, script string, handler http.Handler) http.Handler {
	upgrader := websocket.Upgrader{}
	id := uuid.New().String()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case rootRoute:
			ww := wrapResponseWriter(w)
			handler.ServeHTTP(ww, r)
			if err := ww.inject(script); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		case "/ws/hotreload":
			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			defer c.Close()
			for {
				mt, _, err := c.ReadMessage()
				if err != nil {
					http.Error(w, err.Error(), 500)
					break
				}
				if err := c.WriteMessage(mt, []byte(id)); err != nil {
					http.Error(w, err.Error(), 500)
					break
				}
			}
		default:
			handler.ServeHTTP(w, r)
		}
	})
}

func WithDefaultHotReload(handler http.Handler) http.Handler {
	return WithHotReload("/", clientScript, handler)
}
