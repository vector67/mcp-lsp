package tools

import (
	"strings"

	"github.com/vector67/mcp-language-server/internal/protocol"
)

// AffectedFiles extracts unique file paths from a workspace edit.
// It collects paths from both the Changes and DocumentChanges fields.
func AffectedFiles(edit protocol.WorkspaceEdit) []string {
	seen := make(map[string]struct{})

	for uri := range edit.Changes {
		path := strings.TrimPrefix(string(uri), "file://")
		seen[path] = struct{}{}
	}

	for _, change := range edit.DocumentChanges {
		if change.TextDocumentEdit != nil {
			path := strings.TrimPrefix(string(change.TextDocumentEdit.TextDocument.URI), "file://")
			seen[path] = struct{}{}
		}
	}

	files := make([]string, 0, len(seen))
	for path := range seen {
		files = append(files, path)
	}
	return files
}
