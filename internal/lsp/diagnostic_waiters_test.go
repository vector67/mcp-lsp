package lsp

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/isaacphi/mcp-language-server/internal/protocol"
)

func newTestClient() *Client {
	return &Client{
		diagnosticWaiters: make(map[protocol.DocumentUri][]chan struct{}),
		diagnostics:       make(map[protocol.DocumentUri][]protocol.Diagnostic),
	}
}

func TestNotifyDiagnosticWaiters_ClosesChannelForMatchingURI(t *testing.T) {
	client := newTestClient()

	ch := make(chan struct{})
	uri := protocol.DocumentUri("file:///test.go")
	client.diagnosticWaiters[uri] = []chan struct{}{ch}

	client.notifyDiagnosticWaiters(uri)

	// Channel should be closed
	select {
	case <-ch:
		// good — channel was closed
	default:
		t.Fatal("expected channel to be closed, but it was not")
	}

	// Waiter slice should be cleaned up
	if len(client.diagnosticWaiters[uri]) != 0 {
		t.Fatalf("expected waiter slice to be empty, got %d", len(client.diagnosticWaiters[uri]))
	}
}

func TestNotifyDiagnosticWaiters_DoesNotCloseChannelForDifferentURI(t *testing.T) {
	client := newTestClient()

	ch := make(chan struct{})
	client.diagnosticWaiters["file:///a.go"] = []chan struct{}{ch}

	client.notifyDiagnosticWaiters("file:///b.go")

	select {
	case <-ch:
		t.Fatal("channel for a.go should not have been closed")
	default:
		// good — still open
	}
}

func TestNotifyDiagnosticWaiters_NoWaiters(t *testing.T) {
	client := newTestClient()
	// Should not panic
	client.notifyDiagnosticWaiters("file:///nonexistent.go")
}

func TestWaitForDiagnostics_WaitsForNotificationEvenWithCache(t *testing.T) {
	client := newTestClient()
	uri := protocol.DocumentUri("file:///test.go")

	// Pre-populate cache with stale data
	client.diagnostics[uri] = []protocol.Diagnostic{
		{Message: "stale"},
	}

	// Goroutine: send fresh diagnostics after 50ms
	go func() {
		time.Sleep(50 * time.Millisecond)
		client.diagnosticsMu.Lock()
		client.diagnostics[uri] = []protocol.Diagnostic{
			{Message: "fresh"},
		}
		client.diagnosticsMu.Unlock()
		client.notifyDiagnosticWaiters(uri)
	}()

	ctx := context.Background()
	diags, err := client.WaitForDiagnostics(ctx, uri, 5*time.Second)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diags) != 1 || diags[0].Message != "fresh" {
		t.Fatalf("expected fresh diagnostic, got %v", diags)
	}
}

func TestWaitForDiagnostics_WaitsForNotification(t *testing.T) {
	client := newTestClient()
	uri := protocol.DocumentUri("file:///test.go")

	// Goroutine: populate cache and notify after 50ms
	go func() {
		time.Sleep(50 * time.Millisecond)
		client.diagnosticsMu.Lock()
		client.diagnostics[uri] = []protocol.Diagnostic{
			{Message: "arrived"},
		}
		client.diagnosticsMu.Unlock()
		client.notifyDiagnosticWaiters(uri)
	}()

	ctx := context.Background()
	start := time.Now()
	diags, err := client.WaitForDiagnostics(ctx, uri, 5*time.Second)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diags) != 1 || diags[0].Message != "arrived" {
		t.Fatalf("expected notified diagnostic, got %v", diags)
	}
	// ~50ms for notification + ~1.5s settle period
	if elapsed < 30*time.Millisecond || elapsed > 3*time.Second {
		t.Fatalf("expected notification + settle (~1.5s), took %v", elapsed)
	}
}

