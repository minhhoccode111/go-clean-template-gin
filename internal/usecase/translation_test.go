package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/minhhoccode111/go-clean-template-gin/internal/entity"
	"github.com/minhhoccode111/go-clean-template-gin/internal/usecase"
	"github.com/minhhoccode111/go-clean-template-gin/internal/usecase/translation"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var errInternalServErr = errors.New("internal server error")

const (
	testUserID         = "test-user-123"
	testOriginalText   = "hello"
	testTranslatedText = "привет"
)

type test struct {
	name string
	mock func(repo *MockTranslationRepo, webAPI *MockTranslationWebAPI, cache *MockTranslationCache)
	res  any
	err  error
}

func translationUseCase(
	t *testing.T,
) (usecase.Translation, *MockTranslationRepo, *MockTranslationWebAPI, *MockTranslationCache) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	t.Cleanup(mockCtl.Finish)

	repo := NewMockTranslationRepo(mockCtl)
	webAPI := NewMockTranslationWebAPI(mockCtl)
	cache := NewMockTranslationCache(mockCtl)

	useCase := translation.New(repo, webAPI, cache, nil)

	return useCase, repo, webAPI, cache
}

func TestHistory(t *testing.T) {
	t.Parallel()

	tests := []test{
		{
			name: "cache miss - empty result from db",
			mock: func(repo *MockTranslationRepo, _ *MockTranslationWebAPI, cache *MockTranslationCache) {
				cache.EXPECT().GetHistory(gomock.Any()).Return(nil, false)
				repo.EXPECT().GetHistory(gomock.Any(), testUserID).Return(nil, nil)
				cache.EXPECT().SetHistory(gomock.Any(), []entity.Translation(nil)).Return(true)
			},
			res: entity.TranslationHistory{},
			err: nil,
		},
		{
			name: "cache hit",
			mock: func(_ *MockTranslationRepo, _ *MockTranslationWebAPI, cache *MockTranslationCache) {
				cached := []entity.Translation{
					{Original: testOriginalText, Translation: testTranslatedText},
				}
				cache.EXPECT().GetHistory(gomock.Any()).Return(cached, true)
			},
			res: entity.TranslationHistory{
				History: []entity.Translation{
					{Original: testOriginalText, Translation: testTranslatedText},
				},
			},
			err: nil,
		},
		{
			name: "cache miss - repo error",
			mock: func(repo *MockTranslationRepo, _ *MockTranslationWebAPI, cache *MockTranslationCache) {
				cache.EXPECT().GetHistory(gomock.Any()).Return(nil, false)
				repo.EXPECT().GetHistory(gomock.Any(), testUserID).Return(nil, errInternalServErr)
			},
			res: entity.TranslationHistory{},
			err: errInternalServErr,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, repo, webAPI, cache := translationUseCase(t)
			tc.mock(repo, webAPI, cache)

			res, err := useCase.History(context.Background(), testUserID)

			require.Equal(t, res, tc.res)
			require.ErrorIs(t, err, tc.err)
		})
	}
}

func TestTranslate(t *testing.T) {
	t.Parallel()

	tests := []test{
		{
			name: "success - cache invalidated",
			mock: func(repo *MockTranslationRepo, webAPI *MockTranslationWebAPI, cache *MockTranslationCache) {
				webAPI.EXPECT().
					Translate(gomock.Any(), entity.Translation{}).
					Return(entity.Translation{}, nil)
				repo.EXPECT().Store(gomock.Any(), testUserID, entity.Translation{}).Return(nil)
				cache.EXPECT().InvalidateHistory(gomock.Any())
			},
			res: entity.Translation{},
			err: nil,
		},
		{
			name: "web API error",
			mock: func(_ *MockTranslationRepo, webAPI *MockTranslationWebAPI, _ *MockTranslationCache) {
				webAPI.EXPECT().
					Translate(gomock.Any(), entity.Translation{}).
					Return(entity.Translation{}, errInternalServErr)
			},
			res: entity.Translation{},
			err: errInternalServErr,
		},
		{
			name: "repo error",
			mock: func(repo *MockTranslationRepo, webAPI *MockTranslationWebAPI, _ *MockTranslationCache) {
				webAPI.EXPECT().
					Translate(gomock.Any(), entity.Translation{}).
					Return(entity.Translation{}, nil)
				repo.EXPECT().
					Store(gomock.Any(), testUserID, entity.Translation{}).
					Return(errInternalServErr)
			},
			res: entity.Translation{},
			err: errInternalServErr,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, repo, webAPI, cache := translationUseCase(t)
			tc.mock(repo, webAPI, cache)

			res, err := useCase.Translate(context.Background(), testUserID, entity.Translation{})

			require.EqualValues(t, res, tc.res)
			require.ErrorIs(t, err, tc.err)
		})
	}
}
