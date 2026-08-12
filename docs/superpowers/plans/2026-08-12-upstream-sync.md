# Upstream Sync Port Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Port approved Category A deltas from `upstream/master` (`376bcd4`, evrone/go-clean-template) into go-clean-template-gin — persistent repo subpackages, amqp/nats file splits, missing test files, entity split, integration-test adoption, docker-compose modernization — without losing gin-only changes.

**Architecture:** Pure layout/structural sync. Six commits, one per adoption item, gin-adapted (fiber→gin, sqlc preserved, gin import paths). Gin-only code (sqlc stack, cache, validatorx, middleware extras, uow, migrations, gin controllers) is never touched. Upstream files are copied via `git show upstream/master:<path>` then transformed.

**Tech Stack:** Go 1.26, gin, sqlc, golangci-lint, gofumpt, mockgen, testify, gomock, otel.

**Upstream ref:** `upstream/master` (already fetched; remote `upstream` configured → `git@github.com:evrone/go-clean-template.git`).

## Global Constraints

- Import path prefix: `github.com/minhhoccode111/go-clean-template-gin` (never `github.com/evrone/...`)
- Framework: gin — never port fiber code; HTTP handler tests use `gin.New()` + `httptest`
- sqlc stack (`sqlc.yaml`, `queries/`, `sqlc/*.sql.go`) is gin-only: never modify
- Migrations are gin-only (VARCHAR/TIMESTAMPTZ): never touch
- Gin-only packages never touched: `pkg/cache`, `pkg/validatorx`, `internal/repo/persistent/uow.go`, restapi middleware extras
- `make pre-commit` must be green at end (swag, protoc, mockgen, gofumpt, gci, tests, golangci-lint)
- Work on branch `develop`; one commit per task (74b78ae-style port commits)

---

### Task 0: Baseline commit of current uncommitted state

**Files:**
- All currently modified tracked files (lint fixes from previous session + regenerated artifacts)

**Interfaces:**
- Produces: clean working tree; every later commit contains only its own port changes

- [ ] **Step 1: Inspect current state**

Run: `git status --short`
Expected: ~26 modified files (lint/format fixes, regenerated docs/proto/mocks/api.ts) — the residue of the previous session's `make pre-commit` fixes.

- [ ] **Step 2: Commit baseline**

```bash
git add -A
git commit -m "chore: lint fixes and regenerate artifacts"
```

- [ ] **Step 3: Verify**

Run: `git status --short`
Expected: empty (clean tree)

---

### Task 1 (C1): Move persistent repos into task/translation/user subpackages

Adopts upstream's layout (`internal/repo/persistent/{task,translation,user}/` with `Repo` struct + `New` constructor + `tracedRepo` tracing wrapper). Gin's sqlc implementations move verbatim, renamed to upstream conventions.

**Files:**
- Create: `internal/repo/persistent/task/task.go` (from `task_postgres.go`)
- Create: `internal/repo/persistent/task/tracing.go` (from `tracing_task.go`)
- Create: `internal/repo/persistent/translation/translation.go` (from `translation_postgres.go`)
- Create: `internal/repo/persistent/translation/tracing.go` (from `tracing_translation.go`)
- Create: `internal/repo/persistent/user/user.go` (from `user_postgres.go`)
- Create: `internal/repo/persistent/user/tracing.go` (from `tracing_user.go`)
- Modify: `internal/repo/persistent/helper.go` (export `Timestamptz`)
- Modify: `internal/app/app.go` (repo wiring imports)
- Delete: `internal/repo/persistent/task_postgres.go`, `tracing_task.go`, `translation_postgres.go`, `tracing_translation.go`, `user_postgres.go`, `tracing_user.go`

