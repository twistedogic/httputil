package cache

import (
	"context"

	"github.com/hashicorp/golang-lru/v2/simplelru"
)

type LRU struct {
	lru *simplelru.LRU[string, Item]
}

func NewLRU(size int) (Cache, error) {
	lru, err := simplelru.NewLRU[string, Item](size, nil)
	return LRU{lru: lru}, err
}

func (l LRU) Get(ctx context.Context, key string) (Item, bool) {
	select {
	case <-ctx.Done():
		return Item{}, false
	default:
		return l.lru.Get(key)
	}
}
func (l LRU) Set(ctx context.Context, key string, val Item) {
	select {
	case <-ctx.Done():
		return
	default:
		l.lru.Add(key, val)
	}
}
