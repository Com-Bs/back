package cache

import (
	"sync"
	"time"
)

// CompileCacheEntry represents a cached compile request and its response
type CompileCacheEntry struct {
	ResponseBody []byte
	StatusCode   int
	CreatedAt    time.Time
}

// CompileCache is a thread-safe cache for compile requests
type CompileCache struct {
	mu     sync.RWMutex
	items  map[string]*CompileCacheEntry
	maxAge time.Duration
}

// NewCompileCache creates a new compile cache with the specified max age
func NewCompileCache(maxAge time.Duration) *CompileCache {
	cache := &CompileCache{
		items:  make(map[string]*CompileCacheEntry),
		maxAge: maxAge,
	}

	// Start cleanup goroutine
	go cache.cleanup()

	return cache
}

// Get retrieves a cached response for the given request hash
func (c *CompileCache) Get(hash string) (*CompileCacheEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.items[hash]
	if !exists {
		return nil, false
	}

	// Check if entry is expired
	if time.Since(entry.CreatedAt) > c.maxAge {
		return nil, false
	}

	return entry, true
}

// Set stores a response in the cache
func (c *CompileCache) Set(hash string, responseBody []byte, statusCode int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[hash] = &CompileCacheEntry{
		ResponseBody: responseBody,
		StatusCode:   statusCode,
		CreatedAt:    time.Now(),
	}
}

// cleanup periodically removes expired entries
func (c *CompileCache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for hash, entry := range c.items {
			if now.Sub(entry.CreatedAt) > c.maxAge {
				delete(c.items, hash)
			}
		}
		c.mu.Unlock()
	}
}
