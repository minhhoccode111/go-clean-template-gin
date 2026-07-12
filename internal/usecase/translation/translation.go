package translation

import (
	"context"
	"fmt"

	"github.com/minhhoccode111/go-clean-template-gin/internal/entity"
	"github.com/minhhoccode111/go-clean-template-gin/internal/repo"
	"github.com/minhhoccode111/go-clean-template-gin/internal/usecase"
)

// UseCase -.
type UseCase struct {
	repo   repo.TranslationRepo
	webAPI repo.TranslationWebAPI
	cache  repo.TranslationCache
	uow    repo.UnitOfWork
}

// New returns a Translation usecase instrumented with OpenTelemetry tracing spans.
func New(
	r repo.TranslationRepo,
	w repo.TranslationWebAPI,
	c repo.TranslationCache,
	uow repo.UnitOfWork,
) usecase.Translation {
	return newTraced(&UseCase{
		repo:   r,
		webAPI: w,
		cache:  c,
		uow:    uow,
	})
}

// History - getting translate history from store.
func (uc *UseCase) History(ctx context.Context, userID string) (entity.TranslationHistory, error) {
	if cached, ok := uc.cache.GetHistory(ctx); ok {
		return entity.TranslationHistory{History: cached}, nil
	}

	translations, err := uc.repo.GetHistory(ctx, userID)
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
	userID string,
	t entity.Translation,
) (entity.Translation, error) {
	translation, err := uc.webAPI.Translate(ctx, t)
	if err != nil {
		return entity.Translation{}, fmt.Errorf(
			"TranslationUseCase - Translate - s.webAPI.Translate: %w",
			err,
		)
	}

	err = uc.repo.Store(ctx, userID, translation)
	if err != nil {
		return entity.Translation{}, fmt.Errorf(
			"TranslationUseCase - Translate - s.repo.Store: %w",
			err,
		)
	}

	uc.cache.InvalidateHistory(ctx)

	return translation, nil
}
