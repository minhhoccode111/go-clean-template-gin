package cache

import "time"

// Option -.
type Option func(*Cache)

// DefaultTTL -.
func DefaultTTL(ttl time.Duration) Option {
	return func(c *Cache) {
		c.defaultTTL = ttl
	}
}

// CleanupInterval -.
func CleanupInterval(interval time.Duration) Option {
	return func(c *Cache) {
		c.cleanupInterval = interval
	}
}

// MaxItems -.
func MaxItems(max int) Option {
	return func(c *Cache) {
		c.maxItems = max
	}
}

// OnEviction -.
func OnEviction(fn func(key string, value any)) Option {
	return func(c *Cache) {
		c.onEviction = fn
	}
}
