package translation

import (
	"context"
	"fmt"

	"github.com/minhhoccode111/go-clean-template-gin/internal/entity"
	"github.com/minhhoccode111/go-clean-template-gin/internal/repo"
)

// UseCase -.
type UseCase struct {
	repo   repo.TranslationRepo
	webAPI repo.TranslationWebAPI
	cache  repo.TranslationCache
}

// New -.
func New(r repo.TranslationRepo, w repo.TranslationWebAPI, c repo.TranslationCache) *UseCase {
	return &UseCase{
		repo:   r,
		webAPI: w,
		cache:  c,
	}
}

// History - getting translate history from store.
// Uses cache-aside: returns cached results when available, falls back to DB otherwise.
func (uc *UseCase) History(ctx context.Context) (entity.TranslationHistory, error) {
	if cached, ok := uc.cache.GetHistory(ctx); ok {
		return entity.TranslationHistory{History: cached}, nil
	}

	translations, err := uc.repo.GetHistory(ctx)
	if err != nil {
		return entity.TranslationHistory{}, fmt.Errorf(
			"TranslationUseCase - History - s.repo.GetHistory: %w",
			err,
		)
	}

	uc.cache.SetHistory(ctx, translations)

	return entity.TranslationHistory{History: translations}, nil
}

// Translate -.
func (uc *UseCase) Translate(
	ctx context.Context,
	t entity.Translation,
) (entity.Translation, error) {
	translation, err := uc.webAPI.Translate(t)
	if err != nil {
		return entity.Translation{}, fmt.Errorf(
			"TranslationUseCase - Translate - s.webAPI.Translate: %w",
			err,
		)
	}

	err = uc.repo.Store(ctx, translation)
	if err != nil {
		return entity.Translation{}, fmt.Errorf(
			"TranslationUseCase - Translate - s.repo.Store: %w",
			err,
		)
	}

	// Invalidate the history cache so the next History call hits the DB.
	uc.cache.InvalidateHistory(ctx)

	return translation, nil
}
