package cache

import (
	"context"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

const (
	_defaultTTL             = 10 * time.Minute
	_defaultCleanupInterval = 5 * time.Minute
	_defaultMaxItems        = 1000
)

// ICache defines the operations a cache must support.
type ICache interface {
	// Set adds a value to the cache with the default TTL.
	Set(ctx context.Context, key string, value any)

	// SetWithTTL adds a value to the cache with a custom TTL.
	SetWithTTL(ctx context.Context, key string, value any, ttl time.Duration)

	// Get retrieves a value from the cache.
	Get(ctx context.Context, key string) (any, bool)

	// Delete removes a value from the cache.
	Delete(ctx context.Context, key string)

	// Clear removes all values from the cache.
	Clear(ctx context.Context)

	// Size returns the number of items in the cache.
	Size() int64

	// Close stops all background tasks and releases resources.
	Close() error
}

// item represents a cached value with metadata.
type item struct {
	value      any
	expiration time.Time
	size       int // Approximate size in bytes
}

// Cache is a thread-safe in-memory cache with TTL and memory management.
type Cache struct {
	itemCount atomic.Int64 // Use atomic operations to track item count
	data      sync.Map

	defaultTTL      time.Duration
	cleanupInterval time.Duration
	maxItems        int
	onEviction      func(key string, value any)

	stopChan   chan struct{}
	closedChan chan struct{}
}

// New creates a new memory cache with the given configuration.
func New(opts ...Option) *Cache {
	c := &Cache{
		defaultTTL:      _defaultTTL,
		cleanupInterval: _defaultCleanupInterval,
		maxItems:        _defaultMaxItems,

		stopChan:   make(chan struct{}),
		closedChan: make(chan struct{}),
	}

	// Custom options
	for _, opt := range opts {
		opt(c)
	}

	go c.cleanupLoop()
	return c
}

// Set adds a value to the cache with the default TTL.
func (c *Cache) Set(ctx context.Context, key string, value any) {
	c.SetWithTTL(ctx, key, value, c.defaultTTL)
}

// SetWithTTL adds a value to the cache with a custom TTL.
func (c *Cache) SetWithTTL(_ context.Context, key string, value any, ttl time.Duration) {
	// Estimate size of the item (very rough approximation).
	size := estimateSize(value)

	// Check if the key already exists.
	// If it doesn't exist, we temporarily store a zero value to atomically
	// mark its presence and get a chance to increment itemCount.
	// LoadOrStore returns the actual value stored and a boolean indicating if the key was already present.
	_, loaded := c.data.LoadOrStore(key, item{}) // Placeholder item to mark key presence

	if !loaded { // Key was NOT present, so it's a new entry
		c.itemCount.Add(1)
	}
	// Now, store the actual item, overwriting the placeholder or existing item
	c.data.Store(key, item{
		value:      value,
		expiration: time.Now().Add(ttl),
		size:       size,
	})

	// If we're over the max items, clean up old items.
	if c.maxItems > 0 && c.itemCount.Load() > int64(c.maxItems) {
		c.cleanupOldest()
	}
}

// Get retrieves a value from the cache.
func (c *Cache) Get(_ context.Context, key string) (any, bool) {
	value, ok := c.data.Load(key)
	if !ok {
		return nil, false
	}

	itm, ok := value.(item)
	if !ok {
		// If the value is not of type item, it means it was corrupted or not set correctly.
		c.data.Delete(key)
		c.itemCount.Add(-1) // Decrement count for corrupted item
		return nil, false
	}
	if time.Now().After(itm.expiration) {
		c.data.Delete(key)
		c.itemCount.Add(-1)

		if c.onEviction != nil {
			c.onEviction(key, itm.value)
		}

		return nil, false
	}

	return itm.value, true
}

// Delete removes a value from the cache.
func (c *Cache) Delete(_ context.Context, key string) {
	if value, loaded := c.data.LoadAndDelete(key); loaded {
		c.itemCount.Add(-1)

		if c.onEviction != nil {
			if itm, ok := value.(item); ok {
				c.onEviction(key, itm.value)
			}
		}
	}
}

// Clear removes all values from the cache.
func (c *Cache) Clear(_ context.Context) {
	evicted := make(map[string]any) // Collect evicted items

	if c.onEviction != nil {
		c.data.Range(func(key, value any) bool {
			itm, ok := value.(item)
			if !ok {
				return true
			}
			if keyStr, ok := key.(string); ok {
				evicted[keyStr] = itm.value
			}
			return true
		})
	}

	c.data = sync.Map{}
	c.itemCount.Store(0)

	// Call eviction callbacks outside the loop
	if c.onEviction != nil {
		for k, v := range evicted {
			c.onEviction(k, v)
		}
	}
}

// Size returns the number of items in the cache.
func (c *Cache) Size() int64 {
	return c.itemCount.Load()
}

// Close stops the cache cleanup goroutine.
func (c *Cache) Close() error {
	select {
	case <-c.stopChan:
		// Already closed
		return nil
	default:
		close(c.stopChan)
		<-c.closedChan // Wait for cleanup goroutine to exit
		return nil
	}
}

// cleanupLoop periodically cleans up expired items.
func (c *Cache) cleanupLoop() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer func() {
		ticker.Stop()
		close(c.closedChan)
	}()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stopChan:
			return
		}
	}
}

