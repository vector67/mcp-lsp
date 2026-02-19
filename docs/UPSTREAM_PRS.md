# vector67/mcp-language-server — Open PR Inventory

## Date: 2026-02-17

Full inventory of all 20 open PRs on the upstream repository, with analysis for cherry-picking into a virtuald-based fork.

---

## Architecture / Major Features

### PR #97 — Multi-LSP session support
- **Author:** steffen-heil-secforge (Steffen Heil, secforge GmbH)
- **Opened:** 2025-11-04
- **Size:** +3,361 / -105 across 15 files (1 commit) — **LARGEST PR**
- **Reviews:** None
- **Description:** Complete multi-LSP infrastructure enabling simultaneous management of multiple language servers in a single session. Three modes: Single-MCP (original), Unbounded (dynamic start/stop via config), and Session (pre-configured). Adds tools: `lsp_start`, `lsp_stop`, `lsp_select`, `lsp_list`, `lsp_save`, `lsp_load`. Includes 50+ tests and full documentation.
- **Verdict:** DEFER — Architectural rewrite incompatible with virtuald's additions. Would need reimplementation from scratch.

### PR #116 — Headless mode (connect to existing IDE LSP)
- **Author:** rickbatka (Rick Batka)
- **Opened:** 2026-02-08
- **Size:** +853 / -582 across 15 files (7 commits)
- **Reviews:** 2 self-comments (inline notes). No external reviews.
- **Description:** Adds `--lsp-connect=hostname:port` to connect to an already-running LSP (e.g., the one in your IDE) via TCP. Avoids spawning parallel LSP, shares IDE configuration, works around "too many open files" bug by skipping file watchers.
- **Verdict:** APPLY (Phase 4) — High value. Core new function cherry-picks cleanly; plumbing needs manual integration.

### PR #99 — LSP-compliant glob pattern matching
- **Author:** steffen-heil-secforge
- **Opened:** 2025-11-05
- **Size:** +526 / -133 across 3 files (1 commit)
- **Reviews:** None
- **Description:** Replaces custom glob matching with gopls glob pattern matcher for full LSP 3.17 spec compliance. Supports `*`, `?`, `**`, `{}`, `[]`, `[!...]`. Adds thread-safe caching with safety valve. Fixes upstream issue #54 (stale diagnostics).
- **Verdict:** APPLY (Phase 3) — New `internal/glob/` package applies cleanly. watcher.go integration needs manual work.

---

## New Tools / Capabilities

### PR #100 — `incoming_calls` tool
- **Author:** netbrah (palanisd)
- **Opened:** 2025-11-07
- **Size:** +350 / -2 across 9 files (5 commits)
- **Reviews:** None
- **Description:** Adds incoming_calls tool for LSP call hierarchy support. Includes integration tests for Go (multiple callers, single caller, method callers, no callers) and snapshot files.
- **Verdict:** SKIP — virtuald already has callers+callees (more comprehensive).

### PR #107 — MCP tool annotations
- **Author:** bryankthompson (Bryan Thompson)
- **Opened:** 2025-12-27
- **Size:** +15 / -3 across 3 files (1 commit)
- **Reviews:** None
- **Description:** Adds MCP tool annotations (`ReadOnlyHint`, `DestructiveHint`, human-readable titles) to all 6 tools. Helps MCP clients display better UI, warn before destructive ops, auto-approve safe tools.
- **Verdict:** APPLY (Phase 2) — Small change. Must also annotate virtuald's 3 extra tools.

### PR #71 — Allow disabling individual tools
- **Author:** joel-u410 (Joel Nordell)
- **Opened:** 2025-08-18
- **Size:** +12 / -0 across 2 files (1 commit)
- **Reviews:** None
- **Description:** Adds `-disable-tools` flag. Motivated by aider user whose LLM preferred `edit_file` over aider's built-in editing.
- **Verdict:** APPLY (Phase 2) — Small, conceptually orthogonal.

---

## Bug Fixes

### PR #82 — Fix Windows file URIs
- **Author:** LaelLuo (Lael)
- **Opened:** 2025-09-16
- **Size:** +84 / -84 across 10 files (1 commit)
- **Reviews:** None
- **Description:** Replaces `fmt.Sprintf("file://%s", ...)` with `protocol.URIFromPath(...)` throughout codebase. Fixes completely broken Windows support (invalid `file://C:\path` URIs).
- **Verdict:** APPLY (Phase 3) — Important fix. Do as manual find-and-replace, NOT patch apply.

### PR #80 — Fix logging/setLevel via mcp-go v0.39.1
- **Author:** LaelLuo
- **Opened:** 2025-09-16
- **Size:** +48 / -54 across 3 files (1 commit)
- **Reviews:** None
- **Description:** Upgrades mcp-go to v0.39.1, properly implements `logging/setLevel` (was declared but non-functional). Adapts tool handlers to new `CallToolRequest` API. Fixes issue #79.
- **Verdict:** APPLY (Phase 4, last) — Do as fresh `go get` on merged codebase, not patch apply.

