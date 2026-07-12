// Package repo implements application outer layer logic. Each logic group in own file.
package repo

import (
	"context"

	"github.com/minhhoccode111/go-clean-template-gin/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test

type (
	// TxRepos holds transactional repository accessors that share the same pgx.Tx.
	TxRepos struct {
		Translation TranslationRepo
	}

	// UnitOfWork manages a database transaction lifecycle.
	UnitOfWork interface {
		Do(ctx context.Context, fn func(TxRepos) error) error
	}

	// TranslationRepo -.
	TranslationRepo interface {
		Store(ctx context.Context, userID string, t entity.Translation) error
		GetHistory(ctx context.Context, userID string) ([]entity.Translation, error)
	}

	// TranslationWebAPI -.
	TranslationWebAPI interface {
		Translate(context.Context, entity.Translation) (entity.Translation, error)
	}

	// TranslationCache is an in-process cache for translation data.
	TranslationCache interface {
		GetHistory(ctx context.Context) ([]entity.Translation, bool)
		SetHistory(ctx context.Context, history []entity.Translation) bool
		InvalidateHistory(ctx context.Context)
	}

	// UserRepo -.
	UserRepo interface {
		Store(ctx context.Context, user *entity.User) error
		GetByID(ctx context.Context, id string) (entity.User, error)
		GetByEmail(ctx context.Context, email string) (entity.User, error)
	}

	// TaskRepo -.
	TaskRepo interface {
		Store(ctx context.Context, task *entity.Task) error
		GetByID(ctx context.Context, userID, taskID string) (entity.Task, error)
		List(ctx context.Context, userID string, filter TaskFilter) ([]entity.Task, int, error)
		Update(ctx context.Context, task *entity.Task) error
		Delete(ctx context.Context, userID, taskID string) error
	}

	// TaskFilter -.
	TaskFilter struct {
		Status *entity.TaskStatus
		Limit  uint64
		Offset uint64
	}
)
