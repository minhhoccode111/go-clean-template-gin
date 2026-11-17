# Upstream Parity — Phase 1 + Phase 2 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Port upstream's non-breaking parity features and per-user translation history while preserving this fork's Gin/sqlc architecture.

**Architecture:** Keep all existing Gin/sqlc boundaries. Upgrade toolchain and tracing in place, add Jaeger to the Docker stack, port agent/docs, then replace the global history table/cache with per-user equivalents through a new migration and regenerated sqlc code.

**Tech Stack:** Go 1.27, Gin, sqlc, pgx, OTel/otelpgx, Jaeger, Otter cache, golangci-lint.

**Spec:** `docs/superpowers/specs/2026-11-16-upstream-parity-design.md`

## Global Constraints

- Module path stays `github.com/minhhoccode111/go-clean-template-gin`.
- IDs stay `VARCHAR`; timestamps stay `TIMESTAMPTZ`.
- Keep the existing REST/gRPC/AMQP/NATS route surface and response shapes.
- No Fiber, Squirrel, prefork, UUID migration, README_CN/RU, or JWT refactor.
- Destructive `history` migration is approved.
- Every task ends with build/test/lint verification and a commit.

## Review Focus

- A second user must never receive another user's cached or stored translation history.
- The destructive migration must leave a usable `history` table with a valid `users` foreign key.
- Existing history rows are intentionally discarded; rollback recreates the old table shape.
- Tracing must remain no-op when disabled and must not break integration tests when Jaeger is absent.
- Go 1.27 must be used consistently by local `go.mod`, CI, and the integration Docker image.

---

### Task 1: Toolchain and lint config

**Files:**
- Modify: `go.mod`
- Modify: `.github/workflows/ci.yml`
- Modify: `integration-test/Dockerfile`
- Modify: `.golangci.yml`

**Interfaces:**
- Produces: Go 1.27 toolchain across local, CI, and integration environments; clean linter config.

- [ ] **Step 1: Update Go directive and CI/Docker versions**

Set `go 1.27.0` in `go.mod`, replace all `go-version: '1.26'` / `go-version: 1.26` entries in
`.github/workflows/ci.yml` with `1.27`, and change both `golang:1.26-alpine3.23` images in
`integration-test/Dockerfile` to `golang:1.27-alpine3.23`.

- [ ] **Step 2: Update the gofumpt setting**

Replace:

```yaml
    gofumpt:
      extra-rules: true
```

with:

```yaml
    gofumpt:
      extra:
        group-params: true
```

- [ ] **Step 3: Verify**

Run: `go mod tidy && go build ./... && golangci-lint run --timeout=15m`
Expected: build succeeds, `0 issues.`

- [ ] **Step 4: Commit**

```bash
git add go.mod go.sum .github/workflows/ci.yml integration-test/Dockerfile .golangci.yml
git commit -m "chore: upgrade toolchain to go 1.27"
```

---

### Task 2: PostgreSQL tracing and Jaeger stack

**Files:**
- Modify: `go.mod`, `go.sum`
- Modify: `pkg/postgres/postgres.go`
- Modify: `docker-compose.yml`
- Modify: `.env.example`
- Modify: `README.md`

**Interfaces:**
- Consumes: `github.com/exaring/otelpgx v0.13.0`
- Produces: pgx pool instrumented with OTel; Jaeger reachable at `jaeger:4317` in compose.

- [ ] **Step 1: Add the dependency**

Run: `go get github.com/exaring/otelpgx@v0.13.0 && go mod tidy`
Expected: `github.com/exaring/otelpgx` appears as a direct requirement.

- [ ] **Step 2: Instrument the pgx pool**

In `pkg/postgres/postgres.go`, import `github.com/exaring/otelpgx` and set:

```go
poolConfig.ConnConfig.Tracer = otelpgx.NewTracer()
```

immediately after `poolConfig.MaxConns` is assigned.

- [ ] **Step 3: Add Jaeger and tracing environment**

Add this service to `docker-compose.yml`:

```yaml
  jaeger:
    container_name: jaeger
    image: jaegertracing/all-in-one:1.76.0
    ports:
      - "16686:16686"
      - "4317:4317"
      - "4318:4318"
      - "6831:6831/udp"
    networks:
      app_network:
        aliases:
          - jaeger.lvh.me
```

