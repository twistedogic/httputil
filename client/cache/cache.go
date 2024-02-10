package cache

import (
	"context"
	"fmt"
	"time"
)

type Item struct {
	ExpireAt time.Time
	Content  []byte
}

func (i Item) String() string {
	return fmt.Sprintf("%s %s", i.ExpireAt, i.Content)
}

type Cache interface {
	Get(context.Context, string) (Item, bool)
	Set(context.Context, string, Item)
}
