package cache

import (
	"log"
	"sync"
	"time"
)

type ReadHeavyCache[K comparable, V any] struct {
	Cache
	Map sync.Map
}

func NewReadHeavyCache[K comparable, V any](
	TTLDuration time.Duration,
	Cleanup time.Duration) *ReadHeavyCache[K, V] {

	return &ReadHeavyCache[K, V]{
		Cache: *NewCache(TTLDuration, Cleanup),
	}
}

func (c *ReadHeavyCache[K, V]) Get(key K) (CacheItem[V], bool) {
	var empty CacheItem[V]
	val, ok := c.Map.Load(key)
	if !ok {
		log.Printf("entry with key %v not found", key)
		return empty, false
	}

	return val.(CacheItem[V]), true
}

func (c *ReadHeavyCache[K, V]) Put(key K, value CacheItem[V]) bool {
	c.Map.Store(key, value)
	_, ok := c.Map.Load(key)
	return ok
}

func (c *ReadHeavyCache[K, V]) Delete(key K) bool {
	c.Map.Delete(key)
	_, ok := c.Map.Load(key)
	return !ok
}

func (c *ReadHeavyCache[K, V]) RefreshTTL(key K) {
	if val, ok := c.Map.Load(key); ok {
		entry := val.(*CacheItem[V]) // type-assert to your entry type
		entry.TTL = time.Now().Add(c.TTLDuration)
		c.Map.Store(key, entry)
	}
}
