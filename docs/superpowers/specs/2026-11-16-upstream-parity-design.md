# Upstream Parity — Phase 1 + Phase 2 Design

Date: 2026-11-16
Status: Approved

## Intent

Bring this Gin/sqlc fork closer to upstream `evrone/go-clean-template` without giving up the
fork's intentional differences: Gin instead of Fiber, sqlc instead of Squirrel, string IDs instead
of UUIDs, and the existing cache/validatorx/unit-of-work additions.

Upstream baseline: `d6d2f93` (`v1.19.3`).
Last synced upstream baseline: `376bcd4`.

## Scope

### Phase 1 — non-breaking parity

- Upgrade the project toolchain to Go 1.27 (`go.mod`, CI, integration Dockerfile).
- Replace the deprecated `gofumpt.extra-rules` linter setting with `extra.group-params`.
- Add PostgreSQL OpenTelemetry tracing through `github.com/exaring/otelpgx`.
- Add a Jaeger service to the local Docker stack and document tracing.
- Port upstream's agent documentation as `CLAUDE.md` with `AGENTS.md` symlink, adapted to this
  repository's Gin/sqlc/cache/unit-of-work/migration conventions.
- Add upstream's domain and observability documentation to `README.md`.

### Phase 2 — per-user translation history

- Replace the global translation history with per-user history.
- Use a destructive migration, as approved: drop and recreate `history`, matching upstream's
  approach but adapted to this repository's `VARCHAR` user IDs.
- Update sqlc queries and repository implementations to write and filter by `user_id`.
- Make the translation-history cache key per user.
- Add integration coverage proving one user cannot see another user's history.

## Non-goals

- No Fiber, Squirrel, prefork, or UUID migration.
- No README_CN/README_RU port.
- No JWT refactor.
- No unrelated dependency upgrades.

## Constraints

- Keep module path `github.com/minhhoccode111/go-clean-template-gin`.
- Keep existing `VARCHAR` IDs and `TIMESTAMPTZ` columns.
- Preserve Gin route surface and request/response shapes.
- Keep sqlc as the only SQL access path for repositories.
- Every phase must pass `make test`, `make linter-golangci`, and the Docker integration suite.
- Migration must apply cleanly on a fresh database and roll back cleanly.

## Success criteria

- `go.mod`, CI, and integration Dockerfile use Go 1.27.
- Jaeger receives traces from the app in the full Docker stack.
- `AGENTS.md` and `CLAUDE.md` describe the actual Gin/sqlc repository.
- Translation history rows carry `user_id`; SQL filters by it; cache keys include it.
- Two users translating independently never see each other's history.
- Full CI-equivalent checks pass.
