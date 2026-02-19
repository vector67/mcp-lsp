package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vector67/mcp-language-server/internal/protocol"
)

func TestAffectedFiles_FromChanges(t *testing.T) {
	edit := protocol.WorkspaceEdit{
		Changes: map[protocol.DocumentUri][]protocol.TextEdit{
			"file:///workspace/main.go": {
				{Range: protocol.Range{}, NewText: "new"},
			},
			"file:///workspace/util.go": {
				{Range: protocol.Range{}, NewText: "new"},
			},
		},
	}

	files := AffectedFiles(edit)

	assert.Len(t, files, 2)
	assert.Contains(t, files, "/workspace/main.go")
	assert.Contains(t, files, "/workspace/util.go")
}

func TestAffectedFiles_FromDocumentChanges(t *testing.T) {
	edit := protocol.WorkspaceEdit{
		DocumentChanges: []protocol.DocumentChange{
			{
				TextDocumentEdit: &protocol.TextDocumentEdit{
					TextDocument: protocol.OptionalVersionedTextDocumentIdentifier{
						TextDocumentIdentifier: protocol.TextDocumentIdentifier{
							URI: "file:///workspace/foo.go",
						},
					},
				},
			},
			{
				TextDocumentEdit: &protocol.TextDocumentEdit{
					TextDocument: protocol.OptionalVersionedTextDocumentIdentifier{
						TextDocumentIdentifier: protocol.TextDocumentIdentifier{
							URI: "file:///workspace/bar.go",
						},
					},
				},
			},
		},
	}

	files := AffectedFiles(edit)

	assert.Len(t, files, 2)
	assert.Contains(t, files, "/workspace/foo.go")
	assert.Contains(t, files, "/workspace/bar.go")
}

func TestAffectedFiles_BothSources_Deduplicates(t *testing.T) {
	edit := protocol.WorkspaceEdit{
		Changes: map[protocol.DocumentUri][]protocol.TextEdit{
			"file:///workspace/main.go": {
				{Range: protocol.Range{}, NewText: "new"},
			},
		},
		DocumentChanges: []protocol.DocumentChange{
			{
				TextDocumentEdit: &protocol.TextDocumentEdit{
					TextDocument: protocol.OptionalVersionedTextDocumentIdentifier{
						TextDocumentIdentifier: protocol.TextDocumentIdentifier{
							URI: "file:///workspace/main.go",
						},
					},
				},
			},
		},
	}

	files := AffectedFiles(edit)

	assert.Len(t, files, 1)
	assert.Contains(t, files, "/workspace/main.go")
}

func TestAffectedFiles_Empty(t *testing.T) {
	edit := protocol.WorkspaceEdit{}

	files := AffectedFiles(edit)

	assert.Empty(t, files)
}

func TestAffectedFiles_SkipsNilTextDocumentEdit(t *testing.T) {
	edit := protocol.WorkspaceEdit{
		DocumentChanges: []protocol.DocumentChange{
			{
				TextDocumentEdit: nil,
			},
		},
	}

	files := AffectedFiles(edit)

	assert.Empty(t, files)
}
