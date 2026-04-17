//go:build migrate

package app

import (
	"context"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/minhhoccode111/go-clean-template-gin/internal/repo/persistent/ent"
)

const (
	_defaultAttempts = 20
	_defaultTimeout  = time.Second
)

func init() {
	databaseURL, ok := os.LookupEnv("PG_URL")
	if !ok || len(databaseURL) == 0 {
		log.Fatalf("migrate: environment variable not declared: PG_URL")
	}

	databaseURL += "?sslmode=disable"

	var (
		attempts = _defaultAttempts
		err      error
		client   *ent.Client
	)

	for attempts > 0 {
		client, err = ent.Open("pgx", databaseURL)
		if err == nil {
			break
		}

		log.Printf("Migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		log.Fatalf("Migrate: postgres connect error: %s", err)
	}

	defer client.Close()

	// Keep migrations additive by default; destructive changes must be explicit.
	err = client.Schema.Create(
		context.Background(),
		ent.WithDropColumn(false),
		ent.WithDropIndex(false),
	)
	if err != nil {
		log.Fatalf("Migrate: schema create error: %s", err)
	}

	log.Printf("Migrate: schema sync success")
}