func TestWaitForDiagnostics_RespectsContextCancellation(t *testing.T) {
	client := newTestClient()
	uri := protocol.DocumentUri("file:///test.go")

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	diags, err := client.WaitForDiagnostics(ctx, uri, 5*time.Second)
	elapsed := time.Since(start)

	if err != context.DeadlineExceeded {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
	if len(diags) != 0 {
		t.Fatalf("expected empty diagnostics, got %v", diags)
	}
	if elapsed < 30*time.Millisecond || elapsed > 500*time.Millisecond {
		t.Fatalf("expected ~50ms wait, took %v", elapsed)
	}

	// Waiter should be cleaned up
	client.diagnosticWaitersMu.Lock()
	remaining := len(client.diagnosticWaiters[uri])
	client.diagnosticWaitersMu.Unlock()
	if remaining != 0 {
		t.Fatalf("expected waiter to be cleaned up, got %d remaining", remaining)
	}
}

func TestWaitForDiagnostics_TimesOut(t *testing.T) {
	client := newTestClient()
	uri := protocol.DocumentUri("file:///test.go")

	ctx := context.Background()
	start := time.Now()
	diags, err := client.WaitForDiagnostics(ctx, uri, 100*time.Millisecond)
	elapsed := time.Since(start)

	// Timeout is best-effort, not a hard error
	if err != nil {
		t.Fatalf("expected nil error on timeout, got %v", err)
	}
	if len(diags) != 0 {
		t.Fatalf("expected empty diagnostics, got %v", diags)
	}
	if elapsed < 80*time.Millisecond || elapsed > 500*time.Millisecond {
		t.Fatalf("expected ~100ms wait, took %v", elapsed)
	}
}

func TestHandleDiagnostics_SignalsWaiters(t *testing.T) {
	client := newTestClient()
	uri := protocol.DocumentUri("file:///test.go")

	ch := make(chan struct{})
	client.diagnosticWaiters[uri] = []chan struct{}{ch}

	// Craft raw JSON publishDiagnostics params
	params := protocol.PublishDiagnosticsParams{
		URI: uri,
		Diagnostics: []protocol.Diagnostic{
			{Message: "from handler"},
		},
	}
	raw, _ := json.Marshal(params)

	HandleDiagnostics(client, raw)

	// Channel should be closed
	select {
	case <-ch:
		// good
	default:
		t.Fatal("expected HandleDiagnostics to signal waiters")
	}

	// Cache should be populated
	diags := client.GetFileDiagnostics(uri)
	if len(diags) != 1 || diags[0].Message != "from handler" {
		t.Fatalf("expected cache to be populated, got %v", diags)
	}
}

func TestWaitForDiagnostics_MultipleConcurrentWaiters(t *testing.T) {
	client := newTestClient()
	uri := protocol.DocumentUri("file:///test.go")

	channels := make([]chan struct{}, 3)
	for i := range channels {
		channels[i] = make(chan struct{})
	}
	client.diagnosticWaiters[uri] = channels

	client.notifyDiagnosticWaiters(uri)

	for i, ch := range channels {
		select {
		case <-ch:
			// good
		default:
			t.Fatalf("channel %d was not closed", i)
		}
	}
}

func TestWaitForDiagnostics_SettlesAfterNotification(t *testing.T) {
	client := newTestClient()
	uri := protocol.DocumentUri("file:///test.go")

	// Simulate gopls behavior: two rounds of diagnostics
	go func() {
		// First round: type checker finds 1 error
		time.Sleep(50 * time.Millisecond)
		client.diagnosticsMu.Lock()
		client.diagnostics[uri] = []protocol.Diagnostic{
			{Message: "error 1"},
		}
		client.diagnosticsMu.Unlock()
		client.notifyDiagnosticWaiters(uri)

		// Second round: analysis finds additional warning
		time.Sleep(100 * time.Millisecond)
		client.diagnosticsMu.Lock()
		client.diagnostics[uri] = []protocol.Diagnostic{
			{Message: "error 1"},
			{Message: "warning 1"},
		}
		client.diagnosticsMu.Unlock()
		client.notifyDiagnosticWaiters(uri)
	}()

	ctx := context.Background()
	diags, err := client.WaitForDiagnostics(ctx, uri, 5*time.Second)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diags) != 2 {
		t.Fatalf("expected 2 diagnostics after settle, got %d: %v", len(diags), diags)
	}
}

func TestWaitForDiagnostics_RaceCondition(t *testing.T) {
	client := newTestClient()
	uri := protocol.DocumentUri("file:///test.go")

	// Concurrent notify and wait — run with -race to verify no races.
	go func() {
		time.Sleep(10 * time.Millisecond)
		client.diagnosticsMu.Lock()
		client.diagnostics[uri] = []protocol.Diagnostic{
			{Message: "concurrent"},
		}
		client.diagnosticsMu.Unlock()
		client.notifyDiagnosticWaiters(uri)
	}()

	ctx := context.Background()
	diags, err := client.WaitForDiagnostics(ctx, uri, 5*time.Second)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diags) != 1 || diags[0].Message != "concurrent" {
		t.Fatalf("expected concurrent diagnostic, got %v", diags)
	}
}
