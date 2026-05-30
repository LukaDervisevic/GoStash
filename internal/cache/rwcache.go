package cache

import (
	"log"
	"sync"
	"time"
)

type BalancedCache[K comparable, V any] struct {
	Cache
	Mutex *sync.RWMutex
	Map   map[K]CacheItem[V]
}

func NewBalancedCache[K comparable, V any](
	TTLDuration time.Duration,
	Cleanup time.Duration) *BalancedCache[K, V] {

	return &BalancedCache[K, V]{
		Cache: *NewCache(TTLDuration, Cleanup),
		Mutex: &sync.RWMutex{},
		Map:   make(map[K]CacheItem[V]),
	}
}

func (c *BalancedCache[K, V]) Get(key K) (CacheItem[V], bool) {
	var empty CacheItem[V]
	c.Mutex.RLock()
	val, ok := c.Map[key]
	if !ok {
		log.Printf("entry with key %v not found", key)
		c.Mutex.RUnlock()
		return empty, false
	}

	if val.TTL.Before(time.Now()) {
		log.Printf("entry with key %v has expired... removing key", key)
		c.Mutex.RUnlock()
		c.Delete(key)
		return empty, false
	}

	c.Mutex.Unlock()
	return c.Map[key], true
}

func (c *BalancedCache[K, V]) Put(key K, value CacheItem[V]) bool {
	c.Mutex.Lock()
	if _, ok := c.Map[key]; ok {
		log.Printf("entry with key %v already exists in cache", key)
		c.Mutex.Unlock()
		return false
	}

	if value.TTL.IsZero() {
		value.TTL = time.Now().Add(c.TTLDuration)
	}
	c.Map[key] = value
	c.Mutex.Unlock()
	return true
}

func (c *BalancedCache[K, V]) Delete(key K) bool {
	c.Mutex.Lock()
	delete(c.Map, key)
	c.Mutex.Unlock()

	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	_, ok := c.Map[key]
	return !ok
}

func (c *BalancedCache[K, V]) RefreshTTL(key K) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Map[key].TTL.Add(c.TTLDuration)
}
