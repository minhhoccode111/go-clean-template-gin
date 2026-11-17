# Observability Follow-ups — Jaeger v2, Transport Spans, Unit-of-Work Fix

Date: 2026-11-17
Status: Approved (user said "build it like the upstream repo")

## Intent

Close the three follow-ups left after the upstream parity merge:

1. Replace the end-of-life Jaeger v1 image with Jaeger v2.
2. Add transport-level tracing to REST, AMQP RPC, and NATS RPC, mirroring upstream.
3. Stop `translation.Repo.Store` from opening its own transaction, so a unit-of-work-bound
   repository writes inside the unit-of-work transaction.

## Decisions

- Mirror upstream `d6d2f93` behavior exactly for span names, span kinds, tracer names, and
  attributes. Upstream quirk retained: RPC client spans root at the client's lifecycle context
  (AMQP) or `context.Background()` (NATS), not at the caller's request context.
- Apply Gin tracing only when `TRACING_ENABLED` is true, scoped to the `/v1` group, matching
  upstream's scoping away from `/healthz`, `/metrics`, and `/swagger`.
- Jaeger image `jaegertracing/jaeger:2.23.0`; drop the legacy `6831/udp` agent port.
- `CallHandler` gains `context.Context` in both RPC servers.
- `Store` performs a single statement without an inner transaction; the statement is atomic, and
  a transaction-bound repository then joins the unit-of-work transaction.

## Non-goals

- No `ctx` parameter on `RemoteCall`.
- No retention, ordering, or pagination change for translation history.
- No change to gRPC tracing.

## Success criteria

- Jaeger v2 receives traces and exposes the UI on `16686`.
- One HTTP translate request produces a trace containing the Gin server span and the use-case span.
- One NATS or AMQP request produces one trace containing both the publish and process spans.
- A `Store` call inside a failing unit of work leaves no row behind.
- `make test`, `make linter-golangci`, and the Docker integration suite pass.

## Tasks

### Task 1: Jaeger v2

- `docker-compose.yml`: image `jaegertracing/jaeger:2.23.0`, remove `6831/udp`.
- Verify the container starts and listens on `4317` and `16686`.
- README: no stale v1 references.

### Task 2: REST tracing

- Add `go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin`.
- Test first: minimal Gin engine with `otelgin.Middleware`, in-memory span recorder, assert a
  server span is produced for a request.
- `internal/controller/restapi/router.go`: under `cfg.Tracing.Enabled`, add
  `apiV1Group.Use(otelgin.Middleware(cfg.App.Name))`.

### Task 3: AMQP and NATS transport spans

- AMQP client: start `rmq_rpc.publish <handler>` client span in `RemoteCall`; thread `ctx` into
  `publish`; inject propagation into a `headers` table.
- AMQP server: `CallHandler func(context.Context, *amqp.Delivery)`; extract from delivery headers;
  start `rmq_rpc.process <type>` consumer span; pass `ctx` to the handler.
- NATS client: same with `nats_rpc.publish <handler>` and header injection.
- NATS server: same with `nats_rpc.process <handler>`.
- Update all 20 AMQP/NATS v1 handlers to accept `ctx` and stop calling `context.Background()`.
- Test first: integration test that asserts one trace contains both publish and process spans,
  plus a Gin server span and use-case span for HTTP.

### Task 4: Unit-of-work transaction

- Test first: integration test opens the unit of work, calls `Store`, returns an error, then
  asserts no history row exists. Also assert the success path commits.
- `internal/repo/persistent/translation/translation.go`: remove `Begin`/`Commit` from `Store`.

### Task 5: Verification

- `go mod tidy -diff`, `go build ./...`, `make test`, `make linter-golangci`,
  `make compose-up-integration-test`.
- README observability section reflects the new transport spans.