// cleanup removes expired items.
func (c *Cache) cleanup() {
	evicted := make(map[string]any)
	count := 0

	c.data.Range(func(key, value any) bool {
		itm, ok := value.(item)
		if !ok {
			return true
		}
		if time.Now().After(itm.expiration) {
			c.data.Delete(key)
			count++

			if c.onEviction != nil {
				if keyStr, ok := key.(string); ok {
					evicted[keyStr] = itm.value
				}
			}
		}
		return true
	})

	if count > 0 {
		c.itemCount.Add(-int64(count))

		// Call eviction callbacks outside the loop to avoid blocking the range
		if c.onEviction != nil {
			for k, v := range evicted {
				c.onEviction(k, v)
			}
		}
	}
}

// cleanupOldest removes the oldest items if we're over the max items.
func (c *Cache) cleanupOldest() {
	currentCount := c.itemCount.Load()

	// If we're not over the max items, don't do anything
	if currentCount <= int64(c.maxItems) {
		return
	}

	// Remove 20% of max items at once, but at least 1
	threshold := max(c.maxItems/5, 1)

	// Collect all items into a temporary slice
	type keyExpPair struct {
		key        string
		value      any
		expiration time.Time
	}
	allCandidates := make([]keyExpPair, 0, currentCount)

	c.data.Range(func(key, value any) bool {
		itm, ok := value.(item)
		if !ok {
			return true // Skip corrupted items, they'll be cleaned by Get or cleanup()
		}
		if keyStr, ok := key.(string); ok {
			allCandidates = append(allCandidates, keyExpPair{keyStr, itm.value, itm.expiration})
		}
		return true
	})

	// Sort candidates by expiration time (oldest first)
	sort.Slice(allCandidates, func(i, j int) bool {
		return allCandidates[i].expiration.Before(allCandidates[j].expiration)
	})

	// Determine how many to delete, up to threshold
	toDeleteCount := min(threshold, len(allCandidates))
	if toDeleteCount == 0 {
		return // Nothing to delete
	}

	deletedCount := 0
	evicted := make(map[string]any)

	// Delete the oldest items
	for i := 0; i < toDeleteCount; i++ {
		candidate := allCandidates[i]
		if _, loaded := c.data.LoadAndDelete(candidate.key); loaded {
			deletedCount++
			if c.onEviction != nil {
				evicted[candidate.key] = candidate.value
			}
		}
	}

	// Update count
	if deletedCount > 0 {
		c.itemCount.Add(-int64(deletedCount))
	}

	// Call eviction callbacks outside the loop
	if c.onEviction != nil {
		for k, v := range evicted {
			c.onEviction(k, v)
		}
	}
}

// estimateSize attempts to estimate the memory footprint of a value.
func estimateSize(value any) int {
	switch v := value.(type) {
	case string:
		return len(v) + 24 // base size + string overhead
	case []byte:
		return len(v) + 24 // base size + slice overhead
	case map[string]any:
		return len(v) * 64 // rough estimate
	default:
		return 64 // default conservative estimate
	}
}
