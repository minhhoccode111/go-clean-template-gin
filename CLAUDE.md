# go-clean-template-gin

Go clean-architecture reference service. Three domains (`user`, `task`, `translation`) exposed over
four transports (REST/Gin, gRPC, RabbitMQ RPC, NATS RPC) from one shared use-case layer.

Module path: `github.com/minhhoccode111/go-clean-template-gin`.
This is a Gin/sqlc fork of `evrone/go-clean-template`; keep Gin and sqlc conventions when porting
upstream changes.

## Commands — drive everything through the Makefile

The Makefile loads `.env` (falling back to `.env.example`) and exports it into every target. Use
the target, not the underlying tool.

| Task | Use this | Not this |
| --- | --- | --- |
| Install tool binaries | `make bin-deps` | `go install ...` |
| Unit tests | `make test` | `go test ./...` — the target adds `-race -covermode atomic` and scopes to `./internal/... ./pkg/...` |
| Integration tests | `make compose-up-integration-test` | `make integration-test` — host-side tests cannot resolve compose hostnames |
| Lint | `make linter-golangci` | bare `golangci-lint run` when you need the repo configuration |
| Format | `make format` | `gofmt` — the target runs `go fix`, `gofumpt`, and `gci` |
| Regenerate mocks | `make mock` | `mockgen ...` |
| Regenerate Swagger | `make swag-v1` | `swag init` |
| Regenerate protobuf | `make proto-v1` | `protoc ...` |
| Regenerate sqlc | `make sqlc` | hand-editing generated Go |
| Generate database schema | `make schema` | editing `docs/schema.sql` by hand |
| Generate TypeScript client | `make gen-ts` | hand-editing `api.ts` |
| Tidy / verify modules | `make deps` | `go mod tidy` alone |
| Vulnerability scan | `make deps-audit` | `govulncheck ./...` |
| Start dependencies | `make compose-up` | `docker compose up` |
| Start full stack | `make compose-up-all` | `docker compose up` |
| Tear down | `make compose-down` | `docker compose down` |
| Run locally | `make run` | `go run ./cmd/app` — the target regenerates code and builds with `-tags migrate` |
| Create a migration | `make migrate-create name=<name>` | `migrate create ...` |
| Apply migrations | `make migrate-up` | `migrate -path ... up` |
| Full check before pushing | `make pre-commit` | running steps by hand |

`make help` lists every target.

Traps:

- **`make integration-test` is not the Docker suite.** The real entry point is
  `make compose-up-integration-test`.
- **`make pre-commit` runs generators** (`swag`, protobuf, mocks, sqlc is separate, schema, TypeScript).
  Install tooling with `make bin-deps` first; `protoc` is not installed by that target.
- **`make schema` is best-effort.** It skips when the database is not reachable.

Before claiming a change is complete, run `make format`, `make linter-golangci`, and `make test`.
For transport or database changes, also run `make compose-up-integration-test`.

## Dependency rule

```
cmd/app → internal/app → internal/controller/*  ─┐
                       → internal/repo/*        ─┤→ internal/usecase (interfaces)
                                                 └→ internal/entity
```

- `internal/entity` — domain types and sentinel errors. Imports nothing from this module.
- `internal/usecase/contracts.go` — interfaces controllers call.
- `internal/usecase/<domain>/` — domain rules; tracing decorators live beside them.
- `internal/repo/contracts.go` — interfaces use cases call, including `UnitOfWork` and
  `TranslationCache`.
- `internal/repo/persistent/<domain>/` — Postgres implementations using sqlc.
- `internal/repo/persistent/sqlc/` — generated sqlc code, never hand-edited.
- `internal/repo/persistent/queries/` — hand-written SQL consumed by sqlc.
- `internal/repo/cache/` — Otter-backed implementations of repo cache contracts.
- `internal/repo/webapi/` — outbound HTTP APIs.
- `internal/controller/<transport>/v1/` — one package per transport with its own `request/` and
  `response/` DTOs. Controllers never import each other.
- `pkg/` — transport-agnostic infrastructure: servers, logger, JWT, Postgres, tracing, cache,
  validatorx. Must not import `internal/`.

Inner layers never import outer ones. Wiring happens once in `internal/app/app.go`
(`initUseCases`, `initServers`).

## Where code goes

