package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/isaacphi/mcp-language-server/internal/lsp"
)

func TestDetectServer(t *testing.T) {
	tests := []struct {
		name        string
		markers     []string
		wantCommand string
		wantArgs    []string
		wantErr     bool
	}{
		{
			name:        "go.mod detects gopls",
			markers:     []string{"go.mod"},
			wantCommand: "gopls",
			wantArgs:    nil,
		},
		{
			name:        "go.sum detects gopls",
			markers:     []string{"go.sum"},
			wantCommand: "gopls",
			wantArgs:    nil,
		},
		{
			name:        "Cargo.toml detects rust-analyzer",
			markers:     []string{"Cargo.toml"},
			wantCommand: "rust-analyzer",
			wantArgs:    nil,
		},
		{
			name:        "tsconfig.json detects typescript-language-server",
			markers:     []string{"tsconfig.json"},
			wantCommand: "typescript-language-server",
			wantArgs:    []string{"--stdio"},
		},
		{
			name:        "package.json detects typescript-language-server",
			markers:     []string{"package.json"},
			wantCommand: "typescript-language-server",
			wantArgs:    []string{"--stdio"},
		},
		{
			name:        "pyproject.toml detects pyright-langserver",
			markers:     []string{"pyproject.toml"},
			wantCommand: "pyright-langserver",
			wantArgs:    []string{"--stdio"},
		},
		{
			name:        "setup.py detects pyright-langserver",
			markers:     []string{"setup.py"},
			wantCommand: "pyright-langserver",
			wantArgs:    []string{"--stdio"},
		},
		{
			name:        "requirements.txt detects pyright-langserver",
			markers:     []string{"requirements.txt"},
			wantCommand: "pyright-langserver",
			wantArgs:    []string{"--stdio"},
		},
		{
			name:        "compile_commands.json detects clangd",
			markers:     []string{"compile_commands.json"},
			wantCommand: "clangd",
			wantArgs:    nil,
		},
		{
			name:        "CMakeLists.txt detects clangd",
			markers:     []string{"CMakeLists.txt"},
			wantCommand: "clangd",
			wantArgs:    nil,
		},
		{
			name:    "empty directory returns error",
			markers: nil,
			wantErr: true,
		},
		{
			name:        "priority: go.mod + package.json → gopls",
			markers:     []string{"go.mod", "package.json"},
			wantCommand: "gopls",
			wantArgs:    nil,
		},
		{
			name:        "priority: tsconfig.json + pyproject.toml → typescript-language-server",
			markers:     []string{"tsconfig.json", "pyproject.toml"},
			wantCommand: "typescript-language-server",
			wantArgs:    []string{"--stdio"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()

			for _, marker := range tt.markers {
				path := filepath.Join(dir, marker)
				if err := os.WriteFile(path, []byte{}, 0644); err != nil {
					t.Fatalf("failed to create marker file %s: %v", marker, err)
				}
			}

			result, err := lsp.DetectServer(dir)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Command != tt.wantCommand {
				t.Errorf("command = %q, want %q", result.Command, tt.wantCommand)
			}

			if len(result.Args) != len(tt.wantArgs) {
				t.Errorf("args = %v, want %v", result.Args, tt.wantArgs)
			} else {
				for i := range result.Args {
					if result.Args[i] != tt.wantArgs[i] {
						t.Errorf("args[%d] = %q, want %q", i, result.Args[i], tt.wantArgs[i])
					}
				}
			}
		})
	}
}
