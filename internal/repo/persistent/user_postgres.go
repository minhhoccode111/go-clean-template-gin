package persistent

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/minhhoccode111/go-clean-template-gin/internal/entity"
	"github.com/minhhoccode111/go-clean-template-gin/internal/repo"
	"github.com/minhhoccode111/go-clean-template-gin/internal/repo/persistent/sqlc"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/postgres"
)

// UserRepo implements repo.UserRepo using sqlc.
type UserRepo struct {
	*postgres.Postgres
	queries *sqlc.Queries
}

// NewUserRepo returns a User repository instrumented with OpenTelemetry tracing spans.
func NewUserRepo(pg *postgres.Postgres) repo.UserRepo {
	return newTracedUser(&UserRepo{
		Postgres: pg,
		queries:  sqlc.New(pg.Pool),
	})
}

// Store -.
func (r *UserRepo) Store(ctx context.Context, user *entity.User) error {
	err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    pgTimestamptz(user.CreatedAt),
		UpdatedAt:    pgTimestamptz(user.UpdatedAt),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return entity.ErrUserAlreadyExists
		}

		return fmt.Errorf("UserRepo - Store - r.queries.CreateUser: %w", err)
	}

	return nil
}

// GetByID -.
func (r *UserRepo) GetByID(ctx context.Context, id string) (entity.User, error) {
	row, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, entity.ErrUserNotFound
		}

		return entity.User{}, fmt.Errorf("UserRepo - GetByID - r.queries.GetUserByID: %w", err)
	}

	return userFromRow(&row), nil
}

// GetByEmail -.
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, entity.ErrUserNotFound
		}

		return entity.User{}, fmt.Errorf("UserRepo - GetByEmail - r.queries.GetUserByEmail: %w", err)
	}

	return userFromRow(&row), nil
}

func userFromRow(row *sqlc.User) entity.User {
	return entity.User{
		ID:           row.ID,
		Username:     row.Username,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}
}
