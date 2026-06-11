package persistent

import (
	"context"
	"fmt"

	"github.com/minhhoccode111/go-clean-template-gin/internal/repo"
	"github.com/minhhoccode111/go-clean-template-gin/internal/repo/persistent/sqlc"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/postgres"
)

// PgUnitOfWork implements repo.UnitOfWork using pgx transactions.
type PgUnitOfWork struct {
	pg *postgres.Postgres
}

// NewUnitOfWork creates a new unit of work.
func NewUnitOfWork(pg *postgres.Postgres) *PgUnitOfWork {
	return &PgUnitOfWork{pg: pg}
}

// Do executes fn within a database transaction.
// Commits on success, rolls back on error.
func (u *PgUnitOfWork) Do(ctx context.Context, fn func(repo.TxRepos) error) error {
	tx, err := u.pg.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("PgUnitOfWork - Do - u.pg.Pool.Begin: %w", err)
	}

	defer tx.Rollback(ctx) //nolint:errcheck // rollback is safe to call after commit

	q := sqlc.New(tx)

	txRepos := repo.TxRepos{
		Translation: &TranslationRepo{Postgres: u.pg, queries: q},
	}

	if err := fn(txRepos); err != nil {
		return err // rollback via defer
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("PgUnitOfWork - Do - tx.Commit: %w", err)
	}

	return nil
}