**Interfaces:**
- Consumes: existing `repo.TaskRepo`/`repo.TranslationRepo`/`repo.UserRepo` interfaces (unchanged), `sqlc` package (unchanged)
- Produces: `task.New(pg *postgres.Postgres) repo.TaskRepo`, `translation.New(pg *postgres.Postgres) repo.TranslationRepo`, `user.New(pg *postgres.Postgres) repo.UserRepo`; `persistent.Timestamptz(t time.Time) pgtype.Timestamptz`; app.go calls `persistTaskRepo.New(pg)` etc.

- [ ] **Step 1: Export the timestamp helper in the root package**

Modify `internal/repo/persistent/helper.go`:

```go
package persistent

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func Timestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}
```

- [ ] **Step 2: Move the task repo into its subpackage**

```bash
mkdir -p internal/repo/persistent/task internal/repo/persistent/translation internal/repo/persistent/user
git mv internal/repo/persistent/task_postgres.go internal/repo/persistent/task/task.go
git mv internal/repo/persistent/tracing_task.go internal/repo/persistent/task/tracing.go
```

- [ ] **Step 3: Adapt task.go to subpackage (package, rename, helper call)**

In `internal/repo/persistent/task/task.go`:
- `package persistent` → `package task`
- Add import: `"github.com/minhhoccode111/go-clean-template-gin/internal/repo/persistent"` (needed for `persistent.Timestamptz`)
- `pgTimestamptz(` → `persistent.Timestamptz(` (3 occurrences: lines with `CreatedAt:`, `UpdatedAt:` in `Store` and `Update`)
- `// TaskRepo implements repo.TaskRepo using sqlc.` → `// Repo implements repo.TaskRepo using sqlc.`
- `type TaskRepo struct {` → `type Repo struct {`
- `// NewTaskRepo returns a Task repository instrumented with OpenTelemetry tracing spans.` → `// New returns a Task repository instrumented with OpenTelemetry tracing spans.`
- `func NewTaskRepo(pg *postgres.Postgres) repo.TaskRepo {` → `func New(pg *postgres.Postgres) repo.TaskRepo {`
- `return newTracedTask(&TaskRepo{` → `return newTraced(&Repo{`
- `func (r *TaskRepo)` → `func (r *Repo)` (5 receivers: Store, GetByID, List, Update, Delete)

- [ ] **Step 4: Adapt task/tracing.go**

In `internal/repo/persistent/task/tracing.go`:
- `package persistent` → `package task`
- `const _tracerNameTask = "github.com/minhhoccode111/go-clean-template-gin/internal/repo/persistent/task"` → `const _tracerName = "github.com/minhhoccode111/go-clean-template-gin/internal/repo/persistent/task"`
- `type tracedTaskRepo struct` → `type tracedRepo struct`
- `func newTracedTask(next repo.TaskRepo) repo.TaskRepo` → `func newTraced(next repo.TaskRepo) repo.TaskRepo`
- `return &tracedTaskRepo{next: next}` → `return &tracedRepo{next: next}`
- `func (r *tracedTaskRepo)` → `func (r *tracedRepo)` (6 receivers)

- [ ] **Step 5: Move translation repo into its subpackage**

```bash
git mv internal/repo/persistent/translation_postgres.go internal/repo/persistent/translation/translation.go
git mv internal/repo/persistent/tracing_translation.go internal/repo/persistent/translation/tracing.go
```

In `internal/repo/persistent/translation/translation.go`:
- `package persistent` → `package translation`
- `// TranslationRepo implements repo.TranslationRepo using sqlc.` → `// Repo implements repo.TranslationRepo using sqlc.`
- `type TranslationRepo struct` → `type Repo struct`
- `// NewTranslationRepo returns a Translation repository instrumented with OpenTelemetry tracing spans.` → `// New returns a Translation repository instrumented with OpenTelemetry tracing spans.`
- `func NewTranslationRepo(pg *postgres.Postgres) repo.TranslationRepo` → `func New(pg *postgres.Postgres) repo.TranslationRepo`
- `return newTracedTranslation(&TranslationRepo{` → `return newTraced(&Repo{`
- `func (r *TranslationRepo)` → `func (r *Repo)` (all receivers in file)

