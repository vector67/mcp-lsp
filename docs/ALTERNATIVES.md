# Alternative Approaches (Non-isaacphi Ecosystem)

## Date: 2026-02-17

Other MCP servers that bridge LSP capabilities to AI agents, evaluated as alternatives to the isaacphi/virtuald Go-based approach.

---

## jonrad/lsp-mcp (TypeScript)

**URL:** https://github.com/jonrad/lsp-mcp
**Stars:** 165 | **Forks:** 24
**Language:** TypeScript
**Last commit:** 2025-03-31 (~25 commits over 5 weeks)
**License:** MIT
**Status:** "This is in a POC state" (author's words)

### Architecture

Dynamically generates MCP tools from the full LSP protocol JSON schema. Rather than hand-coding each tool, it parses the LSP spec and auto-generates tool definitions for nearly every LSP method.

### Unique Strengths

1. **Broadest LSP coverage** — Exposes ~30+ LSP methods as MCP tools
2. **`textDocument_codeAction`** — The only project that exposes code actions, which is the gateway to extract method, extract variable, move to file, quick fixes, auto-import, etc. (depends on what the underlying LSP server supports)
3. **`textDocument_formatting`** — Format document/selection
4. **Config file supports multiple LSPs** with extension-based routing
5. **`--methods` flag** to filter which tools are exposed (avoid flooding Claude with 30+ tools)
6. **Docker image** ships with TypeScript + Python LSP pre-installed
7. **Lazy LSP startup** — only spawns when first needed

### Critical Limitations

1. **Sends empty `ClientCapabilities{}`** — LSP servers may run in degraded mode and not advertise full features
2. **Files opened are never closed** — `textDocument/didClose` is never sent, memory grows unbounded
3. **No workspace scanning** — cross-file references may return incomplete results
4. **No reconnection** if LSP server crashes
5. **No tests**
6. **Claude Code integration not documented** (open issue #5)
7. **npm not published** — install via `git+https://github.com/jonrad/lsp-mcp`

### Installation

Docker:
```json
{
  "mcpServers": {
    "lsp": {
      "command": "docker",
      "args": ["run", "-i", "--rm",
               "-v", "/local/path:/workspace",
               "docker.io/jonrad/lsp-mcp:0.3.1"]
    }
  }
}
```

Config file for multiple LSPs:
```json
{
  "lsps": [
    {
      "id": "typescript",
      "extensions": ["ts", "tsx", "js", "jsx"],
      "languages": ["typescript", "javascript"],
      "command": "npx",
      "args": ["-y", "typescript-language-server", "--stdio"]
    },
    {
      "id": "python",
      "extensions": ["py"],
      "languages": ["python"],
      "command": "uvx",
      "args": ["--from", "python-lsp-server", "pylsp"]
    }
  ]
}
```

### Assessment

Most ambitious architecture. The `codeAction` bridge is uniquely valuable — it's the only path to extract/move refactoring without implementing it manually. But POC quality with fundamental gaps (empty capabilities, memory leaks). TypeScript means no single-binary distribution.

**Potential future integration:** If our Go-based fork gets stable, the `codeAction` approach could be back-ported as a new tool.

---

## sminnee/lsp-mcp (TypeScript)

**URL:** https://github.com/sminnee/lsp-mcp
**Stars:** 6 | **Forks:** 1
**Language:** TypeScript
**Total commits:** 2 (single day, Sep 6, 2025)
**License:** MIT
**Status:** Abandoned prototype

### Tools (5)

1. `rename_symbol` — Standard LSP rename
2. `find_references` — Opens every file before querying (perf disaster)
3. `extract_function` — LSP codeAction with refactor.extract.function filter
4. `move_function` — Naive brace-matching (NOT AST), stub import handling
5. `rename_file` — File/folder rename

### Hardcoded Language Support

| Language | Server |
|----------|--------|
| TypeScript/JS | typescript-language-server |
| Python | pyright (default) or ruff (lint mode) |
| Rust | rust-analyzer |
| Go | gopls |

### Critical Problems

- `move_function` import updating is a **stub** — moves break imports
- `move_function` uses brace-counting — breaks on strings/comments with braces
- `find_references` opens ALL project files first — performance disaster
- npm package name `lsp-mcp-server` taken over by different project (ProfessioneIT)
- No config file — all language servers hardcoded
- Custom hand-rolled LSP transport (fragile)
- Zero tests, zero community

### Assessment

The only project with `extract_function` and `move_function` tools, but they're the most broken features. Treat as reference for what tools to build, not as usable code.

---

## Comparison Matrix

| Feature | virtuald (Go) | jonrad (TS) | sminnee (TS) |
|---------|---------------|-------------|--------------|
| definition / go-to | Yes | Yes | No |
| references | Yes | Yes | Yes (slow) |
| rename | Yes | Yes | Yes |
| diagnostics | Yes | Yes | No |
| hover | Yes | Yes | No |
| callers/callees | Yes | No | No |
| code actions | **No** | **Yes** (unique) | Partial |
| extract function | **No** | Via codeAction | Yes (broken) |
| move function | **No** | Via codeAction | Yes (broken) |
| formatting | No | Yes | No |
| completion | No | Yes | No |
| multi-LSP | No (1 per instance) | Yes (config) | No (hardcoded) |
| file management | Proper open/close | Never closes | Opens everything |
| tests | Yes (integration) | None | None |
| single binary | Yes (Go) | No (Node.js) | No (Node.js) |
| maintenance | Slowing | Dormant | Abandoned |

---

## Recommendation

Build on **virtuald's fork** (Go ecosystem) for reliability, then consider back-porting **jonrad's `codeAction` approach** as a future enhancement to unlock extract/move refactoring.
