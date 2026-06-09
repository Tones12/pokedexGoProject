package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	val			[]byte
	createdAt	time.Time
}

type Cache struct {
	mu sync.Mutex
	entries	map[string]cacheEntry
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := cacheEntry{
		val: val,
		createdAt: time.Now(),
	}
	c.entries[key] = entry
}


func (c *Cache) Get(key string) ([]byte, bool) {	
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.entries {
			if time.Since(entry.createdAt) > interval {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}

func NewCache(interval time.Duration) *Cache {
	newCache := Cache{
		entries: make(map[string]cacheEntry),
	}
	go newCache.reapLoop(interval)
	return &newCache
}