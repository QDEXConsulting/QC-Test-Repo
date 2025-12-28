package main

import (
	"sync"
	"time"
)

// CacheEntry represents a single cache entry
type CacheEntry struct {
	Value      interface{}
	ExpiresAt  time.Time
	CreatedAt  time.Time
	AccessCount int64
}

// IsExpired checks if the cache entry has expired
func (e *CacheEntry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// Cache provides an in-memory cache with TTL support
type Cache struct {
	entries map[string]*CacheEntry
	mu      sync.RWMutex
	ttl     time.Duration
	maxSize int
}

// NewCache creates a new cache instance
func NewCache(defaultTTL time.Duration, maxSize int) *Cache {
	cache := &Cache{
		entries: make(map[string]*CacheEntry),
		ttl:     defaultTTL,
		maxSize: maxSize,
	}
	
	// Start cleanup goroutine
	go cache.cleanup()
	
	return cache
}

// Get retrieves a value from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}
	
	if entry.IsExpired() {
		return nil, false
	}
	
	entry.AccessCount++
	return entry.Value, true
}

// Set stores a value in the cache with default TTL
func (c *Cache) Set(key string, value interface{}) {
	c.SetWithTTL(key, value, c.ttl)
}

// SetWithTTL stores a value in the cache with custom TTL
func (c *Cache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Check max size
	if c.maxSize > 0 && len(c.entries) >= c.maxSize {
		c.evictLRU()
	}
	
	c.entries[key] = &CacheEntry{
		Value:      value,
		ExpiresAt:  time.Now().Add(ttl),
		CreatedAt:  time.Now(),
		AccessCount: 0,
	}
}

// Delete removes a key from the cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Clear removes all entries from the cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]*CacheEntry)
}

// Size returns the number of entries in the cache
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}

// cleanup periodically removes expired entries
func (c *Cache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.entries {
			if entry.IsExpired() {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}

// evictLRU evicts the least recently used entry
func (c *Cache) evictLRU() {
	if len(c.entries) == 0 {
		return
	}
	
	var lruKey string
	var minAccess int64 = -1
	
	for key, entry := range c.entries {
		if minAccess == -1 || entry.AccessCount < minAccess {
			minAccess = entry.AccessCount
			lruKey = key
		}
	}
	
	if lruKey != "" {
		delete(c.entries, lruKey)
	}
}

// CacheStats provides statistics about the cache
type CacheStats struct {
	Size        int           `json:"size"`
	MaxSize     int           `json:"max_size"`
	DefaultTTL  time.Duration `json:"default_ttl"`
	HitRate     float64       `json:"hit_rate"`
	TotalHits   int64         `json:"total_hits"`
	TotalMisses int64         `json:"total_misses"`
}

// CacheWithStats extends Cache with statistics tracking
type CacheWithStats struct {
	*Cache
	stats struct {
		hits   int64
		misses int64
		mu     sync.RWMutex
	}
}

// NewCacheWithStats creates a new cache with statistics tracking
func NewCacheWithStats(defaultTTL time.Duration, maxSize int) *CacheWithStats {
	return &CacheWithStats{
		Cache: NewCache(defaultTTL, maxSize),
	}
}

// Get retrieves a value and updates statistics
func (c *CacheWithStats) Get(key string) (interface{}, bool) {
	value, found := c.Cache.Get(key)
	
	c.stats.mu.Lock()
	if found {
		c.stats.hits++
	} else {
		c.stats.misses++
	}
	c.stats.mu.Unlock()
	
	return value, found
}

// GetStats returns cache statistics
func (c *CacheWithStats) GetStats() CacheStats {
	c.stats.mu.RLock()
	hits := c.stats.hits
	misses := c.stats.misses
	c.stats.mu.RUnlock()
	
	total := hits + misses
	hitRate := 0.0
	if total > 0 {
		hitRate = float64(hits) / float64(total)
	}
	
	return CacheStats{
		Size:       c.Size(),
		MaxSize:    c.maxSize,
		DefaultTTL: c.ttl,
		HitRate:    hitRate,
		TotalHits:  hits,
		TotalMisses: misses,
	}
}

// ResetStats resets cache statistics
func (c *CacheWithStats) ResetStats() {
	c.stats.mu.Lock()
	c.stats.hits = 0
	c.stats.misses = 0
	c.stats.mu.Unlock()
}

// InventoryCache provides caching for inventory-related operations
type InventoryCache struct {
	itemsCache    *CacheWithStats
	statsCache    *CacheWithStats
	categoriesCache *CacheWithStats
}

// NewInventoryCache creates a new inventory cache
func NewInventoryCache() *InventoryCache {
	return &InventoryCache{
		itemsCache:    NewCacheWithStats(5*time.Minute, 1000),
		statsCache:    NewCacheWithStats(1*time.Minute, 100),
		categoriesCache: NewCacheWithStats(10*time.Minute, 100),
	}
}

// GetItem retrieves an item from cache
func (ic *InventoryCache) GetItem(itemID string) (*Item, bool) {
	value, found := ic.itemsCache.Get("item:" + itemID)
	if !found {
		return nil, false
	}
	
	item, ok := value.(*Item)
	return item, ok
}

// SetItem stores an item in cache
func (ic *InventoryCache) SetItem(item *Item) {
	ic.itemsCache.Set("item:"+item.ID, item)
}

// InvalidateItem removes an item from cache
func (ic *InventoryCache) InvalidateItem(itemID string) {
	ic.itemsCache.Delete("item:" + itemID)
	ic.statsCache.Clear() // Invalidate stats when items change
}

// GetStats retrieves inventory stats from cache
func (ic *InventoryCache) GetStats() (InventoryStats, bool) {
	value, found := ic.statsCache.Get("stats")
	if !found {
		return InventoryStats{}, false
	}
	
	stats, ok := value.(InventoryStats)
	return stats, ok
}

// SetStats stores inventory stats in cache
func (ic *InventoryCache) SetStats(stats InventoryStats) {
	ic.statsCache.Set("stats", stats)
}

// GetCategories retrieves categories from cache
func (ic *InventoryCache) GetCategories() ([]string, bool) {
	value, found := ic.categoriesCache.Get("categories")
	if !found {
		return nil, false
	}
	
	categories, ok := value.([]string)
	return categories, ok
}

// SetCategories stores categories in cache
func (ic *InventoryCache) SetCategories(categories []string) {
	ic.categoriesCache.Set("categories", categories)
}

// Clear clears all caches
func (ic *InventoryCache) Clear() {
	ic.itemsCache.Clear()
	ic.statsCache.Clear()
	ic.categoriesCache.Clear()
}