In `internal/repo/persistent/translation/tracing.go`:
- `package persistent` → `package translation`
- `_tracerNameTranslation` → `_tracerName` (declaration and 3 uses)
- `tracedTranslationRepo` → `tracedRepo` (declaration + constructor + receiver + `&tracedRepo{`)

- [ ] **Step 6: Move user repo into its subpackage**

```bash
git mv internal/repo/persistent/user_postgres.go internal/repo/persistent/user/user.go
git mv internal/repo/persistent/tracing_user.go internal/repo/persistent/user/tracing.go
```

In `internal/repo/persistent/user/user.go`:
- `package persistent` → `package user`
- Add import: `"github.com/minhhoccode111/go-clean-template-gin/internal/repo/persistent"`
- `pgTimestamptz(` → `persistent.Timestamptz(` (2 occurrences in `Store`)
- `// UserRepo implements repo.UserRepo using sqlc.` → `// Repo implements repo.UserRepo using sqlc.`
- `type UserRepo struct` → `type Repo struct`
- `// NewUserRepo returns a User repository instrumented with OpenTelemetry tracing spans.` → `// New returns a User repository instrumented with OpenTelemetry tracing spans.`
- `func NewUserRepo(pg *postgres.Postgres) repo.UserRepo` → `func New(pg *postgres.Postgres) repo.UserRepo`
- `return newTracedUser(&UserRepo{` → `return newTraced(&Repo{`
- `func (r *UserRepo)` → `func (r *Repo)` (all receivers)

In `internal/repo/persistent/user/tracing.go`:
- `package persistent` → `package user`
- `_tracerNameUser` → `_tracerName` (declaration + 2 uses)
- `tracedUserRepo` → `tracedRepo` (declaration + constructor + receiver + `&tracedRepo{`)

- [ ] **Step 7: Rewire app.go repo wiring**

In `internal/app/app.go`, replace the imports block for repos and the `initUseCases` repo construction:

```go
import (
	...
	persistTaskRepo "github.com/minhhoccode111/go-clean-template-gin/internal/repo/persistent/task"
	persistTranslationRepo "github.com/minhhoccode111/go-clean-template-gin/internal/repo/persistent/translation"
	persistUserRepo "github.com/minhhoccode111/go-clean-template-gin/internal/repo/persistent/user"
	"github.com/minhhoccode111/go-clean-template-gin/internal/repo/persistent/sqlc"
	...
)
```

Replace in `initUseCases`:

```go
	translationRepo := persistTranslationRepo.New(pg)
	taskRepo := persistTaskRepo.New(pg)
	userRepo := persistUserRepo.New(pg)
```

Remove now-unused imports from app.go (e.g. `persistent` if only used for repo constructors — keep it if `persistent.NewUnitOfWork` still needs it).

- [ ] **Step 8: Build + format + lint**

```bash
gofumpt -l -w internal/repo/persistent internal/app/app.go
gci write internal/repo/persistent internal/app/app.go --skip-generated -s standard -s default
go build ./...
go test ./internal/... 2>&1 | grep -E "FAIL|ok" | tail -8
golangci-lint run 2>&1 | tail -3
```

Expected: build OK, no FAIL, `0 issues.`

- [ ] **Step 9: Regenerate mocks and verify**

Run: `make mock`
Expected: `mocks_repo_test.go`/`mocks_usecase_test.go` unchanged (contracts untouched); confirm with `git status --short | grep mocks` → no output.

- [ ] **Step 10: Commit**

```bash
git add -A
git commit -m "feat(repo): move persistent repos into task/translation/user subpackages"
```

---

### Task 2 (C2): amqp/nats request/response file splits + auth helper

Adopts upstream's file layout: per-domain `request/*.go` + `response/*.go` files and `v1/auth.go` with `extractUserID` (payload envelope stays `{"token": ..., "data": ...}` — wire-compatible). `user.go` in both transports already uses `request`/`response` packages (from the earlier auth pull) — only `task.go` and `translation.go` need rewriting, plus the missing files.

