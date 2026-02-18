# Future Roadmap & Ideas

## Date: 2026-02-17

Ideas and priorities for after the initial fork + PR merge is complete.

---

## Priority 1: Correctness

### Fix cross-file rename (upstream #104)
The `rename_symbol` tool only renames in the defining file, not across the project. This is likely caused by missing workspace edit client capabilities in the LSP initialization. Need to:
1. Audit the `ClientCapabilities` sent during `initialize`
2. Ensure `workspace.workspaceEdit.documentChanges` is declared
3. Test with gopls, pyright, typescript-language-server, rust-analyzer

### Fix stale diagnostics
Even with the glob fix (PR #99), verify that `textDocument/didChange` notifications are properly sent for all file modifications. The file watcher needs to correctly detect changes made by external tools (git, other editors).

### Fix CI
Upstream CI is broken (issue #78). Set up working GitHub Actions with:
- Go build + test
- Integration tests against gopls, pyright, typescript-language-server, rust-analyzer
- Linting (golangci-lint)

---

## Priority 2: Reliability

### Resource management
- Fix "too many open files" (upstream #83) for large codebases
- Implement proper file open/close lifecycle tracking
- Add connection health monitoring and LSP server restart on crash
- Consider the headless mode (PR #116) as the primary solution for large projects

### Error handling
- Graceful degradation when LSP server doesn't support a capability
- Better error messages when LSP server is not installed or fails to start
- Timeout handling for slow LSP initialization (especially clangd)

---

## Priority 3: New Capabilities

### textDocument/codeAction support (the big one)
This is the gateway to extract method, extract variable, move to file, auto-import, quick fixes, etc. jonrad/lsp-mcp is the only project that exposes this.

Implementation approach:
1. Add `code_actions` tool that takes a file + range and returns available actions
2. Add `apply_code_action` tool that executes a selected action
3. Handle `workspace/applyEdit` requests from the LSP server
4. The available actions depend entirely on the underlying LSP server:
   - TypeScript: extract function, extract constant, move to file, organize imports
   - Python (pyright): limited code actions
   - Go (gopls): extract function, extract variable, fill struct
   - Rust (rust-analyzer): extract function, extract variable, inline variable

### textDocument/formatting
Simple but useful. Add `format_document` and `format_range` tools.

### Workspace symbols
Add `workspace_symbols` tool for searching symbols across the project (not just within a file).

---

## Priority 4: Architecture

### Multi-LSP support
The pragmatic approach (multiple instances) works but is clunky. Options:
1. **Config-based routing** (like t3ta's approach) — one MCP server, multiple LSP servers, route by file extension
2. **Session management** (like steffen's PR #97) — dynamic start/stop of LSP servers
3. **Keep it simple** — document how to configure multiple instances in Claude Code settings

Steffen's PR #97 (+3,361 lines) is the most comprehensive attempt but is architecturally incompatible with the current codebase. If multi-LSP becomes critical, consider reaching out to him about reimplementing on the current architecture.

### Headless mode refinement
PR #116 adds basic headless mode. Future enhancements:
- Auto-discovery of running LSP servers
- Protocol negotiation (some LSPs use different JSON-RPC transports)
- Shared state management with IDE

---

## Priority 5: Distribution & Docs

### Releases
- Set up GoReleaser for automated binary builds (Linux, macOS, Windows)
- Publish to Homebrew
- Nix flake (PR #70)

### Documentation
- Complete Claude Code setup guide (all supported languages)
- Troubleshooting guide (common issues: PATH not set, LSP not found, stale diagnostics)
- Performance tuning guide (large codebases, multiple languages)

### Testing
- Expand integration test coverage to all supported languages
- Add regression tests for known bugs (cross-file rename, stale diagnostics)
- Benchmark suite for startup time and response latency

---

## People to Potentially Engage

### Steffen Heil (steffen-heil-secforge)
- Wrote the multi-LSP PR and glob matching PR
- Active on other MCP projects (chrome-devtools-mcp, ms-365-mcp-server)
- Works at secforge GmbH
- Might be interested in contributing if there's actually a maintained fork to contribute to

### Rick Batka (rickbatka)
- Wrote the headless mode PR (Feb 2026, very recent)
- May still be actively working on it
- Could be interested in collaborating

### Dustin Spicuzza (virtuald)
- Original fork author
- May appreciate someone carrying his work forward
- Could provide context on design decisions

### creatorrr
- User who tested virtuald's fork extensively (issue #18)
- Could be an early adopter/tester

### bclermont
- Submitted HTTP support PR (#31) to virtuald's fork
- Interested in alternative transport mechanisms
