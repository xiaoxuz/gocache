package cache

import (
	"testing"
	"time"
)

func TestLRUCache(t *testing.T) {

	cache := NewLRUCache(3, 2, 1*time.Second)

	cache.Set("mykey", []string{"hello", "cache"})

	go func() {
		if _, ok := cache.Get("mykey"); !ok {
			t.Fail()
		}
		if _, ok := cache.Get("mykey"); !ok {
			t.Fail()
		}
		// age > 2 fail
		if _, ok := cache.Get("mykey"); ok {
			t.Fail()
		}
	}()
	time.Sleep(1 * time.Second)
	// 2次 Get 后基于 lru 删除 存在则 fail
	if _, ok := cache.Get("mykey"); ok {
		t.Fail()
	}

	// 统计
	if i, q, h := cache.Count(); i != 0 || q != 4 || h != 2 {
		t.Fail()
	}
}
