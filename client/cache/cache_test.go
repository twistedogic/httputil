package cache

import (
	"context"
	"reflect"
	"testing"
	"time"
)

type entry struct {
	key  string
	item Item
}

func testCache(t *testing.T, c Cache) {
	cases := map[string]struct {
		items []entry
		key   string
		exist bool
		want  Item
	}{
		"base": {
			items: []entry{
				{
					key: "a",
					item: Item{
						ExpireAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
						Content:  []byte(`ok`),
					},
				},
			},
			key:   "a",
			exist: true,
			want: Item{
				ExpireAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				Content:  []byte(`ok`),
			},
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			for _, v := range tc.items {
				c.Set(context.TODO(), v.key, v.item)
			}
			got, ok := c.Get(context.TODO(), tc.key)
			if tc.exist != ok {
				t.Fatalf("key %q, exist to be %v, got: %v", tc.key, tc.exist, ok)
			}
			if ok && !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("want: %v, got: %v", tc.want, got)
			}
		})
	}
}
