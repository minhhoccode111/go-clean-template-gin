// Package repo implements application outer layer logic. Each logic group in own file.
package repo

import (
	"context"

	"github.com/minhhoccode111/go-clean-template-gin/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test

type (
	// TranslationRepo -.
	TranslationRepo interface {
		Store(context.Context, entity.Translation) error
		GetHistory(context.Context) ([]entity.Translation, error)
	}

	// TranslationWebAPI -.
	TranslationWebAPI interface {
		Translate(entity.Translation) (entity.Translation, error)
	}

	// TranslationCache is an in-process cache for translation data.
	// Implementations must be safe for concurrent use.
	TranslationCache interface {
		// GetHistory returns the cached translation history and whether it was found.
		GetHistory(ctx context.Context) ([]entity.Translation, bool)
		// SetHistory stores the translation history in the cache.
		// Returns false if the entry was not admitted by the cache.
		SetHistory(ctx context.Context, history []entity.Translation) bool
		// InvalidateHistory removes the cached history so the next read hits the DB.
		InvalidateHistory(ctx context.Context)
	}
)
