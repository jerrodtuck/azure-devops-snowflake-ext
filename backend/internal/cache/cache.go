package cache

import (
	"sync"
	"time"
	"snowflake-dropdown-api/internal/models"
)

// Cache stores the dropdown data with expiration
type Cache struct {
	mu         sync.RWMutex
	data       map[string]models.DropdownResponse
	expiration time.Duration
}

// Global cache instance
var Instance = &Cache{
	data:       make(map[string]models.DropdownResponse),
	expiration: 1 * time.Hour, // Cache for 1 hour
}

// Get retrieves cached data
func (c *Cache) Get(key string) (models.DropdownResponse, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, exists := c.data[key]
	return val, exists
}

// Set stores data in cache
func (c *Cache) Set(key string, value models.DropdownResponse) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = cacheItem{
        value:     value,
        expiresAt: time.Now().Add(c.expiration),
    }

    // Start cleanup goroutine once
    c.cleanupOnce.Do(func() {
        go c.cleanupExpired()
    })
}

// Clear removes all cached data
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]models.DropdownResponse)
}

// SetExpiration sets the cache expiration duration
func (c *Cache) SetExpiration(duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.expiration = duration
}