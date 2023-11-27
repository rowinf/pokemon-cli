package pokecache

import (
	"fmt"
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

func NewCache(interval time.Duration) Cache {
	c := Cache{}
	c.mu = sync.Mutex{}
	c.entries = make(map[string]cacheEntry)
	reapLoop(&c, interval)
	return c
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
	done := make(chan bool)
	go func() {
		time.Sleep(interval)
		done <- true
	}()
	fmt.Println("start reaping")
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			c.mu.Lock()
			for key, entry := range c.entries {
				if t.After(entry.createdAt.Add(interval)) {
					fmt.Println("reap: ", entry)
					delete(c.entries, key)
				}
			}
			c.mu.Unlock()
		}
	}
}
