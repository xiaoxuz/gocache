// @Date: 2020-06-16

// LRU Cache
package cache

import (
	"container/list"
	"sync"
	"time"
)

type LRUCache struct {
	lruMap   map[interface{}]*item
	lruList  *list.List
	expires  time.Duration
	age      int
	size     int
	m        sync.RWMutex
	cntQuery uint64
	cntHit   uint64
}

type item struct {
	el         *list.Element
	age        int
	expiration time.Time
	val        interface{}
}

func NewLRUCache(size int, age int, expires time.Duration) *LRUCache {
	if size < 0 {
		panic("Size of LRUCache should not less than zero")
	}
	if age < 0 {
		panic("Age of LRUCache should not less than zero")
	}
	if expires < 0 {
		panic("Expire of LRUCache should not less than zero")
	}
	return &LRUCache{
		lruMap:  make(map[interface{}]*item),
		lruList: list.New(),
		size:    size,
		age:     age,
		expires: expires,
	}
}

func (lc *LRUCache) Get(key interface{}) (interface{}, bool) {
	lc.m.Lock()
	defer lc.m.Unlock()

	lc.cntQuery++
	if v, ok := lc.lruMap[key]; ok {
		// check age
		if lc.age > 0 {
			if v.age >= lc.age {
				delete(lc.lruMap, key)
				lc.lruList.Remove(v.el)
				return nil, false
			}
			v.age++
		}
		// check expir
		if lc.expires > 0 {
			if v.expiration.Before(time.Now()) {
				delete(lc.lruMap, key)
				lc.lruList.Remove(v.el)
				return nil, false
			}
		}
		// hit
		lc.lruList.MoveToFront(v.el)
		lc.cntHit++
		return v.val, true
	}
	return nil, false
}

func (lc *LRUCache) Set(key, val interface{}) {
	lc.m.Lock()
	defer lc.m.Unlock()

	if v, ok := lc.lruMap[key]; ok {
		v.age = 0
		v.expiration = time.Now().Add(lc.expires)
		lc.lruList.MoveToFront(v.el)
		return
	}

	el := lc.lruList.PushFront(key)
	lc.lruMap[key] = &item{
		el:         el,
		val:        val,
		age:        0,
		expiration: time.Now().Add(lc.expires),
	}

	for lc.lruList.Len() > lc.size {
		last := lc.lruList.Back()
		delete(lc.lruMap, last.Value)
		lc.lruList.Remove(last)
	}
}

func (lc *LRUCache) Purge() {
	lc.m.Lock()
	defer lc.m.Unlock()

	lc.lruMap = make(map[interface{}]*item)
	lc.lruList = list.New()
}

func (lc *LRUCache) Count() (cntItems, cntQuery, cntHit uint64) {
	lc.m.RLock()
	defer lc.m.RUnlock()

	return uint64(lc.lruList.Len()), lc.cntQuery, lc.cntHit
}
