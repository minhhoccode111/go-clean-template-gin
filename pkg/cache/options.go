package cache

import "time"

// Option configures a Cache instance.
type Option func(*options)

type options struct {
	maxCost int
	ttl     time.Duration
}

// MaxCost sets the maximum number of entries the cache can hold.
func MaxCost(size int) Option {
	return func(o *options) {
		o.maxCost = size
	}
}

// TTL sets how long each entry lives before eviction.
func TTL(ttl time.Duration) Option {
	return func(o *options) {
		o.ttl = ttl
	}
}
