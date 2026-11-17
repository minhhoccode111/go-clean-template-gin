package integration_test

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	natsClient "github.com/minhhoccode111/go-clean-template-gin/pkg/nats/nats_rpc/client"
	rmqClient "github.com/minhhoccode111/go-clean-template-gin/pkg/rabbitmq/rmq_rpc/client"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/tracing"
)

const (
	jaegerURL         = "http://jaeger:16686"
	tracingService    = "go-clean-template-gin"
	tracingHTTPPath   = "/translation/history"
	tracingHTTPRoute  = "/v1" + tracingHTTPPath
	tracingHTTPMethod = http.MethodGet
)

// jaegerTracesResponse is the subset of the Jaeger v3 query API used here.
// The v3 API returns OTLP JSON, not the legacy Jaeger span format.
type jaegerTracesResponse struct {
	Result struct {
		ResourceSpans []struct {
			ScopeSpans []struct {
				Spans []struct {
					TraceID string `json:"traceId"`
					Name    string `json:"name"`
				} `json:"spans"`
			} `json:"scopeSpans"`
		} `json:"resourceSpans"`
	} `json:"result"`
}

// TestTracingTransportSpans verifies that each transport emits spans and that
// RPC callers and handlers share one trace through header propagation.
func TestTracingTransportSpans(t *testing.T) {
	startTestTracing(t)

	token := registerTracedUser(t)

	assertHTTPTrace(t, token)
	assertNATSTrace(t, token)
	assertRMQTrace(t, token)
}

// startTestTracing gives the test binary its own tracer provider and
// propagator: it acts as the RPC caller, so it must emit client spans and
// inject trace context.
func startTestTracing(t *testing.T) {
	t.Helper()

	shutdown, err := tracing.New(t.Context(), tracing.Config{
		Enabled:     true,
		ServiceName: "integration-test-client",
		Version:     "test",
		Endpoint:    "jaeger:4317",
		Insecure:    true,
		SampleRate:  1.0,
	})
	if err != nil {
		t.Fatalf("startTestTracing: tracing.New: %v", err)
	}

	t.Cleanup(func() {
		if serr := shutdown(context.Background()); serr != nil {
			t.Errorf("startTestTracing: tracing shutdown: %v", serr)
		}
	})
}

// registerTracedUser creates a unique user so the test can be re-run against a
// persistent database.
func registerTracedUser(t *testing.T) string {
	t.Helper()

	name := sanitizeTestName(t) + "_tracing_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	email := name + "@test.com"

	resp := registerUser(t, name, email, testPassword)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("registerTracedUser: register expected %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	resp.Body.Close()

	return loginUser(t, email, testPassword)
}

func assertHTTPTrace(t *testing.T, token string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
	defer cancel()

	resp, err := doAuthenticatedRequest(ctx, tracingHTTPMethod, basePathV1+tracingHTTPPath, http.NoBody, token)
	if err != nil {
		t.Fatalf("assertHTTPTrace: request failed: %v", err)
	}

	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("assertHTTPTrace: expected %d, got %d", http.StatusOK, resp.StatusCode)
	}

	waitForTrace(t, "HTTP history",
		tracingHTTPMethod+" "+tracingHTTPRoute,
		"TranslationUseCase.History",
	)
}

func assertNATSTrace(t *testing.T, token string) {
	t.Helper()

	client, err := natsClient.New(natsURL, rpcServerExchange)
	if err != nil {
		t.Fatalf("assertNATSTrace: natsClient.New: %v", err)
	}

	defer func() {
		if serr := client.Shutdown(); serr != nil {
			t.Fatalf("assertNATSTrace: client.Shutdown: %v", serr)
		}
	}()

	var history struct {
		History []any `json:"history"`
	}

	if err := client.RemoteCall("v1.translation.getHistory", authenticatedPayload(token, nil), &history); err != nil {
		t.Fatalf("assertNATSTrace: RemoteCall: %v", err)
	}

	waitForTrace(t, "NATS history",
		"nats_rpc.publish v1.translation.getHistory",
		"nats_rpc.process v1.translation.getHistory",
	)
}

func assertRMQTrace(t *testing.T, token string) {
	t.Helper()

	client, err := rmqClient.New(rmqURL, rpcServerExchange, rpcClientExchange)
	if err != nil {
		t.Fatalf("assertRMQTrace: rmqClient.New: %v", err)
	}

	defer func() {
		if serr := client.Shutdown(); serr != nil {
			t.Fatalf("assertRMQTrace: client.Shutdown: %v", serr)
		}
	}()

	var history struct {
		History []any `json:"history"`
	}

	if err := client.RemoteCall("v1.translation.getHistory", authenticatedPayload(token, nil), &history); err != nil {
		t.Fatalf("assertRMQTrace: RemoteCall: %v", err)
	}

	waitForTrace(t, "AMQP history",
		"rmq_rpc.publish v1.translation.getHistory",
		"rmq_rpc.process v1.translation.getHistory",
	)
}

// fetchTraceSpanNames groups span names by trace ID.
func fetchTraceSpanNames(t *testing.T) map[string]map[string]bool {
	t.Helper()

	query := url.Values{}
	query.Set("query.service_name", tracingService)
	query.Set("query.start_time_min", time.Now().Add(-10*time.Minute).UTC().Format(time.RFC3339Nano))
	query.Set("query.start_time_max", time.Now().Add(time.Minute).UTC().Format(time.RFC3339Nano))
	query.Set("query.num_traces", "100")

	ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, jaegerURL+"/api/v3/traces?"+query.Encode(), http.NoBody)
	if err != nil {
		t.Fatalf("fetchTraceSpanNames: http.NewRequestWithContext: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("fetchTraceSpanNames: request failed: %v", err)
	}

	defer resp.Body.Close()

	// Jaeger answers 404 for "no traces found"; keep polling until the batch
	// exporter flushes instead of failing the first empty query.
	if resp.StatusCode == http.StatusNotFound {
		return map[string]map[string]bool{}
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("fetchTraceSpanNames: expected %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body := parseJSON[jaegerTracesResponse](t, resp)
	traces := make(map[string]map[string]bool)

	for _, resource := range body.Result.ResourceSpans {
		for _, scope := range resource.ScopeSpans {
			for _, span := range scope.Spans {
				if traces[span.TraceID] == nil {
					traces[span.TraceID] = make(map[string]bool)
				}

				traces[span.TraceID][span.Name] = true
			}
		}
	}

	return traces
}

// waitForTrace polls Jaeger until one trace contains every wanted span name.
func waitForTrace(t *testing.T, label string, wanted ...string) {
	t.Helper()

	deadline := time.Now().Add(60 * time.Second)

	for {
		for _, spans := range fetchTraceSpanNames(t) {
			if traceHasSpans(spans, wanted) {
				return
			}
		}

		if time.Now().After(deadline) {
			t.Fatalf("waitForTrace: no %s trace contained all spans %v", label, wanted)
		}

		time.Sleep(2 * time.Second)
	}
}

func traceHasSpans(spans map[string]bool, wanted []string) bool {
	for _, name := range wanted {
		if !spans[name] {
			return false
		}
	}

	return true
}
