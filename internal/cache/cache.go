package cache

import (
	"time"
)

type CacheItem[V any] struct {
	Value V
	TTL   time.Time
}

type ICache[K comparable, V any] interface {
	Get(key K) (CacheItem[V], bool)
	Put(key K, value CacheItem[V]) bool
	Delete(key K) bool
	RefreshTTL(key K)
	CleanUp()
	DeleteExpired()
	Stop()
}

type Cache struct {
	TTLDuration     time.Duration
	CleanupDuration time.Duration
	stop            chan struct{}
}

func NewCache(
	TTLDuration time.Duration,
	CleanupDuration time.Duration) *Cache {

	return &Cache{
		TTLDuration:     TTLDuration,
		CleanupDuration: CleanupDuration,
		stop:            make(chan struct{}),
	}
}
