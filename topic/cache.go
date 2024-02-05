package topic

import (
	"sync"
)

var topics *Cache

type Cache struct {
	store sync.Map
}

func init() {
	topics = new(Cache)
}

func (t *Cache) Set(key string, value interface{}) {
	t.store.Store(key, value)
}

func (t *Cache) Get(key string) (interface{}, bool) {
	return t.store.Load(key)
}
