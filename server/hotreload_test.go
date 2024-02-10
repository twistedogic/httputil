package server

import (
	"bytes"
	"strings"
	"testing"
)

func Test_injectClientHotReload(t *testing.T) {
	cases := map[string]struct {
		input, content, want string
	}{
		"base": {
			input:   `<html><head></head><body><h1>hi</h1></body></html>`,
			content: `alert("h1")`,
			want:    `<html><head><script>alert("h1")</script></head><body><h1>hi</h1></body></html>`,
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			r, w := strings.NewReader(tc.input), &bytes.Buffer{}
			if err := injectClient(tc.content, r, w); err != nil {
				t.Fatal(err)
			}
			if got := w.String(); got != tc.want {
				t.Fatalf("want: %s, got: %s", tc.want, got)
			}
		})
	}
}
