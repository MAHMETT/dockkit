package registry

import (
	"sync"
	"time"
)

type cacheItem struct {
	tags      []TagInfo
	expiresAt time.Time
}

type tagCache struct {
	items map[string]cacheItem
	mu    sync.RWMutex
	ttl   time.Duration
}

func newTagCache(ttl time.Duration) *tagCache {
	return &tagCache{
		items: make(map[string]cacheItem),
		ttl:   ttl,
	}
}

func (c *tagCache) Get(image string) ([]TagInfo, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[image]
	if !ok {
		return nil, false
	}
	if time.Now().After(item.expiresAt) {
		delete(c.items, image)
		return nil, false
	}
	return item.tags, true
}

func (c *tagCache) Set(image string, tags []TagInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[image] = cacheItem{
		tags:      tags,
		expiresAt: time.Now().Add(c.ttl),
	}
}

func (c *tagCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]cacheItem)
}

// Len returns the number of items in cache.
func (c *tagCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}