**Files:**
- Modify: `internal/controller/amqp_rpc/v1/request/auth.go` (add `AuthenticatedRequest`)
- Create: `internal/controller/amqp_rpc/v1/request/task.go`
- Create: `internal/controller/amqp_rpc/v1/request/translate.go`
- Create: `internal/controller/amqp_rpc/v1/response/error.go`
- Create: `internal/controller/amqp_rpc/v1/response/task.go`
- Create: `internal/controller/amqp_rpc/v1/auth.go`
- Modify: `internal/controller/amqp_rpc/v1/task.go`, `internal/controller/amqp_rpc/v1/translation.go`
- Mirror all of the above in `internal/controller/nats_rpc/v1/...`

**Interfaces:**
- Consumes: existing `V1` struct (`r.j *jwt.Manager`, `r.tk`, `r.t`, `r.u`, `r.v`, `r.l`), `server.CallHandler` from `pkg/rabbitmq/rmq_rpc/server` and `pkg/nats/nats_rpc/server`
- Produces: `extractUserID(d *amqp.Delivery, jwtManager *jwt.Manager) (string, json.RawMessage, error)`; `request.AuthenticatedRequest{Token string; Data json.RawMessage}`; typed structs `request.CreateTask`, `request.GetTask`, `request.ListTasks`, `request.UpdateTask`, `request.TransitionTask`, `request.DeleteTask`, `request.Translate`; `response.TaskList`, `response.DeleteStatus`, `response.Error`

- [ ] **Step 1: Add AuthenticatedRequest to amqp request/auth.go**

Append to `internal/controller/amqp_rpc/v1/request/auth.go`:

```go
package request

import "github.com/goccy/go-json"

// AuthenticatedRequest is the envelope for all authenticated RPC calls.
type AuthenticatedRequest struct {
	Token string          `json:"token" validate:"required"`
	Data  json.RawMessage `json:"data"`
}

// Register -.
type Register struct {
	Username string `json:"username" validate:"required,min=3,max=255"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// Login -.
