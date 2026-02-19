package watcher

import (
	"testing"

	"github.com/vector67/mcp-language-server/internal/protocol"
)

func TestMatchesPattern_BracedGlob(t *testing.T) {
	w := &WorkspaceWatcher{}

	// gopls registers patterns like **/*.{go,mod,sum,work}
	pattern := protocol.GlobPattern{
		Value: "**/*.{go,mod,sum,work}",
	}

	tests := []struct {
		path string
		want bool
	}{
		{"/workspace/main.go", true},
		{"/workspace/internal/foo/bar.go", true},
		{"/workspace/go.mod", true},
		{"/workspace/go.sum", true},
		{"/workspace/go.work", true},
		{"/workspace/readme.txt", false},
		{"/workspace/main.rs", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := w.matchesPattern(tt.path, pattern)
			if got != tt.want {
				t.Errorf("matchesPattern(%q, %q) = %v, want %v", tt.path, "**/*.{go,mod,sum,work}", got, tt.want)
			}
		})
	}
}

func TestMatchesPattern_SimpleGlob(t *testing.T) {
	w := &WorkspaceWatcher{}

	// Simple non-braced pattern should still work
	pattern := protocol.GlobPattern{
		Value: "**/*.go",
	}

	tests := []struct {
		path string
		want bool
	}{
		{"/workspace/main.go", true},
		{"/workspace/internal/foo.go", true},
		{"/workspace/readme.txt", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := w.matchesPattern(tt.path, pattern)
			if got != tt.want {
				t.Errorf("matchesPattern(%q, %q) = %v, want %v", tt.path, "**/*.go", got, tt.want)
			}
		})
	}
}
