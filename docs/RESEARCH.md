# MCP Language Server — Research & Analysis

## Date: 2026-02-17

This document captures all research into MCP servers that bridge LSP (Language Server Protocol) to MCP (Model Context Protocol), enabling AI coding agents like Claude Code to perform IDE-style refactoring operations.

---

## 1. Problem Statement

We want a CLI tool / MCP server that allows Claude Code to perform LSP-style refactoring operations (rename, extract method, find references, go to definition, etc.) similar to JetBrains IDEs. No production-ready solution exists. All projects in this space are early-stage or abandoned.

---

## 2. Landscape Overview

Five projects were evaluated:

| Project | Stars | Language | Last Human Commit | Status |
|---------|-------|----------|-------------------|--------|
| [vector67/mcp-language-server](https://github.com/vector67/mcp-language-server) | 1,447 | Go | May 2025 | Unmaintained (author confirmed) |
| [virtuald/mcp-language-server](https://github.com/virtuald/mcp-language-server) | 26 | Go | Jul 2025 | Best fork, also slowing |
| [jonrad/lsp-mcp](https://github.com/jonrad/lsp-mcp) | 165 | TypeScript | Mar 2025 | POC, dormant |
| [sminnee/lsp-mcp](https://github.com/sminnee/lsp-mcp) | 6 | TypeScript | Sep 2025 | Abandoned prototype |
| [t3ta/mcp-language-server](https://github.com/t3ta/mcp-language-server) | 3 | Go | Mar 2025 | Abandoned fork of isaacphi |

---

## 3. Detailed Project Analyses

### 3.1 vector67/mcp-language-server (upstream)

**Repository:** https://github.com/vector67/mcp-language-server
**Author:** Isaac Phi
**Stars:** 1,447 | **Forks:** 114 | **Contributors:** 4 (isaacphi: 98 commits, +3 minor)
**License:** BSD 3-Clause
**Created:** 2024-12-30
**Last release:** v0.1.1 (2025-05-16)
**Last human commit:** 2025-05-16
**Open issues:** 25 | **Closed issues:** 21
**Open PRs:** 22 (ZERO have been reviewed)

#### Tools (6)

1. `definition` — Retrieves complete source code where a symbol is defined
2. `references` — Locates all usages of a symbol across the codebase
3. `diagnostics` — Provides diagnostic information (warnings, errors) for a file
4. `hover` — Displays documentation, type hints, or hover info for a location
5. `rename_symbol` — Renames a symbol across a project
6. `edit_file` — Makes multiple text edits to a file based on line numbers

#### Supported Languages (documented)

| Language | LSP Server | Notes |
|----------|-----------|-------|
| Go | gopls | Best tested |
| Rust | rust-analyzer | Works well |
| Python | pyright-langserver | Needs `-- --stdio` |
| TypeScript | typescript-language-server | Needs `-- --stdio` |
| C/C++ | clangd | Requires compile_commands.json |

Community-attempted (varying success): Zig (zls), Java (jdtls, issue #101), Vue (broken, #85), PHP/Intelephense (broken, #85), C#/Omnisharp (partial, #73), Coq (open PR #109).

#### Critical Bugs

- **Issue #104:** `rename_symbol` only renames in the defining file, NOT across the project. The one refactoring tool is broken for its primary use case.
- **Issue #83:** "Too many open files" — file watcher exhausts file descriptors on large codebases.
- **Issue #54:** Diagnostics are stale — file watcher glob matching is broken, changes not detected. Requires MCP server restart.
- **Issue #86:** TypeScript references fail for class references (methods work).
- **Issue #85:** Vue and PHP language servers don't work at all.
- **Issue #60:** Unreliable diagnostics for TypeScript.
- **Issue #78:** CI pipeline is broken.

#### Maintenance Status

The author explicitly stated in **issue #92** (Dec 2025):
> "Unfortunately, I don't have too much time to maintain this repo at the moment."

He invited community forks and offered merge access to interested Go developers. No action was taken. 22 PRs sit unreviewed, the oldest from June 2025.

#### Installation

```bash
go install github.com/vector67/mcp-language-server@latest
```

Configuration for Claude Code (`.claude/settings.json`):
```json
{
  "mcpServers": {
    "language-server": {
      "command": "mcp-language-server",
      "args": ["--workspace", "/path/to/project/", "--lsp", "gopls"],
      "env": {
        "PATH": "/opt/homebrew/bin:/Users/you/go/bin",
        "GOPATH": "/Users/you/go"
      }
    }
  }
}
```

Note: Claude Desktop/Code don't inherit shell environment; PATH, GOPATH, GOCACHE, GOMODCACHE must be explicit.

---

### 3.2 virtuald/mcp-language-server (best fork)

**Repository:** https://github.com/virtuald/mcp-language-server
**Author:** Dustin Spicuzza — seasoned OSS developer (account since 2011, 212 public repos, maintains pyhcl with 340 stars)
**Stars:** 26 | **Forks:** 10
**Commits ahead of upstream:** 29 | **Behind:** 0
**Last commit:** 2025-07-10
**Last push activity:** 2026-02-01 (Dependabot)
**Open issues:** 1 real issue (#18 — user feedback)

#### Additional Tools (3, on top of upstream's 6 = 9 total)

7. `content` — Retrieves source code at a specific file/line location
8. `callers` — Call hierarchy incoming (who calls this function?)
9. `callees` — Call hierarchy outgoing (what does this function call?)

#### Additional Improvements

- **`-open` flag** — Triggers initial opening of workspace files for LSP servers like clangd that need it
- **Workspace symbol matching** in the `definition` tool
- **File watcher glob fix** (addresses upstream issue #54)
- **Hover retry logic** for flaky LSP responses
- **clangd-specific improvements** — Uses smallest matching symbol for definition results, strips `struct ` prefix from symbols
- **Fix for duplicate `initialized` notifications** (upstream PR #49)
- **Upgraded mcp-go to v0.33.0** (upstream is on v0.25.0)

#### Usage Signals

- 10 people have forked this fork
- Issue #18: Real user tested extensively, reported rename works great, definition is solid
- PR #31: External contributor submitted HTTP transport support
- Dependabot enabled and configured (ongoing maintenance infrastructure)

#### Assessment

Dustin was genuinely using this, likely with clangd (C/C++ — that's where most custom fixes are). He built what he needed over 3 weeks (June 19 – July 10, 2025) and stopped. Most "production-ready" option available, but drifting toward unmaintained.

---

### 3.3 jonrad/lsp-mcp

**Repository:** https://github.com/jonrad/lsp-mcp
**Author:** Jon Rad (solo developer)
**Stars:** 165 | **Forks:** 24
**Language:** TypeScript
**Last commit:** 2025-03-31 (~25 commits over 5 weeks)
**Open issues:** 5 (3 real, 2 spam)
**License:** MIT

**README explicitly states: "This is in a POC state."**

#### Tools

Dynamically generated from the full LSP protocol JSON schema. Exposes nearly every LSP method as an MCP tool:

**Navigation:** `textDocument_definition`, `textDocument_declaration`, `textDocument_typeDefinition`, `textDocument_implementation`, `textDocument_references`, `textDocument_documentSymbol`, `workspace_symbol`, `textDocument_hover`

**Refactoring:** `textDocument_rename`, `textDocument_prepareRename`, `textDocument_codeAction` (gateway to extract, move, quick fixes — depends on underlying LSP server), `workspace_applyEdit`, `workspace_executeCommand`

**Completion & Diagnostics:** `textDocument_completion`, `completionItem_resolve`, `textDocument_signatureHelp`, `textDocument_publishDiagnostics`

**Formatting:** `textDocument_formatting`, `textDocument_rangeFormatting`, `textDocument_onTypeFormatting`

**Plus:** Code lens, document links, folding ranges, document lifecycle, workspace notifications, and 2 custom tools (`lsp_info`, `file_contents_to_uri`).

**`--methods` flag** allows filtering which tools are exposed.

#### Key Limitations

1. Sends **empty `ClientCapabilities{}`** during initialization — LSP servers may run in degraded mode
2. Files opened are **never closed** — memory leak in long sessions
3. No workspace scanning — cross-file references may be incomplete
4. No reconnection if LSP server crashes
5. No tests
6. Claude Code integration not documented (open issue #5)

#### Installation

Docker (recommended):
```json
{
  "mcpServers": {
    "lsp": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "-v", "/local/path:/workspace", "docker.io/jonrad/lsp-mcp:0.3.1"]
    }
  }
}
```

Or npx from GitHub:
```json
{
  "mcpServers": {
    "lsp": {
      "command": "npx",
      "args": ["-y", "--silent", "git+https://github.com/jonrad/lsp-mcp",
               "--lsp", "npx -y typescript-language-server --stdio"]
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

#### Assessment

Most ambitious architecture. The `textDocument_codeAction` bridge is the only thing in this landscape that *could* do extract/move refactoring (if the LSP server supports it). But POC quality with real gaps. The TypeScript implementation also means it doesn't benefit from Go's single-binary distribution.

---

### 3.4 sminnee/lsp-mcp

**Repository:** https://github.com/sminnee/lsp-mcp
**Author:** Sam Minnee
**Stars:** 6 | **Forks:** 1
**Language:** TypeScript
**Total commits:** 2 (both on Sep 6, 2025 — "Initial commit" and "initial version")
**Open issues:** 0 | **License:** MIT

#### Tools (5)

1. `rename_symbol` — LSP textDocument/rename
2. `find_references` — LSP textDocument/references
3. `extract_function` — LSP textDocument/codeAction with refactor.extract.function filter
4. `move_function` — **Hybrid/manual implementation** using naive brace-matching (NOT AST)
5. `rename_file` — File/folder rename

#### Supported Languages (hardcoded)

TypeScript/JavaScript (typescript-language-server), Python (pyright or ruff), Rust (rust-analyzer), Go (gopls)

#### Critical Problems

- `move_function` uses brace-counting to find function boundaries — breaks on strings with braces, template literals, comments
- Import updating after `move_function` is a **stub** — moves break imports
- `find_references` opens **every file** in the project before querying — major perf issue
- Custom hand-rolled LSP transport instead of established libraries
- npm package name `lsp-mcp-server` has been **taken over by a different project** on npm
- No configuration file — language servers are hardcoded
- Zero tests, zero community

#### Assessment

Interesting proof-of-concept with unique features (extract_function, move_function) that are the most broken parts. Abandoned single-day prototype.

---

### 3.5 t3ta/mcp-language-server

**Repository:** https://github.com/t3ta/mcp-language-server
**Stars:** 3 | **Language:** Go
**Fork of:** vector67/mcp-language-server
**Total development:** 2 days (Mar 26-28, 2025)
**Self-described:** "Pre-beta Quality"

#### Unique Value

Multi-LSP in one process via config file (TypeScript + Python + Go simultaneously).

#### Tools (8)

`read_definition`, `find_references`, `rename_symbol`, `get_diagnostics`, `get_codelens`, `execute_codelens`, `apply_text_edit`, `find_symbols`

#### Problems

- `find_symbols` with workspace scope picks the first LSP client arbitrarily
- Workspace watcher only uses first configured language server
- `rename_symbol` likely has same single-file bug as upstream (#104)
- Zero community engagement

#### Assessment

The multi-LSP config approach is interesting but the implementation is half-baked. Upstream (isaacphi) was already exploring multi-LSP via issue #97.

---

## 4. Why Are None of These Maintained?

### The Structural Problem

1. **LSP servers are designed for editors, not headless tools.** They expect a long-lived client that manages file open/close lifecycle, sends `didChange` notifications, negotiates capabilities, and handles workspace indexing. Every project cuts corners here — leading to bugs that are architectural, not incidental.

2. **Bug surface area scales multiplicatively.** Each language server implements a different subset of LSP with different quirks. One maintainer can't keep up with (languages x features x edge cases).

3. **The hardest problems are the most important ones.** Cross-file rename is broken. Move function has stub import handling. These aren't things you forgot to finish; they're genuinely difficult requiring proper capability negotiation, workspace-wide edit handling, and deep per-server testing.

4. **Governance failure.** isaacphi has 1,447 stars, 22 unmerged PRs, and an author who said he doesn't have time. The community *wants* this tool but nobody wants to *maintain* it. steffen-heil-secforge wrote +3,361 lines with 50+ tests, got ghosted, and moved on.

5. **Moving target.** MCP is still young. AI tools keep getting better at code understanding natively. Solo developers don't want to invest months into something that might be obsoleted.

6. **The "cool demo" to "production tool" gap.** Getting rename working on a simple project takes an afternoon. Making it reliable across languages and large codebases takes months of unglamorous compatibility work.

---

## 5. Key People

### isaacphi (Isaac Phi)
- Upstream author, 98 commits
- Acknowledged he can't maintain it (issue #92, Dec 2025)
- Offered merge access, no takers

### virtuald (Dustin Spicuzza)
- Best fork maintainer
- Seasoned OSS dev (account since 2011, 212 repos, pyhcl with 340 stars)
- Built what he needed for clangd over 3 weeks (Jun-Jul 2025)
- 10 people forked his fork
- Last human commit: Jul 2025

### steffen-heil-secforge (Steffen Heil)
- Works at secforge GmbH (security company)
- Actively building MCP tools (chrome-devtools-mcp, ms-365-mcp-server, ios-simulator-mcp)
- Submitted the largest PR (#97, +3,361 lines, 50+ tests) — multi-LSP support
- Also submitted #99 (glob matching, +526 lines)
- Got zero response, moved on
- Last active: Feb 2026 (on other projects)

### jonrad (Jon Rad)
- Built the most architecturally ambitious approach (full LSP protocol bridge)
- Explicitly labeled it POC
- 5 weeks of work, then stopped
- 165 stars shows community interest

### rickbatka (Rick Batka)
- Submitted #116 (headless mode, +853 lines) — Feb 2026, very recent
- Allows connecting to an already-running IDE LSP server
- Still actively working on it

### bclermont
- Submitted HTTP support PR (#31) to virtuald's fork
- Nov 2025
