package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/minhhoccode111/go-clean-template-gin/internal/entity"
	"github.com/minhhoccode111/go-clean-template-gin/internal/usecase/translation"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var errInternalServErr = errors.New("internal server error")

type test struct {
	name string
	mock func()
	res  any
	err  error
}

func translationUseCase(t *testing.T) (*translation.UseCase, *MockTranslationRepo, *MockTranslationWebAPI, *MockTranslationCache) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := NewMockTranslationRepo(mockCtl)
	webAPI := NewMockTranslationWebAPI(mockCtl)
	cache := NewMockTranslationCache(mockCtl)

	useCase := translation.New(repo, webAPI, cache)

	return useCase, repo, webAPI, cache
}

func TestHistory(t *testing.T) { //nolint:tparallel // data races here
	t.Parallel()

	translationUseCase, repo, _, cache := translationUseCase(t)

	tests := []test{
		{
			name: "cache miss - empty result from db",
			mock: func() {
				cache.EXPECT().GetHistory(context.Background()).Return(nil, false)
				repo.EXPECT().GetHistory(context.Background()).Return(nil, nil)
				cache.EXPECT().SetHistory(context.Background(), []entity.Translation(nil)).Return(true)
			},
			res: entity.TranslationHistory{},
			err: nil,
		},
		{
			name: "cache hit",
			mock: func() {
				cached := []entity.Translation{{Original: "hello", Translation: "привет"}}
				cache.EXPECT().GetHistory(context.Background()).Return(cached, true)
			},
			res: entity.TranslationHistory{History: []entity.Translation{{Original: "hello", Translation: "привет"}}},
			err: nil,
		},
		{
			name: "cache miss - repo error",
			mock: func() {
				cache.EXPECT().GetHistory(context.Background()).Return(nil, false)
				repo.EXPECT().GetHistory(context.Background()).Return(nil, errInternalServErr)
			},
			res: entity.TranslationHistory{},
			err: errInternalServErr,
		},
	}

	for _, tc := range tests { //nolint:paralleltest // data races here
		localTc := tc

		t.Run(localTc.name, func(t *testing.T) {
			localTc.mock()

			res, err := translationUseCase.History(context.Background())

			require.Equal(t, res, localTc.res)
			require.ErrorIs(t, err, localTc.err)
		})
	}
}

func TestTranslate(t *testing.T) { //nolint:tparallel // data races here
	t.Parallel()

	translationUseCase, repo, webAPI, cache := translationUseCase(t)

	tests := []test{
		{
			name: "success - cache invalidated",
			mock: func() {
				webAPI.EXPECT().Translate(entity.Translation{}).Return(entity.Translation{}, nil)
				repo.EXPECT().Store(context.Background(), entity.Translation{}).Return(nil)
				cache.EXPECT().InvalidateHistory(context.Background())
			},
			res: entity.Translation{},
			err: nil,
		},
		{
			name: "web API error",
			mock: func() {
				webAPI.EXPECT().Translate(entity.Translation{}).Return(entity.Translation{}, errInternalServErr)
			},
			res: entity.Translation{},
			err: errInternalServErr,
		},
		{
			name: "repo error",
			mock: func() {
				webAPI.EXPECT().Translate(entity.Translation{}).Return(entity.Translation{}, nil)
				repo.EXPECT().Store(context.Background(), entity.Translation{}).Return(errInternalServErr)
			},
			res: entity.Translation{},
			err: errInternalServErr,
		},
	}

	for _, tc := range tests { //nolint:paralleltest // data races here
		localTc := tc

		t.Run(localTc.name, func(t *testing.T) {
			localTc.mock()

			res, err := translationUseCase.Translate(context.Background(), entity.Translation{})

			require.EqualValues(t, res, localTc.res)
			require.ErrorIs(t, err, localTc.err)
		})
	}
}
