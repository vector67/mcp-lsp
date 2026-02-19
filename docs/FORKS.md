# Fork Analysis

## Date: 2026-02-17

Detailed analysis of all notable forks of vector67/mcp-language-server.

---

## virtuald/mcp-language-server (CHOSEN BASE)

**URL:** https://github.com/virtuald/mcp-language-server
**Author:** Dustin Spicuzza (virtuald)
**Profile:** http://www.virtualroadside.com/ | GitHub since 2011 | 212 public repos
**Stars:** 26 | **Forks of fork:** 10
**Relationship to upstream:** 29 commits ahead, 0 behind
**Active development period:** June 19 – July 10, 2025 (3 weeks)
**Last human commit:** 2025-07-10 (mcp-go v0.33.0 upgrade)
**Last push:** 2026-02-01 (Dependabot)

### Commit History (virtuald-only, chronological)

```
2025-06-19  Don't send multiple 'initialized' notifications (PR #4)
2025-06-19  Bump golang.org/x/text 0.25.0 -> 0.26.0 (PR #1)
2025-06-19  Fix golang unit tests (PR #3)
2025-06-19  Bump rust timeout
2025-06-19  Bump mcp-go 0.25.0 -> 0.32.0 (PR #5)
2025-06-19  Fix compilation errors in upgraded mcp-go
2025-06-06  Add hack to remove 'struct ' prefix from symbols (PR #6)
2025-06-06  Add -open flag for initial workspace file opening (PR #8)
2025-05-29  Add callers and callees tools (PR #9)
2025-06-20  Add content tool (PR #10)
2025-06-20  definition: match workspace symbols too (PR #11)
2025-06-10  clangd: use smallest matching symbol (PR #7)
2025-06-30  Remove duplicated glob matching logic (PR #12)
2025-06-30  Retry hover requests (PR #13)
2025-06-30  Update rust tests
2025-07-09  Fix nil dereference in content tool
2025-07-10  Upgrade mcp-go to v0.33.0
```

### Files Modified by virtuald (vs upstream)

**Core source:**
- `main.go` — Added `-open` flag, `StringArrayFlag` type, `openInitialFiles()`
- `tools.go` — Added callers, callees, content tool registrations

**LSP internals:**
- `internal/lsp/client.go` — Fixed duplicate initialized, modified initialization
- `internal/lsp/transport.go` — Minor changes

**Tool implementations:**
- `internal/tools/definition.go` — Workspace symbol matching, smallest symbol for clangd
- `internal/tools/hover.go` — Retry logic
- `internal/tools/references.go` — Modified
- `internal/tools/lsp-utilities.go` — Modified
- `internal/tools/utilities.go` — Modified

**New files:**
- `internal/tools/call_hierarchy.go` — callers/callees implementation
- `internal/tools/content.go` — content tool implementation

**Watcher:**
- `internal/watcher/watcher.go` — Removed duplicated glob logic

**Dependencies:**
- `go.mod`, `go.sum` — mcp-go v0.33.0, golang.org/x/text v0.26.0

**Tests/snapshots:**
- Many integration test files added/modified

### Open PRs on virtuald's fork

| PR | Title | Author | Date |
|----|-------|--------|------|
| #31 | Add HTTP support | bclermont | 2025-11-13 |
| #21 | Bump testify 1.10.0 -> 1.11.1 | Dependabot | 2025-09-01 |
| #23 | Bump actions/setup-python v5 -> v6 | Dependabot | 2025-10-01 |
| #24 | Bump actions/setup-go v5 -> v6 | Dependabot | 2025-10-01 |
| #28 | Bump actions/setup-node v4 -> v6 | Dependabot | 2025-11-01 |
| #34 | Bump actions/checkout v4 -> v6 | Dependabot | 2025-12-01 |
| #35 | Bump mcp-go v0.33.0 -> v0.43.2 | Dependabot | 2026-01-01 |
| #37 | Bump doublestar/v4 v4.8.1 -> v4.10.0 | Dependabot | 2026-02-01 |
| #38 | Bump golang.org/x/text v0.26.0 -> v0.33.0 | Dependabot | 2026-02-01 |

### Open Issues

| # | Title | Author | Date | Comments |
|---|-------|--------|------|----------|
| #18 | User feedback (from Claude itself) | creatorrr | 2025-08-14 | 2 |

Issue #18 contains real user testing feedback:
> "Symbol Renaming - The rename_symbol tool is genuinely excellent! It renamed health to vitality across all occurrences correctly. This is better than find-and-replace."
> "Definition Lookup - Works reliably for finding where symbols are defined, showing full implementation."
> "Editing with Line Numbers - The edit_file tool with line-based editing is cleaner than string-based Edit tool for some use cases."

### Forks of virtuald's fork

| Fork | Pushed | Stars |
|------|--------|-------|
| kym6464/mcp-language-server | 2025-06-20 | 0 |
| rkunnamp/mcp-language-server | 2025-06-24 | 0 |
| svenna/mcp-language-server | 2025-08-07 | 0 |
| noamsto/mcp-language-server | 2025-07-20 | 0 |
| maharjun/mcp-language-server | 2025-07-30 | 0 |
| titanous/mcp-language-server | 2025-08-08 | 0 |
| korri123/mcp-language-server | 2025-09-17 | 0 |
| asimihsan/mcp-language-server | 2025-10-09 | 0 |
| transform-ia/mcp-language-server | 2025-11-15 | 0 |
| p3psi-boo/mcp-language-server | 2025-12-18 | 0 |

---

## steffen-heil-secforge/mcp-language-server

**URL:** https://github.com/steffen-heil-secforge/mcp-language-server
**Author:** Steffen Heil (secforge GmbH)
**Stars:** 0 | **Forks:** 0
**Relationship to upstream:** main branch is identical to upstream
**Active work:** Feature branches only (PR workflow)
**Last push:** 2025-11-05

### Branches

- `main` — Identical to upstream
- `feat/multi-lsp-session-support` — 1 commit ahead (PR #97)
- `feat/lsp-glob-matching` — 1 commit ahead (PR #99)
- Plus 7 Dependabot branches

### Assessment

Pure PR workflow fork. Steffen never merged his own work into his own main — he was contributing upstream, not building a personal fork. When his PRs got no response, he moved on to other MCP projects (chrome-devtools-mcp, ms-365-mcp-server, both active as of Feb 2026).

---

## t3ta/mcp-language-server

**URL:** https://github.com/t3ta/mcp-language-server
**Stars:** 3 | **Forks:** 0
**Active development:** 2 days (Mar 26-28, 2025)
**Self-described:** "Pre-beta Quality"

### Unique feature

Multi-LSP via config file:
```json
{
  "workspaceDir": "/path/to/project",
  "languageServers": [
    {
      "language": "typescript",
      "command": "typescript-language-server",
      "args": ["--stdio"],
      "extensions": [".ts", ".tsx"]
    },
    {
      "language": "go",
      "command": "gopls",
      "args": [],
      "extensions": [".go"]
    }
  ]
}
```

### Tools (8)

`read_definition`, `find_references`, `rename_symbol`, `get_diagnostics`, `get_codelens`, `execute_codelens`, `apply_text_edit`, `find_symbols`

### Problems

- `find_symbols` queries only one LSP server arbitrarily
- Workspace watcher only uses first configured server
- No community
- Abandoned

### Assessment

Interesting multi-LSP config approach but half-baked. Not worth building on.

---

## Other Notable Forks of isaacphi (all 0 stars)

None show significant independent development beyond the three analyzed above.
