package internal

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	mu      sync.Mutex
	entries map[string]cacheEntry
}

func NewCache(interval time.Duration) *Cache {
	c := Cache{}
	c.mu = sync.Mutex{}
	c.entries = make(map[string]cacheEntry)
	go reapLoop(&c, interval)
	return &c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = cacheEntry{time.Now(), val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entries[key]
	return entry.val, ok
}

func reapLoop(c *Cache, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for t := range ticker.C {
		c.mu.Lock()
		for key, entry := range c.entries {
			if t.After(entry.createdAt.Add(interval)) {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}