### PR #49 — Don't send multiple initialized notifications
- **Author:** micahscopes (Micah)
- **Opened:** 2025-06-19 (oldest open PR)
- **Size:** +2 / -8 across 1 file (1 commit)
- **Reviews:** 1 comment from virtuald noting he added this to his fork.
- **Description:** Fixes duplicate `initialized` LSP notifications causing panics in strict LSP servers (e.g., `async-lsp` in Rust).
- **Verdict:** SKIP — Already fixed in virtuald PR #4.

---

## Dependency Updates (Dependabot)

| PR | Title | Opened | Status |
|----|-------|--------|--------|
| #114 | Bump `golang.org/x/text` 0.25.0 -> 0.33.0 | 2026-02-01 | Open |
| #110 | Bump `mcp-go` 0.25.0 -> 0.43.2 | 2026-01-01 | Open |
| #105 | Bump `actions/checkout` v4 -> v6 | 2025-12-02 | Open |
| #93 | Bump `actions/setup-node` v4 -> v6 | 2025-11-01 | Open |
| #89 | Bump `actions/setup-go` v5 -> v6 | 2025-10-01 | Open |
| #87 | Bump `actions/setup-python` v5 -> v6 | 2025-10-01 | Open |
| #77 | Bump `stretchr/testify` 1.10.0 -> 1.11.1 | 2025-09-02 | Open |

**Verdict:** Handle as part of Phase 4 mcp-go upgrade. Bump all deps to latest at that point.

---

## SDK Upgrades

### PR #103 — Update mcp-go to v0.37.0
- **Author:** chloelee767 (Chloe Lee)
- **Opened:** 2025-11-25
- **Size:** +64 / -81 across 3 files (5 commits)
- **Reviews:** None
- **Description:** Migrates from removed APIs to new ones. Needed for MCP Inspector compatibility.
- **Verdict:** SKIP — Subsumed by #80 (v0.39.1).

---

## Documentation / README

### PR #62 — Claude Code setup instructions
- **Author:** Munawwar (Munawwar Firoz)
- **Opened:** 2025-07-07
- **Size:** +10 / -0 across 1 file (2 commits)
- **Reviews:** None
- **Description:** Adds Claude Code configuration instructions to README.
- **Verdict:** APPLY (Phase 2)

### PR #106 — Fix typo in README
- **Author:** jhstatewide (Joshua Harding)
- **Opened:** 2025-12-25
- **Size:** +1 / -1
- **Verdict:** APPLY (Phase 1)

### PR #113 — Fix typo in README
- **Author:** jhstatewide
- **Opened:** 2026-01-23
- **Size:** +1 / -1
- **Verdict:** APPLY (Phase 1)

---

## New Language Support

### PR #109 — Coq (.v) file support
- **Author:** larsr (Lars Rasmusson)
- **Opened:** 2025-12-30
- **Size:** +2 / -0 across 1 file (1 commit)
- **Reviews:** None
- **Description:** Adds `.v` extension mapping for `coq-lsp`.
- **Verdict:** APPLY (Phase 1) — Trivial, clean apply.

---

## Packaging / Distribution

### PR #70 — Nix flake
- **Author:** elsirion
- **Opened:** 2025-08-11
- **Size:** +130 / -0 across 3 files (1 commit)
- **Reviews:** None
- **Description:** Adds Nix flake with default package and dev shell. Enables `nix run github:vector67/mcp-language-server`.
- **Verdict:** APPLY (Phase 1) — All new files, zero conflict.

---

## Summary Table

| PR | Title | Phase | Action |
|----|-------|-------|--------|
| #49 | Duplicate initialized fix | — | SKIP (already in virtuald) |
| #62 | Claude Code README | 2 | APPLY |
| #70 | Nix flake | 1 | APPLY |
| #71 | Disable individual tools | 2 | APPLY |
| #77 | Bump testify | 4 | WITH DEP UPGRADE |
| #80 | mcp-go v0.39.1 | 4 | APPLY (last) |
| #82 | Windows file URIs | 3 | APPLY (manual) |
| #87 | Bump setup-python | 4 | WITH CI FIX |
| #89 | Bump setup-go | 4 | WITH CI FIX |
| #93 | Bump setup-node | 4 | WITH CI FIX |
| #97 | Multi-LSP sessions | — | DEFER |
| #99 | LSP-compliant glob | 3 | APPLY |
| #100 | incoming_calls | — | SKIP (superseded) |
| #103 | mcp-go v0.37.0 | — | SKIP (subsumed by #80) |
| #105 | Bump checkout | 4 | WITH CI FIX |
| #106 | Typo fix | 1 | APPLY |
| #107 | Tool annotations | 2 | APPLY |
| #109 | Coq support | 1 | APPLY |
| #110 | Bump mcp-go v0.43.2 | 4 | WITH DEP UPGRADE |
| #113 | Typo fix | 1 | APPLY |
| #114 | Bump golang.org/x/text | 4 | WITH DEP UPGRADE |
| #116 | Headless mode | 4 | APPLY |
