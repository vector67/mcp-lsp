package lsp

import (
	"context"
	"testing"
)

func TestNotifyChange_UnopenedFile_ReturnsNil(t *testing.T) {
	client := &Client{
		openFiles: make(map[string]*OpenFileInfo),
	}

	// NotifyChange for a file that isn't in openFiles should return nil
	// (skip silently) rather than returning an error.
	err := client.NotifyChange(context.Background(), "/some/unopened/file.go")
	if err != nil {
		t.Fatalf("expected nil error for unopened file, got: %v", err)
	}
}
