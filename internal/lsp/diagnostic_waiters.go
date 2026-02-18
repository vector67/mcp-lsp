package lsp

import (
	"context"
	"time"

	"github.com/isaacphi/mcp-language-server/internal/protocol"
)

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
		// First notification received — settle to catch additional rounds
		return c.settleDiagnostics(ctx, uri)
	case <-ctx.Done():
		return c.GetFileDiagnostics(uri), ctx.Err()
	case <-time.After(timeout):
		return c.GetFileDiagnostics(uri), nil
	}
}

// settleDiagnostics waits for additional diagnostic notifications to stop
// arriving. Each new notification resets the settle timer. Returns when
// no new notifications arrive within diagnosticSettleTime or the context
// is cancelled.
func (c *Client) settleDiagnostics(ctx context.Context, uri protocol.DocumentUri) ([]protocol.Diagnostic, error) {
	timer := time.NewTimer(diagnosticSettleTime)
	defer timer.Stop()

	for {
		ch := make(chan struct{})
		c.diagnosticWaitersMu.Lock()
		c.diagnosticWaiters[uri] = append(c.diagnosticWaiters[uri], ch)
		c.diagnosticWaitersMu.Unlock()

		select {
		case <-ch:
			// New notification — reset settle timer
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(diagnosticSettleTime)
			continue
		case <-timer.C:
			c.removeWaiter(uri, ch)
			return c.GetFileDiagnostics(uri), nil
		case <-ctx.Done():
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
