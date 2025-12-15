package cardingester

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

type Cache struct {
	cache map[string]cacheEntry
	mux *sync.Mutex
}

func NewCache(interval time.Duration) Cache {
	c := Cache{
		cache: make(map[string]cacheEntry),
		mux:   &sync.Mutex{},
	}

	go c.reapLoop(interval)

	return c
}

func (cache *Cache) Add(key string, value []byte) {
	cache.mux.Lock()
	cache.cache[key] = cacheEntry{
		createdAt: time.Now().UTC(),
		val: value,
	}
	cache.mux.Unlock()
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mux.Lock()
	defer cache.mux.Unlock()
	value, result := cache.cache[key]

	return value.val, result
}
func (cache *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		cache.reap(time.Now().UTC(), interval)
	}
}

func (cache *Cache) reap(now time.Time, last time.Duration) {
	cache.mux.Lock()
	defer cache.mux.Unlock()

	for k, v := range cache.cache {
		if v.createdAt.Before(now.Add(-last)) {
			delete(cache.cache, k)
		}
	}
}
