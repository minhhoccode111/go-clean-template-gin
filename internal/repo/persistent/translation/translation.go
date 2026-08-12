package translation

import (
	"context"
	"fmt"

	"github.com/minhhoccode111/go-clean-template-gin/internal/entity"
	"github.com/minhhoccode111/go-clean-template-gin/internal/repo"
	"github.com/minhhoccode111/go-clean-template-gin/internal/repo/persistent/sqlc"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/postgres"
)

// Repo implements repo.TranslationRepo using sqlc.
type Repo struct {
	*postgres.Postgres
	queries *sqlc.Queries
}

// New returns a Translation repository instrumented with OpenTelemetry tracing spans.
func New(pg *postgres.Postgres) repo.TranslationRepo {
	return newTraced(&Repo{
		Postgres: pg,
		queries:  sqlc.New(pg.Pool),
	})
}

// NewTx returns a raw Translation repository bound to the given queries
// (e.g. transaction-scoped), used by the unit of work.
func NewTx(pg *postgres.Postgres, q *sqlc.Queries) repo.TranslationRepo {
	return &Repo{
		Postgres: pg,
		queries:  q,
	}
}

// GetHistory -.
func (r *Repo) GetHistory(ctx context.Context, _ string) ([]entity.Translation, error) {
	rows, err := r.queries.GetHistory(ctx)
	if err != nil {
		return nil, fmt.Errorf("TranslationRepo - GetHistory - r.queries.GetHistory: %w", err)
	}

	entities := make([]entity.Translation, 0, len(rows))

	for _, row := range rows {
		entities = append(entities, entity.Translation{
			Source:      row.Source,
			Destination: row.Destination,
			Original:    row.Original,
			Translation: row.Translation,
		})
	}

	return entities, nil
}

// Store -.
func (r *Repo) Store(ctx context.Context, _ string, t entity.Translation) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("TranslationRepo - Store - r.Pool.Begin: %w", err)
	}

	defer tx.Rollback(ctx) //nolint:errcheck // rollback error is expected; commit decides the outcome

	queriesTx := r.queries.WithTx(tx)

	err = queriesTx.Store(ctx, sqlc.StoreParams{
		Source:      t.Source,
		Destination: t.Destination,
		Original:    t.Original,
		Translation: t.Translation,
	})
	if err != nil {
		return fmt.Errorf("TranslationRepo - Store - qtx.Store: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("TranslationRepo - Store - tx.Commit: %w", err)
	}

	return nil
}
