package cache

import (
	"testing"
)

func Test_LRU(t *testing.T) {
	c, err := NewLRU(10)
	if err != nil {
		t.Fatal(err)
	}
	testCache(t, c)
}
