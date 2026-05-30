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
	CleanupDuration time.Duration) *ReadHeavyCache[K, V] {

	cache := &ReadHeavyCache[K, V]{
		Cache: *NewCache(TTLDuration, CleanupDuration),
	}

	go cache.CleanUp()
	return cache
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
		entry := val.(*CacheItem[V])
		entry.TTL = time.Now().Add(c.TTLDuration)
		c.Map.Store(key, entry)
	}
}

func (c *ReadHeavyCache[K, V]) CleanUp() {
	ticker := time.NewTicker(c.CleanupDuration)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-c.stop:
			ticker.Stop()
			return
		}
	}
}

func (c *ReadHeavyCache[K, V]) DeleteExpired() {
	for key, item := range c.Map.Range {
		if item.(*CacheItem[V]).TTL.Before(time.Now()) {
			c.Map.Delete(key)
		}
	}
}

func (c *ReadHeavyCache[K, V]) Stop() {
	close(c.stop)
}
