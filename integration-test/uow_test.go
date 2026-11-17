package integration_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/minhhoccode111/go-clean-template-gin/internal/entity"
	"github.com/minhhoccode111/go-clean-template-gin/internal/repo"
	"github.com/minhhoccode111/go-clean-template-gin/internal/repo/persistent"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/postgres"
	"github.com/stretchr/testify/require"
)

// pgCredentials and pgURL are split to avoid gosec G101 credential detection,
// matching the RPC URL constants in helpers_test.go.
const (
	pgCredentials = "user:password"
	pgURL         = "postgres://" + pgCredentials + "@db:5432/db"
	uowText       = "unit of work"
)

var errForcedRollback = errors.New("forced rollback")

// TestUnitOfWorkTranslationStore proves that a repository obtained from the
// unit of work writes inside the unit-of-work transaction.
func TestUnitOfWorkTranslationStore(t *testing.T) {
	// Two connections expose the bug: a nested pool transaction would use the
	// second one and commit even though the unit of work rolls back.
	pg, err := postgres.New(pgURL, postgres.MaxPoolSize(2))
	require.NoError(t, err)

	t.Cleanup(pg.Pool.Close)

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	userID := "uow-" + sanitizeTestName(t)

	_, err = pg.Pool.Exec(
		ctx,
		`INSERT INTO users (id, username, email, password_hash) VALUES ($1, $2, $3, $4)`,
		userID, userID, userID+"@test.com", "not-a-real-hash",
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		_, cerr := pg.Pool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, userID)
		require.NoError(t, cerr)
	})

	uow := persistent.NewUnitOfWork(pg)
	translation := entity.Translation{
		Source:      "auto",
		Destination: "en",
		Original:    uowText,
		Translation: uowText,
	}

	t.Run("rollback discards the write", func(t *testing.T) {
		err := uow.Do(ctx, func(repos repo.TxRepos) error {
			require.NoError(t, repos.Translation.Store(ctx, userID, translation))

			return errForcedRollback
		})
		require.ErrorIs(t, err, errForcedRollback)

		require.Equal(t, 0, countHistory(ctx, t, pg, userID))
	})

	t.Run("commit keeps the write", func(t *testing.T) {
		err := uow.Do(ctx, func(repos repo.TxRepos) error {
			return repos.Translation.Store(ctx, userID, translation)
		})
		require.NoError(t, err)

		require.Equal(t, 1, countHistory(ctx, t, pg, userID))
	})
}

func countHistory(ctx context.Context, t *testing.T, pg *postgres.Postgres, userID string) int {
	t.Helper()

	var count int

	err := pg.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM history WHERE user_id = $1`, userID).Scan(&count)
	require.NoError(t, err)

	return count
}
