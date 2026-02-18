package lsp

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ServerDetection struct {
	Command string
	Args    []string
}

type detectionRule struct {
	marker  string
	command string
	args    []string
}

var detectionRules = []detectionRule{
	{"go.mod", "gopls", nil},
	{"go.sum", "gopls", nil},
	{"Cargo.toml", "rust-analyzer", nil},
	{"tsconfig.json", "typescript-language-server", []string{"--stdio"}},
	{"package.json", "typescript-language-server", []string{"--stdio"}},
	{"pyproject.toml", "pyright-langserver", []string{"--stdio"}},
	{"setup.py", "pyright-langserver", []string{"--stdio"}},
	{"requirements.txt", "pyright-langserver", []string{"--stdio"}},
	{"compile_commands.json", "clangd", nil},
	{"CMakeLists.txt", "clangd", nil},
}

func DetectServer(workspaceDir string) (*ServerDetection, error) {
	for _, rule := range detectionRules {
		path := filepath.Join(workspaceDir, rule.marker)
		if _, err := os.Stat(path); err == nil {
			return &ServerDetection{
				Command: rule.command,
				Args:    rule.args,
			}, nil
		}
	}

	var markers []string
	for _, rule := range detectionRules {
		markers = append(markers, rule.marker)
	}

	return nil, fmt.Errorf(
		"could not detect LSP server for workspace %s: no recognized project files found. "+
			"Looked for: %s. Use --lsp to specify the server manually",
		workspaceDir,
		strings.Join(markers, ", "),
	)
}