type Login struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
```

- [ ] **Step 2: Create amqp request files (copy from upstream, swap imports)**

```bash
git show upstream/master:internal/controller/amqp_rpc/v1/request/task.go > internal/controller/amqp_rpc/v1/request/task.go
git show upstream/master:internal/controller/amqp_rpc/v1/request/translate.go > internal/controller/amqp_rpc/v1/request/translate.go
git show upstream/master:internal/controller/amqp_rpc/v1/response/error.go > internal/controller/amqp_rpc/v1/response/error.go
git show upstream/master:internal/controller/amqp_rpc/v1/response/task.go > internal/controller/amqp_rpc/v1/response/task.go
```

These files are import-free except `response/task.go` (imports `github.com/evrone/go-clean-template/internal/entity`). Fix it:

```bash
sed -i 's|github.com/evrone/go-clean-template|github.com/minhhoccode111/go-clean-template-gin|' internal/controller/amqp_rpc/v1/response/task.go
```

- [ ] **Step 3: Create amqp v1/auth.go**

```bash
git show upstream/master:internal/controller/amqp_rpc/v1/auth.go > internal/controller/amqp_rpc/v1/auth.go
sed -i 's|github.com/evrone/go-clean-template|github.com/minhhoccode111/go-clean-template-gin|' internal/controller/amqp_rpc/v1/auth.go
```

Verify content matches the upstream helper (it must reference `request.AuthenticatedRequest` and `jwtManager.ParseToken`).

- [ ] **Step 4: Rewrite amqp task.go to upstream handler pattern**

Replace the inline data structs and per-handler token parsing in `internal/controller/amqp_rpc/v1/task.go`:

1. Delete the 6 inline data structs (`createTaskData`, `getTaskData`, `listTasksData`, `updateTaskData`, `transitionTaskData`, `deleteTaskData`).
2. Update imports: remove `"github.com/goccy/go-json"`, remove `"github.com/minhhoccode111/go-clean-template-gin/internal/entity"` (if unused after rewrite — `listTasks` uses `entity.TaskStatus`; keep it), add `"github.com/minhhoccode111/go-clean-template-gin/internal/controller/amqp_rpc/v1/request"` and `"github.com/minhhoccode111/go-clean-template-gin/internal/controller/amqp_rpc/v1/response"`.
3. Each handler becomes (pattern for `createTask`; repeat for all 6):

```go
func (r *V1) createTask() server.CallHandler {
	return func(d *amqp.Delivery) (any, error) {
		userID, rawData, err := extractUserID(d, r.j)
		if err != nil {
			r.l.Error(err, "amqp_rpc - V1 - createTask")

			return nil, fmt.Errorf("amqp_rpc - V1 - createTask - extractUserID: %w", err)
		}

		var reqData request.CreateTask

		if err := json.Unmarshal(rawData, &reqData); err != nil {
			r.l.Error(err, "amqp_rpc - V1 - createTask")

			return nil, fmt.Errorf("amqp_rpc - V1 - createTask - json.Unmarshal: %w", err)
		}

		task, err := r.tk.Create(context.Background(), userID, reqData.Title, reqData.Description)
		if err != nil {
			r.l.Error(err, "amqp_rpc - V1 - createTask")

			return nil, fmt.Errorf("amqp_rpc - V1 - createTask: %w", err)
		}

		return task, nil
	}
}
```

Handler → request struct + usecase call mapping (identical for nats):

| Handler | request struct | Usecase call |
|---|---|---|
| `createTask` | `request.CreateTask` | `r.tk.Create(ctx, userID, reqData.Title, reqData.Description)` |
| `getTask` | `request.GetTask` | `r.tk.Get(ctx, userID, reqData.ID)` |
| `listTasks` | `request.ListTasks` | keep the `entity.TaskStatus`/`status` pointer block; `r.tk.List(ctx, userID, status, reqData.Limit, reqData.Offset)`; return `response.TaskList{Tasks: tasks, Total: total}` |
| `updateTask` | `request.UpdateTask` | `r.tk.Update(ctx, userID, reqData.ID, reqData.Title, reqData.Description)` |
| `transitionTask` | `request.TransitionTask` | `r.tk.Transition(ctx, userID, reqData.ID, entity.TaskStatus(reqData.Status))` |
| `deleteTask` | `request.DeleteTask` | `r.tk.Delete(ctx, userID, reqData.ID)`; return `response.DeleteStatus{Status: "deleted"}, nil` |

Note: `listTasks` previously returned `map[string]any{"tasks": tasks, "total": total}` — the `response.TaskList` JSON is identical (`{"tasks": [...], "total": n}`).

- [ ] **Step 5: Rewrite amqp translation.go**

Same pattern in `internal/controller/amqp_rpc/v1/translation.go`:

- `getHistory`: no data payload — use `userID, _, err := extractUserID(d, r.j)` (blank the raw data, otherwise the `rawData` variable is unused and the build fails), then call `r.t.History(context.Background(), userID)`.
- `translate`: delete `translateReqData`; `var reqData request.Translate`; call `r.t.Translate(context.Background(), userID, entity.Translation{Source: reqData.Source, Destination: reqData.Destination, Original: reqData.Original})`.

- [ ] **Step 6: Mirror all amqp changes into nats_rpc**

Repeat Steps 1–5 with paths `internal/controller/nats_rpc/v1/...`, and:
- `auth.go`: replace `*amqp.Delivery` with `*nats.Msg` in `extractUserID`
- `server.CallHandler` import: `"github.com/minhhoccode111/go-clean-template-gin/pkg/nats/nats_rpc/server"`
- nats handler signatures: `func(msg *nats.Msg) (any, error)`

- [ ] **Step 7: Build + lint**

```bash
gofumpt -l -w internal/controller/amqp_rpc internal/controller/nats_rpc
gci write internal/controller/amqp_rpc internal/controller/nats_rpc --skip-generated -s standard -s default
go build ./...
golangci-lint run 2>&1 | tail -3
```

Expected: build OK, `0 issues.`

- [ ] **Step 8: Verify payload compatibility (unit-level sanity)**

Run: `go vet ./internal/controller/amqp_rpc/... ./internal/controller/nats_rpc/...`
Expected: no output (struct tags/envelope unchanged).

- [ ] **Step 9: Commit**

```bash
git add -A
git commit -m "refactor(amqp,nats): split request/response structs into per-domain files, add auth helper"
```

---

### Task 3 (C3): Port upstream test files

**Files:**
- Create: `internal/controller/restapi/middleware/auth_test.go` (gin-adapted rewrite)
- Create: `internal/controller/grpc/v1/response/response_test.go` (path swap)
- Create: `internal/usecase/task_test.go` (path swap)
- Create: `internal/usecase/user_test.go` (path swap)

**Interfaces:**
- Consumes: `middleware.Auth(jwtManager *jwt.Manager) gin.HandlerFunc`, `jwt.Manager.GenerateToken(userID string) (string, error)`, gin's generated mocks `NewMockTaskRepo`/`NewMockUserRepo` (exist in `internal/usecase/mocks_repo_test.go`), `task.New(r repo.TaskRepo) usecase.Task`, `user.New(r repo.UserRepo, j *jwt.Manager) usecase.User`

- [ ] **Step 1: Write the gin-adapted restapi auth middleware test**

Create `internal/controller/restapi/middleware/auth_test.go`:

```go
package middleware_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minhhoccode111/go-clean-template-gin/internal/controller/restapi/middleware"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestApp(t *testing.T) (*gin.Engine, *jwt.Manager) {
	t.Helper()

	gin.SetMode(gin.TestMode)

	jwtManager := jwt.New("test-secret", time.Hour)

	app := gin.New()
	app.Use(middleware.Auth(jwtManager))
	app.GET("/test", func(c *gin.Context) {
		v, ok := c.Get("userID")
		if !ok {
			c.Status(http.StatusUnauthorized)

			return
		}

		userID, ok := v.(string)
		if !ok {
			c.Status(http.StatusUnauthorized)

			return
		}

		c.String(http.StatusOK, userID)
	})

	return app, jwtManager
}

