// Package cache implements the repo.TranslationCache interface using an in-process
// otter-backed cache. It is the only place in the codebase that knows about ottercache.
package cache

import (
	"context"

	"github.com/minhhoccode111/go-clean-template-gin/internal/entity"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/cache"
)

func historyKey(userID string) string {
	return "translation:history:" + userID
}

// TranslationCache wraps an ottercache.Cache and implements repo.TranslationCache.
type TranslationCache struct {
	c *cache.Cache[string, []entity.Translation]
}

// New creates a new TranslationCache.
// The TTL and max-cost are configured on the ottercache.Cache itself via ottercache.TTL / ottercache.MaxCost.
func New(c *cache.Cache[string, []entity.Translation]) *TranslationCache {
	return &TranslationCache{c: c}
}

// GetHistory returns the cached translation history for a user.
func (tc *TranslationCache) GetHistory(_ context.Context, userID string) ([]entity.Translation, bool) {
	return tc.c.Get(historyKey(userID))
}

// SetHistory stores the translation history for a user in the cache.
func (tc *TranslationCache) SetHistory(_ context.Context, userID string, history []entity.Translation) bool {
	return tc.c.Set(historyKey(userID), history)
}

// InvalidateHistory removes the user's cached history entry so the next read hits the DB.
func (tc *TranslationCache) InvalidateHistory(_ context.Context, userID string) {
	tc.c.Delete(historyKey(userID))
}
