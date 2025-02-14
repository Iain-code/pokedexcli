package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entry    map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{}

	cache.interval = interval
	cache.entry = make(map[string]cacheEntry)

	go cache.reapLoop(1 * time.Millisecond)

	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	cacheentry := cacheEntry{}
	cacheentry.val = val
	cacheentry.createdAt = time.Now()
	c.entry[key] = cacheentry
}

func (c *Cache) Get(key string) (val []byte, response bool) {

	c.mu.Lock()
	defer c.mu.Unlock()
	value, exists := c.entry[key]
	if !exists {
		return nil, false
	}
	return value.val, true
}

func (c *Cache) reapLoop(tickInterval time.Duration) {

	ticker := time.NewTicker(tickInterval)

	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()

		for key, value := range c.entry {
			if time.Since(value.createdAt) > c.interval {
				delete(c.entry, key)
			}
		}

		c.mu.Unlock()
	}
}