Add `jaeger` to `app.depends_on`. Add tracing env values to the `x-backend-app-environment`
anchor and `.env.example`:

```yaml
  TRACING_ENABLED: "true"
  TRACING_OTLP_ENDPOINT: "jaeger:4317"
  TRACING_OTLP_INSECURE: "true"
  TRACING_SAMPLE_RATE: "1.0"
```

- [ ] **Step 4: Document tracing**

Add a README `Observability` section covering `TRACING_ENABLED`, `TRACING_OTLP_ENDPOINT`,
`TRACING_OTLP_INSECURE`, `TRACING_SAMPLE_RATE`, and Jaeger UI at `http://jaeger.lvh.me`.

- [ ] **Step 5: Verify**

Run: `go build ./... && make test`
Expected: exit 0.

- [ ] **Step 6: Commit**

```bash
git add go.mod go.sum pkg/postgres/postgres.go docker-compose.yml .env.example README.md
git commit -m "feat(tracing): instrument postgres and add jaeger stack"
```

---

### Task 3: Agent documentation and README parity

**Files:**
- Create: `CLAUDE.md`
- Create: `AGENTS.md` (symlink to `CLAUDE.md`)
- Modify: `README.md`

**Interfaces:**
- Produces: repository-specific agent instructions and user-facing domain/observability docs.

- [ ] **Step 1: Write `CLAUDE.md`**

Adapt upstream's `CLAUDE.md` to this repository. Cover: Makefile-first commands, dependency rule,
Gin/sqlc boundaries, cache/validatorx/unit-of-work ownership, generated files (`sqlc`, Swagger,
protobuf, mocks, `api.ts`), migrations, testing, API versioning, and adding a domain.

- [ ] **Step 2: Create the `AGENTS.md` symlink**

Run: `ln -s CLAUDE.md AGENTS.md`
Verify: `readlink AGENTS.md` prints `CLAUDE.md`.

- [ ] **Step 3: Add README domain and observability sections**

Add a concise table of the three domains and their REST/gRPC routes, plus the tracing section from
Task 2. Keep existing Gin/sqlc/cache overview text.

- [ ] **Step 4: Verify**

Run: `git diff --check && test -L AGENTS.md && test -f CLAUDE.md`
Expected: no whitespace errors; symlink exists.

- [ ] **Step 5: Commit**

```bash
git add CLAUDE.md AGENTS.md README.md
git commit -m "docs: add agent guide and upstream domain/observability docs"
```

---

### Task 4: Destructive per-user history migration and sqlc repository

**Files:**
- Create: `migrations/20261116000001_add_user_id_to_history.up.sql`
- Create: `migrations/20261116000001_add_user_id_to_history.down.sql`
- Modify: `internal/repo/persistent/queries/translation.sql`
- Regenerate: `internal/repo/persistent/sqlc/translation.sql.go`
- Modify: `internal/repo/persistent/translation/translation.go`
- Regenerate: `docs/schema.sql` if the database is reachable

**Interfaces:**
- Produces: `sqlc.GetHistoryParams{UserID string}`, `sqlc.StoreParams{UserID string, ...}`;
  `translation.Repo.GetHistory(ctx, userID)` filters by user and `Store` writes `user_id`.

- [ ] **Step 1: Write the destructive migration**

`up.sql`:

```sql
DROP TABLE IF EXISTS history;

CREATE TABLE IF NOT EXISTS history (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    source VARCHAR(255),
    destination VARCHAR(255),
    original VARCHAR(255),
    translation VARCHAR(255)
);

CREATE INDEX IF NOT EXISTS idx_history_user_id ON history(user_id);
```

`down.sql`:

```sql
DROP TABLE IF EXISTS history;

CREATE TABLE IF NOT EXISTS history(
    id serial PRIMARY KEY,
    source VARCHAR(255) NOT NULL,
    destination VARCHAR(255) NOT NULL,
    original VARCHAR(255) NOT NULL,
    translation VARCHAR(255) NOT NULL
);
```

- [ ] **Step 2: Update sqlc queries**

`internal/repo/persistent/queries/translation.sql`:

```sql
-- name: GetHistory :many
SELECT source, destination, original, translation FROM history WHERE user_id = $1;

-- name: Store :exec
INSERT INTO history (user_id, source, destination, original, translation)
VALUES ($1, $2, $3, $4, $5);
```