func TestAuthMiddleware(t *testing.T) {
	t.Parallel()

	app, jwtManager := newTestApp(t)

	validToken, err := jwtManager.GenerateToken("user-id-123")
	require.NoError(t, err)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "missing header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid format",
			authHeader:     "Basic xxx",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer invalid",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "valid token",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
			expectedBody:   "user-id-123",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/test", http.NoBody)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}

			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)

			if tc.expectedBody != "" {
				body, readErr := io.ReadAll(w.Result().Body)
				require.NoError(t, readErr)
				assert.Equal(t, tc.expectedBody, string(body))
			}
		})
	}
}
```

- [ ] **Step 2: Run the new test**

Run: `go test ./internal/controller/restapi/middleware/... -run TestAuthMiddleware -v`
Expected: PASS (4 subtests)

- [ ] **Step 3: Port grpc response_test.go (path swap only)**

```bash
git show upstream/master:internal/controller/grpc/v1/response/response_test.go > internal/controller/grpc/v1/response/response_test.go
sed -i 's|github.com/evrone/go-clean-template|github.com/minhhoccode111/go-clean-template-gin|' internal/controller/grpc/v1/response/response_test.go
```

- [ ] **Step 4: Run it**

Run: `go test ./internal/controller/grpc/v1/response/... -v`
Expected: PASS

- [ ] **Step 5: Port usecase task_test.go + user_test.go**

```bash
git show upstream/master:internal/usecase/task_test.go > internal/usecase/task_test.go
git show upstream/master:internal/usecase/user_test.go > internal/usecase/user_test.go
sed -i 's|github.com/evrone/go-clean-template|github.com/minhhoccode111/go-clean-template-gin|' internal/usecase/task_test.go internal/usecase/user_test.go
```

- [ ] **Step 6: Run them — fix compile mismatches if any**

Run: `go test ./internal/usecase/... -v 2>&1 | grep -E "FAIL|PASS: Test(Task|User)|ok|error" | head -20`
Expected: all Task/User subtests PASS.

If a compile error appears (e.g. a usecase signature differing from upstream — `Create`, `Register`, `GetUser`, `List`), the fix is to align the test call to gin's actual signature in `internal/usecase/contracts.go` (e.g. `Create(ctx, userID, title, description)`), NOT to change the usecase. Document any deviation in the commit message.

- [ ] **Step 7: Lint**

Run: `golangci-lint run 2>&1 | tail -3`
Expected: `0 issues.` (if a ported test trips paralleltest/tparallel, add `t.Parallel()` per the linter's complaint, matching the existing test style in `internal/usecase/translation_test.go`)

- [ ] **Step 8: Commit**

```bash
git add -A
git commit -m "test: port upstream middleware, grpc response, and usecase tests"
```

---

### Task 4 (C4): Split TranslationHistory entity

**Files:**
- Create: `internal/entity/translation.history.go`
- Modify: `internal/entity/translation.go` (remove `TranslationHistory`)

**Interfaces:**
- Produces: `entity.TranslationHistory` (same shape: `History []Translation`) — package-internal move, no import changes anywhere

- [ ] **Step 1: Create the split file**

Create `internal/entity/translation.history.go`:

```go
package entity

