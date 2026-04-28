package cache

import (
	"sync"
	"time"
)

// Interface defines the contract for caching implementations
type Interface interface {
	// Get retrieves a value from the cache
	Get(key string) (interface{}, bool)
	// Set stores a value in the cache with TTL
	Set(key string, value interface{}, ttl time.Duration)
	// Delete removes a value from the cache
	Delete(key string)
	// Clear removes all values from the cache
	Clear()
}

// NoOpCache is a cache that doesn't cache anything
type NoOpCache struct{}

// Get always returns false
func (c *NoOpCache) Get(key string) (interface{}, bool) {
	return nil, false
}

// Set is a no-op
func (c *NoOpCache) Set(key string, value interface{}, ttl time.Duration) {
}

// Delete is a no-op
func (c *NoOpCache) Delete(key string) {
}

// Clear is a no-op
func (c *NoOpCache) Clear() {
}

// item represents a cached item with expiration time
type item struct {
	value      interface{}
	expiration int64
}

// InMemoryCache is a simple in-memory cache implementation with TTL support
type InMemoryCache struct {
	data map[string]item
	mu   interface {
		Lock()
		Unlock()
		RLock()
		RUnlock()
	}
}

// NewInMemoryCache creates a new in-memory cache
func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[string]item),
		mu:   &sync.RWMutex{},
	}
}

// Get retrieves a value from cache if it exists and hasn't expired
func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	it, found := c.data[key]
	if !found || (it.expiration > 0 && time.Now().UnixNano() > it.expiration) {
		return nil, false
	}
	return it.value, true
}

// Set stores a value in the cache with TTL
func (c *InMemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiration int64
	if ttl > 0 {
		expiration = time.Now().Add(ttl).UnixNano()
	}

	c.data[key] = item{
		value:      value,
		expiration: expiration,
	}
}

// Delete removes a value from the cache
func (c *InMemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

// Clear removes all values from the cache
func (c *InMemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]item)
}

// Close implements io.Closer for compatibility
func (c *InMemoryCache) Close() error {
	c.Clear()
	return nil
}