- [ ] **Step 3: Regenerate sqlc code**

Run: `sqlc generate`
Expected: `GetHistory(ctx, userID string)` and `StoreParams.UserID` appear in generated code.

- [ ] **Step 4: Update the repository**

In `internal/repo/persistent/translation/translation.go`:

- `GetHistory` calls `r.queries.GetHistory(ctx, userID)` and no longer ignores its parameter.
- `Store` passes `UserID: userID` into `sqlc.StoreParams`.
- Keep the existing transaction wrapper and tracing behavior.

- [ ] **Step 5: Verify**

Run: `go build ./... && go test ./internal/repo/...`
Expected: exit 0.

- [ ] **Step 6: Commit**

```bash
git add migrations/20261116000001_add_user_id_to_history.*.sql \
  internal/repo/persistent/queries/translation.sql \
  internal/repo/persistent/sqlc/translation.sql.go \
  internal/repo/persistent/translation/translation.go docs/schema.sql
git commit -m "feat(db): scope translation history by user"
```

---

### Task 5: Per-user cache and usecase wiring

**Files:**
- Modify: `internal/repo/contracts.go`
- Modify: `internal/repo/cache/translation_cache.go`
- Modify: `internal/usecase/translation/translation.go`
- Modify: `internal/usecase/translation_test.go`
- Regenerate: `internal/usecase/mocks_repo_test.go`

**Interfaces:**
- Produces: `TranslationCache.GetHistory(ctx, userID)`, `SetHistory(ctx, userID, history)`,
  `InvalidateHistory(ctx, userID)`; cache key `translation:history:<userID>`.

- [ ] **Step 1: Update the cache contract**

Change the `TranslationCache` interface methods to accept `userID string`.

- [ ] **Step 2: Update the Otter cache implementation**

Use `historyKey(userID)` returning `"translation:history:" + userID` for get/set/invalidate.

- [ ] **Step 3: Update the usecase and unit tests**

Pass `userID` into every cache call in `History` and `Translate`. Update the gomock expectations in
`internal/usecase/translation_test.go` to include `testUserID`.

- [ ] **Step 4: Regenerate mocks**

Run: `make mock`
Expected: `mocks_repo_test.go` compiles with the new cache interface.

- [ ] **Step 5: Verify**

Run: `make test`
Expected: all unit tests pass.

- [ ] **Step 6: Commit**

```bash
git add internal/repo/contracts.go internal/repo/cache/translation_cache.go \
  internal/usecase/translation/translation.go internal/usecase/translation_test.go \
  internal/usecase/mocks_repo_test.go
git commit -m "feat(cache): scope translation history cache by user"
```

---

### Task 6: Integration coverage and final verification

**Files:**
- Modify: `integration-test/translation_test.go`

**Interfaces:**
- Consumes: HTTP auth/translation helpers from `integration-test/helpers_test.go`.
- Produces: `TestHTTPHistoryUserIsolation` proving per-user isolation.

- [ ] **Step 1: Add the failing isolation test**

Add a test that:
1. Registers/logs in user A with `registerAndLogin(t)`.
2. Registers/logs in user B with a unique `_second` name/email.
3. User A translates a unique marker text.
4. Asserts A's history contains the marker and B's history does not.
5. User B translates a different marker.
6. Asserts B's history contains only B's marker and A's history does not contain B's marker.

- [ ] **Step 2: Run it before the implementation**

Run: `make compose-up-integration-test`
Expected before Tasks 4/5: FAIL (global history/cache). If Tasks 4/5 are already applied, this is
the acceptance test and must PASS.

- [ ] **Step 3: Run the complete verification gate**

Run:

```bash
go mod tidy -diff
go build ./...
make test
make linter-golangci
make compose-up-integration-test
```

Expected: no tidy diff, build exit 0, unit tests pass, `0 issues`, integration suite PASS.

- [ ] **Step 4: Update the sync spec**

Append a short "Latest sync" note to
`docs/superpowers/specs/2026-08-12-upstream-sync-design.md` pointing at the new parity spec.

- [ ] **Step 5: Commit**

```bash
git add integration-test/translation_test.go \
  docs/superpowers/specs/2026-08-12-upstream-sync-design.md
git commit -m "test(integration): prove per-user translation history isolation"
```