| The code knows about… | It belongs in | Shape |
| --- | --- | --- |
| nothing but the domain | `internal/entity/<name>.go` | struct + methods |
| a rule spanning repositories | `internal/usecase/<domain>/<domain>.go` | method on `UseCase` |
| SQL, a table, a driver error | `internal/repo/persistent/<domain>/<domain>.go` | method on `Repo`, calls sqlc |
| a new SQL statement | `internal/repo/persistent/queries/<domain>.sql` | sqlc query, then `make sqlc` |
| an outbound API | `internal/repo/webapi/<name>.go` | webapi method |
| caching | `internal/repo/cache/<name>.go` | implements a repo cache contract |
| an HTTP/gRPC/RPC status or envelope | `internal/controller/<transport>/v1/<domain>.go` | handler |
| a server, pool, client, or middleware with no domain knowledge | `pkg/<name>/` | reusable package |

## Handler contract

- Gin handlers are methods on `V1`, receive `*gin.Context`, and write with `c.JSON` or
  `errorResponse`.
- Validate request bodies with the injected `validatorx`/`validator.Validate`.
- Extract the authenticated user from middleware context; never trust a user ID from the body.
- Return transport DTOs, not entities, when the transport shape differs.
- Keep route registration in the transport's `router.go`; keep handlers in per-domain files.
- gRPC, AMQP, and NATS handlers call the same use-case interfaces as REST.

## Conventions that will bite you

- IDs are `VARCHAR(255)` strings in entities, migrations, and sqlc. Do not introduce UUID types
  without changing the schema and all transports.
- Timestamps are `TIMESTAMPTZ`; use `persistent.Timestamptz` when passing `time.Time` to sqlc.
- Translation history is per user. SQL filters by `user_id`; cache keys include the user ID.
- Use the unit of work for multi-repository transactions. Do not open ad-hoc transactions in
  use cases.
- Keep tracing decorators in place when refactoring repositories or use cases.
- Use sentinel errors from `internal/entity`; controllers map them to transport-specific statuses.
- Never hand-edit generated files.

## Generated files — do not hand-edit

- `internal/repo/persistent/sqlc/*.go` — `make sqlc`
- `internal/usecase/mocks_*_test.go` — `make mock`
- `docs/proto/**/*.pb.go` — `make proto-v1`
- `docs/swagger.*` — `make swag-v1`
- `api.ts` — `make gen-ts`
- `docs/schema.sql` — `make schema` when the database is reachable

## Database and migrations

- Migrations live in `migrations/` and are applied automatically at app startup with
  `-tags migrate`, or manually with `make migrate-up`.
- Create with `make migrate-create name=<name>`.
- Add or change SQL in `internal/repo/persistent/queries/*.sql`, then run `make sqlc`.
- Update `docs/schema.sql` through `make schema`; do not patch it by hand when the database is
  available.
- Integration tests run against a fresh database in `docker-compose-integration-test.yml`.

## Testing

- Unit tests: `make test` (`internal/...` and `pkg/...`).
- Integration tests: `make compose-up-integration-test`. They cover HTTP, gRPC, AMQP RPC, and
  NATS RPC.
- Add integration coverage for cross-user isolation, auth failures, and transport-specific
  behavior.
- Test names are used to build unique user emails; keep tests independent.

## Adding a use case to an existing domain

1. Add the method to the use-case contract in `internal/usecase/contracts.go`.
2. Implement it in `internal/usecase/<domain>/`.
3. Add or extend the tracing decorator.
4. Add repository contracts and implementations if new persistence is needed.
5. Wire it in `internal/app/app.go`.
6. Add unit tests and, for transport behavior, integration tests.
7. Add handlers in every transport that should expose it.
8. Regenerate Swagger/protobuf/TypeScript when the API changes.

## Adding a whole domain

Create `internal/entity/<domain>.go`, `internal/usecase/<domain>/`, repository contracts and
implementations, transport DTO/handler packages, migrations, sqlc queries, and integration tests.
Follow the existing user/task/translation layout rather than inventing a new one.

## API versioning

- REST and gRPC live under `v1` packages.
- AMQP and NATS use `v1.<domain>.<method>` subjects/routing keys.
- Add `v2` packages for breaking changes; keep `v1` wired until clients migrate.
- Regenerate `docs/proto`, Swagger, and `api.ts` after API changes.

## CI

CI runs on pull requests: golangci-lint, yamllint, hadolint, dotenv-linter, dependency audit, unit
tests, and Docker integration tests. There is no push-to-main workflow.
