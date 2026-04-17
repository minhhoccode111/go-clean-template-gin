//go:build migrate

package app

import (
	"log"
	"os"
	"os/exec"
	"time"
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

	attempts := _defaultAttempts
	var err error

	atlasBin := "atlas"
	if _, lookErr := exec.LookPath(atlasBin); lookErr != nil {
		atlasBin = "/atlas"
	}

	for attempts > 0 {
		cmd := exec.Command( //nolint:gosec // arguments are constants/environment values
			atlasBin,
			"migrate",
			"apply",
			"--url",
			databaseURL,
			"--dir",
			"file://migrations",
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err == nil {
			log.Printf("Migrate: apply success")
			return
		}

		log.Printf("Migrate: atlas apply failed, attempts left: %d", attempts)
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		log.Fatalf("Migrate: atlas apply error: %s", err)
	}
}
