package cache

import "sync"

type CacheMap struct {
	Sessions     map[string]Session
	ChangeStream chan CacheOperation
}

var lock = &sync.Mutex{}
var instance *CacheMap

func getInstance() *CacheMap {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		instance = &CacheMap{
			Sessions:     make(map[string]Session),
			ChangeStream: make(chan CacheOperation),
		}
	}
	return instance
}

type CacheOperation struct {
	Operation string
	Key       string
	Vault     struct{}
	Reply     chan CacheOperationResult
}

type CacheOperationResult struct {
	Value   struct{}
	Message string
}
