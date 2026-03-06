// Package cache implements the repo.TranslationCache interface using an in-process
// otter-backed cache. It is the only place in the codebase that knows about ottercache.
package cache

import (
	"context"

	"github.com/minhhoccode111/go-clean-template-gin/internal/entity"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/cache"
)

const historyKey = "translation:history"

// TranslationCache wraps an ottercache.Cache and implements repo.TranslationCache.
type TranslationCache struct {
	c *cache.Cache[string, []entity.Translation]
}

// New creates a new TranslationCache.
// The TTL and max-cost are configured on the ottercache.Cache itself via ottercache.TTL / ottercache.MaxCost.
func New(c *cache.Cache[string, []entity.Translation]) *TranslationCache {
	return &TranslationCache{c: c}
}

// GetHistory returns the cached translation history.
func (tc *TranslationCache) GetHistory(_ context.Context) ([]entity.Translation, bool) {
	return tc.c.Get(historyKey)
}

// SetHistory stores the translation history in the cache.
func (tc *TranslationCache) SetHistory(_ context.Context, history []entity.Translation) bool {
	return tc.c.Set(historyKey, history)
}

// InvalidateHistory removes the cached history entry so the next read hits the DB.
func (tc *TranslationCache) InvalidateHistory(_ context.Context) {
	tc.c.Delete(historyKey)
}
