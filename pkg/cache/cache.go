// Package ottercache provides a thin wrapper around github.com/maypok86/otter/v2,
// suitable for use as an injectable infrastructure dependency.
package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/maypok86/otter/v2"
)

const (
	_defaultMaxCost = 10_000
	_defaultTTL     = 5 * time.Minute
)

// Cache is a generic, in-process cache with variable per-entry TTL.
// K must be comparable (e.g. string); V can be any type.
type Cache[K comparable, V any] struct {
	c   *otter.Cache[K, V]
	ttl time.Duration
}

// New creates a cache with the given options.
// Defaults: MaximumSize=10_000, TTL=5m.
func New[K comparable, V any](opts ...Option) (*Cache[K, V], error) {
	o := &options{
		maxCost: _defaultMaxCost,
		ttl:     _defaultTTL,
	}

	for _, opt := range opts {
		opt(o)
	}

	c, err := otter.New(&otter.Options[K, V]{
		MaximumSize:      o.maxCost,
		ExpiryCalculator: otter.ExpiryWriting[K, V](o.ttl),
	})
	if err != nil {
		return nil, fmt.Errorf("ottercache: build: %w", err)
	}

	return &Cache[K, V]{c: c, ttl: o.ttl}, nil
}

// Get returns the value stored for key and whether it was found.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	return c.c.GetIfPresent(key)
}

// Set stores value under key with the cache's configured TTL.
// Returns false if the entry was not admitted (e.g. cost > capacity).
func (c *Cache[K, V]) Set(key K, value V) bool {
	_, ok := c.c.Set(key, value)

	return ok
}

// Delete removes the entry for key, if it exists.
func (c *Cache[K, V]) Delete(key K) {
	c.c.Invalidate(key)
}

// GetOrLoad returns the cached value for key.
// On a cache miss it calls load, stores the result with the cache's default TTL,
// and returns it. Concurrent calls for the same key are deduplicated — only one
// loader is invoked and the rest wait for that result (singleflight semantics).
func (c *Cache[K, V]) GetOrLoad(
	ctx context.Context,
	key K,
	load func(context.Context, K) (V, error),
) (V, error) {
	return c.c.Get(ctx, key, otter.LoaderFunc[K, V](load))
}

// SetWithTTL stores value under key with an explicit TTL, overriding the
// cache-wide default for this entry only.
// Returns false if the entry was not admitted (e.g. cost > capacity).
func (c *Cache[K, V]) SetWithTTL(key K, value V, ttl time.Duration) bool {
	_, ok := c.c.Set(key, value)
	if ok {
		c.c.SetExpiresAfter(key, ttl)
	}

	return ok
}
