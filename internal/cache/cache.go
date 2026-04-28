package cache

import (
	"sync"
	"time"
)

type item struct {
	value      interface{}
	expiration int64
}

type Cache struct {
	data map[string]item
	mu   sync.RWMutex
}

func New() *Cache {
	return &Cache{data: make(map[string]item)}
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = item{
		value:      value,
		expiration: time.Now().Add(ttl).UnixNano(),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	it, found := c.data[key]
	if !found || time.Now().UnixNano() > it.expiration {
		return nil, false
	}
	return it.value, true
}