// TranslationHistory -.
type TranslationHistory struct {
	History []Translation `json:"history"`
}
```

- [ ] **Step 2: Remove the struct from translation.go**

In `internal/entity/translation.go`, delete:

```go
// TranslationHistory -.
type TranslationHistory struct {
	History []Translation `json:"history"`
}
```

- [ ] **Step 3: Verify build + tests**

```bash
go build ./...
go test ./internal/entity/... ./internal/usecase/... ./internal/controller/restapi/... 2>&1 | grep -E "FAIL|ok" | tail -6
```

Expected: build OK, no FAIL (same package — all references still resolve).

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "refactor(entity): split TranslationHistory into its own file"
```

---

### Task 5 (C5): Adopt upstream integration-test suite

Upstream's integration tests are black-box HTTP/gRPC/RMQ/NATS client tests (no fiber imports — verified). Routes/proto methods are identical between the repos (verified). The upstream suite is a superset of gin's: register/login/profile, task CRUD + transitions, translation — across all 4 transports.

**Files:**
- Create: `integration-test/helpers_test.go`, `integration-test/task_test.go`, `integration-test/translation_test.go`, `integration-test/user_test.go`
- Delete: `integration-test/integration_test.go`

**Interfaces:**
- Consumes: existing `Makefile` `test` target (runs `go test ./integration-test/...` only under docker-compose integration env — compile-check locally with `go build`/`go vet`)
- Produces: upstream-equivalent integration coverage, gin-imported

- [ ] **Step 1: Copy + transform the four files**

```bash
git show upstream/master:integration-test/helpers_test.go > integration-test/helpers_test.go
git show upstream/master:integration-test/task_test.go > integration-test/task_test.go
git show upstream/master:integration-test/translation_test.go > integration-test/translation_test.go
git show upstream/master:integration-test/user_test.go > integration-test/user_test.go
sed -i 's|github.com/evrone/go-clean-template|github.com/minhhoccode111/go-clean-template-gin|' integration-test/helpers_test.go integration-test/task_test.go integration-test/translation_test.go integration-test/user_test.go
```

