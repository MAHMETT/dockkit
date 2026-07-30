package perf

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value     any
	ExpiresAt time.Time
}

type Cache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
	ttl   time.Duration
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		items: make(map[string]CacheItem),
		ttl:   ttl,
	}
}

func (c *Cache) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if !ok {
		return nil, false
	}
	if time.Now().After(item.ExpiresAt) {
		delete(c.items, key)
		return nil, false
	}
	return item.Value, true
}

func (c *Cache) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]CacheItem)
}

// Len returns the number of items in cache (including expired).
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// Delete removes a specific key.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}
