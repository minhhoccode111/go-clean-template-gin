# Upstream Sync Port: evrone/go-clean-template → go-clean-template-gin

Date: 2026-08-12
Status: Approved (design), pending implementation plan

## Goal

Sync the current upstream state of [evrone/go-clean-template](https://github.com/evrone/go-clean-template)
(`upstream/master` @ `376bcd4`) into this project (go-clean-template-gin), adopting upstream's
structural direction without losing gin-specific changes. This is the latest iteration of the
user's established sync pattern (previous: commit `74b78ae` "pull upstream features").

## Repository relationship (verified)

- The gin project has **no shared git history** with upstream — it began as a snapshot
  (2026-02-26) of the fiber template and evolved independently. Plain `git merge` is impossible;
  sync is manual porting.
- The repos **co-evolve**: upstream ported the gin project's user/task/JWT features to fiber
  (`2860e4c`, 2026-04-04 — line-identical files modulo framework), then restructured beyond
  what gin pulled on 2026-07-12 (`74b78ae`).
- Upstream changes since the last gin sync: OTel repo/usecase tracing restructure (2026-07-11),
  docker-compose `version:` removal, security updates, go1.26.5, rabbitmq/nginx bumps.

## Delta categorization

### Category A — Adopt from upstream (approved)

| Item | What | Commit |
|---|---|---|
| Persistent subpackage restructure | `repo/persistent/{task,translation,user}/` + per-repo `tracing.go` | C1 |
| amqp/nats file splits | `request/{task,translate,user}.go`, `response/{error,task}.go`, `v1/auth.go` | C2 |
| Missing test files | restapi `auth_test.go`, grpc `response_test.go`, usecase `task_test.go`/`user_test.go` | C3 |
| Entity split | `TranslationHistory` → `entity/translation.history.go` | C4 |
| Integration-test adoption | upstream's 25-test suite (user/task/translation × HTTP/gRPC/RMQ/NATS), gin-adapted | C5 |
| docker-compose modernization | drop `version:` spec | C6 |

### Category B — Keep gin-only (never port over)

- sqlc stack (`sqlc.yaml`, `queries/`, `sqlc/*.sql.go`), `pkg/cache` (otter), `pkg/validatorx`
- Middleware extras: `authz.go`, `body_size.go`, `cors.go`, `rate_limit.go`, `sleep.go`, `helper.go`
- `uow.go`, `entity/metadata.go`, dev tooling (`.air.toml`, `exec.sh`, `logs.sh`, `simplify.sh`, `.vscode/launch.json`)
- **Migrations** — incompatible with upstream: upstream uses `UUID`/`TIMESTAMP`, gin uses
  `VARCHAR`/`TIMESTAMPTZ` matching gin's string-ID entities and sqlc. Upstream's `NOT NULL`
  columns are already present in gin's migrations.

### Category C — Reconcile (keep gin's version)

- All framework diffs (fiber→gin): restapi controllers, middleware, `pkg/httpserver`, swagger
- Cosmetic gin changes: config struct names, jwt field naming, doc comments
- `entity/task.go`, `entity/user.go` are identical to upstream
- go.mod/go.sum: gin is ahead on most deps via its own dependabot; nothing to back-port
- Makefile: gin's version is ahead (`ENV_FILE` indirection, `db-up` target)

## Port design

### Pre-work

1. Commit current uncommitted state as baseline: `chore: lint fixes and regenerate artifacts`
   (contains the previous session's lint/format fixes and regenerated files)
2. `upstream` remote: already configured → `git@github.com:evrone/go-clean-template.git`

### Commit C1 — Persistent subpackage restructure

- Move `internal/repo/persistent/task_postgres.go` + `tracing_task.go` →
  `internal/repo/persistent/task/{task.go,tracing.go}`; same for translation and user repos
- Port upstream's per-repo `tracing.go` shape (from `533a0ab`), gin-adapted: sqlc queries stay
  in the repo implementations, tracing wraps the same `repo.TaskRepo`/`repo.TranslationRepo`/
  `repo.UserRepo` interfaces
- Rewire `internal/app/app.go` imports
- Regenerate mocks (`make mock`), verify build + tests + `make lint`

### Commit C2 — amqp/nats file splits

- Split inline request/response structs into `request/{task,translate,user}.go` and
  `response/{error,task}.go` for both amqp_rpc and nats_rpc v1 packages
- Port upstream's `v1/auth.go` handlers (from `2860e4c`), gin-adapted to gin's V1 struct and
  jwt Manager

### Commit C3 — Missing test files

- `internal/controller/restapi/middleware/auth_test.go`: port upstream test, gin-adapted
  (gin engine + `httptest` instead of fiber app)
- `internal/controller/grpc/v1/response/response_test.go`: import-path swap only (verified
  framework-free)
- `internal/usecase/task_test.go`, `internal/usecase/user_test.go`: port against gin's existing
  mockgen-generated mocks (same mock names and interface shapes — verified compatible)

### Commit C4 — Entity split

- Move `TranslationHistory` from `entity/translation.go` to `entity/translation.history.go`
  (upstream layout), update imports

### Commit C5 — Integration-test adoption

- Adopt upstream's full integration suite: helpers/task/translation/user split files, 25 test
  funcs covering register/login/profile, task CRUD + transitions, translation, across
  HTTP/gRPC/RMQ/NATS transports
- Gin-adapt: fiber→gin in HTTP tests; keep gin's existing translation coverage; remove the
  consolidated `integration_test.go` once content is carried over
- Largest diff volume; tests only run in docker-compose integration env (unchanged behavior)

### Commit C6 — docker-compose

- Drop `version: "3.9"` from `docker-compose.yml` and `docker-compose-integration-test.yml`

### Adaptation rules

- Every ported file: gin framework, `github.com/minhhoccode111/go-clean-template-gin` import
  paths, sqlc-consistent types
- Category B files are never touched; upstream code conflicting with them gets gin-adapted
- Category C files keep gin's version

### Verification

- After each commit: `go build ./...`, targeted `go test`, `make lint`
- Final gate: `make pre-commit` fully green (swag, protoc, mockgen, format, tests, lint)
- No DB/schema changes; migrations untouched

### Out of scope

- Migrations (schema-incompatible by design, see Category B)
- README/README_CN/README_RU sync
- Dockerfile/nginx/CI sync beyond docker-compose `version:` removal

## Recurring sync process (documented for the future)

1. `git fetch upstream`
2. `git diff --stat upstream/master HEAD` — identify new deltas
3. Classify into A (adopt) / B (keep) / C (reconcile) using this spec's categories
4. Port approved items as one commit per item, gin-adapted, `make pre-commit` gate