- [ ] **Step 2: Delete gin's consolidated test file**

```bash
git rm integration-test/integration_test.go
```

- [ ] **Step 3: Compile-check**

```bash
go vet ./integration-test/...
```

Expected: no errors. If a vet error appears, it will be a signature/type mismatch (e.g. a gRPC client method name) — fix by aligning to gin's generated protos in `docs/proto/v1/` (they match upstream's; the earlier grpc client check confirmed the same method set).

- [ ] **Step 4: Cross-check gin-specific expectations**

Grep the ported tests for anything gin-specific that differs from upstream behavior (auth header handling, response envelope):
- `loginUser`/`registerUser` helpers use `/auth/register` + `/auth/login` — gin's router has these exact routes (verified in `internal/controller/restapi/v1/router.go`)
- Response bodies: gin returns `response.Token{Token: ...}` — check `user_test.go`'s token extraction matches (upstream JSON shape `{"token": "..."}` — same tags)

If any assertion shape differs, fix the test, not the app.

- [ ] **Step 5: Format + lint**

```bash
gofumpt -l -w integration-test
golangci-lint run 2>&1 | tail -3
```

Note: `.golangci.yml` already excludes `integration-test` from `godot` and `paralleltest`. Expected: `0 issues.`

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "test(integration): adopt upstream integration suite, gin-adapted"
```

---

### Task 6 (C6): Remove docker-compose version spec

**Files:**
- Modify: `docker-compose.yml`, `docker-compose-integration-test.yml`

- [ ] **Step 1: Remove the version line**

```bash
sed -i '/^version: "3.9"$/d' docker-compose.yml docker-compose-integration-test.yml
```

- [ ] **Step 2: Verify**

Run: `head -3 docker-compose.yml docker-compose-integration-test.yml`
Expected: no `version:` line; `services:` at top (matching upstream).

- [ ] **Step 3: Commit**

```bash
git add -A
git commit -m "chore(docker): remove version specification from docker-compose files"
```

---

### Task 7: Final gate — full pre-commit

- [ ] **Step 1: Run the full pipeline**

Run: `make pre-commit`
Expected: exit 0 — swag docs regenerate, protoc compiles, mocks regenerate, gofumpt/gci pass, all `go test -race` suites PASS, `golangci-lint run` reports `0 issues.`

- [ ] **Step 2: Handle regeneration noise**

After pre-commit, regenerated files (docs/, api.ts, *.pb.go, mocks) may show as modified. Verify they are content-identical to what pre-commit produced earlier:

```bash
git status --short | head -20
```

If only regenerated artifacts changed, commit them:

```bash
git add -A
git commit -m "chore: regenerate artifacts"
```

- [ ] **Step 3: Verify upstream-sync completeness**

```bash
git log --oneline develop --not upstream/master | head -12
```

Expected: the 6 port commits + baseline visible. Then confirm the adopted files match upstream:

```bash
git diff upstream/master HEAD --stat | tail -3
```

Expected: no lines for `integration-test/helpers_test.go` etc. (adopted files now identical except import paths — diff should show nothing for files that were pure copies; restructured files will still differ by design).

- [ ] **Step 4: Update the spec's status**

In `docs/superpowers/specs/2026-08-12-upstream-sync-design.md`, change `Status: Approved (design), pending implementation plan` → `Status: Implemented (2026-08-12)`. Commit:

```bash
git add docs/superpowers/specs/2026-08-12-upstream-sync-design.md
git commit -m "docs: mark upstream sync spec implemented"
```

---

## Out of Scope (from spec, do not implement)

- Migrations (UUID vs VARCHAR schema — incompatible by design)
- README/README_CN/README_RU sync
- Dockerfile/nginx/CI sync beyond `version:` removal
- Category B files (sqlc stack, cache, validatorx, middleware extras, uow, tooling)
- Category C files (gin's versions of controllers/config/jwt/Makefile/go.mod)
