package lsp

import (
	"context"
	"time"

	"github.com/vector67/mcp-language-server/internal/logging"
	"github.com/vector67/mcp-language-server/internal/protocol"
)

var diagLogger = logging.NewLogger(logging.LSP)

// diagnosticSettleTime is how long to wait after each notification for
// additional diagnostics. LSP servers like gopls send diagnostics in
// multiple rounds (e.g., type-check first, then analysis). This settle
// period catches subsequent rounds without a fixed sleep.
const diagnosticSettleTime = 1500 * time.Millisecond

// notifyDiagnosticWaiters closes all waiter channels for the given URI
// and removes them from the map.
func (c *Client) notifyDiagnosticWaiters(uri protocol.DocumentUri) {
	c.diagnosticWaitersMu.Lock()
	defer c.diagnosticWaitersMu.Unlock()

	waiterCount := len(c.diagnosticWaiters[uri])
	if waiterCount > 0 {
		diagLogger.Debug("notifyDiagnosticWaiters: closing %d waiter(s) for %s", waiterCount, uri)
	}
	for _, ch := range c.diagnosticWaiters[uri] {
		close(ch)
	}
	delete(c.diagnosticWaiters, uri)
}

// WaitForDiagnostics waits for a publishDiagnostics notification for the
// given URI, then settles briefly to catch additional rounds (LSP servers
// like gopls send diagnostics incrementally: type-check first, then
// analysis). Returns on context cancellation or timeout.
func (c *Client) WaitForDiagnostics(ctx context.Context, uri protocol.DocumentUri, timeout time.Duration) ([]protocol.Diagnostic, error) {
	diagLogger.Debug("WaitForDiagnostics: waiting for first notification for %s (timeout=%v)", uri, timeout)
	start := time.Now()
	ch := make(chan struct{})

	// 1. Register waiter
	c.diagnosticWaitersMu.Lock()
	c.diagnosticWaiters[uri] = append(c.diagnosticWaiters[uri], ch)
	c.diagnosticWaitersMu.Unlock()

	// 2. Cleanup on exit
	defer c.removeWaiter(uri, ch)

	// 3. Wait for first notification
	select {
	case <-ch:
		diagLogger.Debug("WaitForDiagnostics: first notification received for %s after %v, entering settle", uri, time.Since(start))
		diags, err := c.settleDiagnostics(ctx, uri)
		diagLogger.Debug("WaitForDiagnostics: settled for %s after %v total, returning %d diagnostics", uri, time.Since(start), len(diags))
		return diags, err
	case <-ctx.Done():
		diags := c.GetFileDiagnostics(uri)
		diagLogger.Debug("WaitForDiagnostics: context cancelled for %s after %v, returning %d cached diagnostics", uri, time.Since(start), len(diags))
		return diags, ctx.Err()
	case <-time.After(timeout):
		diags := c.GetFileDiagnostics(uri)
		diagLogger.Debug("WaitForDiagnostics: timed out for %s after %v, returning %d cached diagnostics", uri, time.Since(start), len(diags))
		return diags, nil
	}
}

// settleDiagnostics waits for additional diagnostic notifications to stop
// arriving. Each new notification resets the settle timer. Returns when
// no new notifications arrive within diagnosticSettleTime or the context
// is cancelled.
func (c *Client) settleDiagnostics(ctx context.Context, uri protocol.DocumentUri) ([]protocol.Diagnostic, error) {
	timer := time.NewTimer(diagnosticSettleTime)
	defer timer.Stop()

	rounds := 0
	start := time.Now()
	for {
		ch := make(chan struct{})
		c.diagnosticWaitersMu.Lock()
		c.diagnosticWaiters[uri] = append(c.diagnosticWaiters[uri], ch)
		c.diagnosticWaitersMu.Unlock()

		select {
		case <-ch:
			rounds++
			diagLogger.Debug("settleDiagnostics: round %d for %s (elapsed=%v), resetting settle timer", rounds, uri, time.Since(start))
			// New notification — reset settle timer
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(diagnosticSettleTime)
			continue
		case <-timer.C:
			diagLogger.Debug("settleDiagnostics: settled for %s after %d rounds, %v elapsed", uri, rounds, time.Since(start))
			c.removeWaiter(uri, ch)
			return c.GetFileDiagnostics(uri), nil
		case <-ctx.Done():
			diagLogger.Debug("settleDiagnostics: context cancelled for %s after %d rounds, %v elapsed", uri, rounds, time.Since(start))
			c.removeWaiter(uri, ch)
			return c.GetFileDiagnostics(uri), ctx.Err()
		}
	}
}

// removeWaiter removes a specific channel from the waiters for a URI.
func (c *Client) removeWaiter(uri protocol.DocumentUri, ch chan struct{}) {
	c.diagnosticWaitersMu.Lock()
	defer c.diagnosticWaitersMu.Unlock()

	waiters := c.diagnosticWaiters[uri]
	for i, w := range waiters {
		if w == ch {
			c.diagnosticWaiters[uri] = append(waiters[:i], waiters[i+1:]...)
			break
		}
	}
	if len(c.diagnosticWaiters[uri]) == 0 {
		delete(c.diagnosticWaiters, uri)
	}
}
