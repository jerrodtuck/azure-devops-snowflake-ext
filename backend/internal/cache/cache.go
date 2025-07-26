package cache

import (
	"snowflake-dropdown-api/internal/models"
	"sync"
	"time"
)

// cacheItem represents a cached item with expiration
type cacheItem struct {
	value     models.DropdownResponse
	expiresAt time.Time
}

// Cache stores the dropdown data with expiration
type Cache struct {
	mu          sync.RWMutex
	data        map[string]cacheItem
	expiration  time.Duration
	cleanupOnce sync.Once
}

// Global cache instance
var Instance = &Cache{
	data:       make(map[string]cacheItem),
	expiration: 1 * time.Hour, // Cache for 1 hour
}

// Get retrieves cached data
func (c *Cache) Get(key string) (models.DropdownResponse, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return models.DropdownResponse{}, false
	}

	// Check if expired
	if time.Now().After(item.expiresAt) {
		// Item expired, remove it
		c.mu.RUnlock()
		c.mu.Lock()
		delete(c.data, key)
		c.mu.Unlock()
		c.mu.RLock()
		return models.DropdownResponse{}, false
	}

	return item.value, true
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
	c.data = make(map[string]cacheItem)
}

// SetExpiration sets the cache expiration duration
func (c *Cache) SetExpiration(duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.expiration = duration
}

// cleanupExpired removes expired items from cache periodically
func (c *Cache) cleanupExpired() {
	ticker := time.NewTicker(30 * time.Minute) // Cleanup every 30 minutes
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.data {
			if now.After(item.expiresAt) {
				delete(c.data, key)
			}
		}
		c.mu.Unlock()
	}
}
