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
	CleanupDuration time.Duration) *BalancedCache[K, V] {

	cache := &BalancedCache[K, V]{
		Cache: *NewCache(TTLDuration, CleanupDuration),
		Mutex: &sync.RWMutex{},
		Map:   make(map[K]CacheItem[V]),
	}
	go cache.CleanUp()
	return cache
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

func (c *BalancedCache[K, V]) CleanUp() {
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

func (c *BalancedCache[K, V]) DeleteExpired() {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	for key, item := range c.Map {
		if item.TTL.Before(time.Now()) {
			delete(c.Map, key)
		}
	}
}

func (c *BalancedCache[K, V]) Stop() {
	close(c.stop)
}
