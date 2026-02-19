# MCP Language Server — Fork & Merge Plan

## Date: 2026-02-17

## Goal

Fork virtuald/mcp-language-server and apply all viable open PRs from vector67/mcp-language-server to create the most complete community fork.

## Starting Point

**Base:** virtuald/mcp-language-server (29 commits ahead of upstream, 0 behind)

**What we get for free:**
- All 6 upstream tools: definition, references, diagnostics, hover, rename_symbol, edit_file
- 3 additional tools: content, callers, callees
- File watcher glob fix
- Hover retry logic
- clangd improvements (smallest symbol matching, struct prefix stripping)
- Duplicate initialized notification fix
- mcp-go v0.33.0
- `-open` flag for initial workspace file opening
- Workspace symbol matching in definition tool

---

## PR Triage

### Skip (already done or superseded)

| PR | Title | Reason |
|----|-------|--------|
| #49 | Don't send duplicate initialized | Already fixed in virtuald PR #4 |
| #100 | incoming_calls tool | Superseded by virtuald's callers+callees (more comprehensive) |
| #103 | Update mcp-go to v0.37.0 | Subsumed by #80 (v0.39.1). virtuald already on v0.33.0 |

### Phase 1 — Trivial / Clean Apply (~15 min)

These touch files virtuald hasn't modified, or are entirely new files.

#### PR #70 — Nix flake
- **Author:** elsirion
- **Files:** `.gitignore`, `flake.lock`, `flake.nix` (all NEW)
- **Conflict risk:** NONE
- **Action:** Cherry-pick or apply patch directly

#### PR #109 — Coq (.v) file support
- **Author:** larsr
- **Files:** `internal/lsp/detect-language.go` only (+2 lines)
- **Conflict risk:** NONE — file is unmodified by virtuald
- **Action:** Cherry-pick or apply patch directly

#### PR #106 — Fix typo in README.md
- **Author:** jhstatewide
- **Files:** `README.md` (+1/-1)
- **Conflict risk:** TRIVIAL
- **Action:** Manual apply (1 line)

#### PR #113 — Fix typo in README.md
- **Author:** jhstatewide
- **Files:** `README.md` (+1/-1)
- **Conflict risk:** TRIVIAL
- **Action:** Manual apply (1 line)

### Phase 2 — Low Effort (~1 hr)

#### PR #62 — Claude Code README instructions
- **Author:** Munawwar
- **Files:** `README.md` (+10 lines)
- **Conflict risk:** LOW — both modified README but changes are additive
- **Action:** Manually add Claude Code setup section to virtuald's README

#### PR #107 — MCP tool annotations
- **Author:** bryankthompson
- **Files:** `go.mod`, `go.sum`, `tools.go`
- **Dependency:** Requires mcp-go v0.28.0+ (virtuald is on v0.33.0, already satisfied)
- **Conflict risk:** MODERATE on tools.go — virtuald has additional tools (callers, callees, content)
- **Action:** Manually add `WithTitleAnnotation()`, `WithReadOnlyHintAnnotation()`, `WithDestructiveHintAnnotation()` to each tool definition in virtuald's tools.go. Must also annotate the 3 extra tools (callers, callees, content) that the original PR doesn't cover.
- **Details:** Adds human-readable titles and hints (read-only vs destructive) to all tool definitions. `edit_file` and `rename_symbol` get `DestructiveHint: true`. All others get `ReadOnlyHint: true`.

#### PR #71 — Allow disabling individual tools
- **Author:** joel-u410
- **Files:** `main.go`, `tools.go` (+12 lines)
- **Conflict risk:** MODERATE — both files modified by virtuald
- **Action:** Manually apply concept:
  1. Add `-disable-tools` flag to virtuald's config struct in `main.go`
  2. Add parsing logic
  3. Add `DeleteTools()` call at end of `registerTools()` in `tools.go`
- **Motivation:** An `aider` user found Claude kept preferring `edit_file` over aider's built-in editing tool. Being able to disable specific tools prevents this.

### Phase 3 — Moderate Effort (~2 hrs)

#### PR #99 — LSP-compliant glob pattern matching
- **Author:** steffen-heil-secforge
- **Files:** `internal/watcher/watcher.go` + NEW `internal/glob/glob.go`, `internal/glob/glob_test.go`
- **Conflict risk:** MODERATE
- **What it does:** Replaces the broken `filepath.Match` (which doesn't support LSP 3.17 glob patterns like `**`, `{}`, `[!...]`) with a proper glob parser using gopls' implementation. Adds thread-safe caching. This fixes upstream issue #54 (stale diagnostics).
- **Action:**
  1. Copy the new `internal/glob/` package (applies cleanly — new directory)
  2. Manually integrate into virtuald's `watcher.go` — replace `filepath.Match` calls with the new glob matcher
- **Note:** virtuald's PR #12 partially addressed this by "removing duplicated glob matching logic" but still uses `filepath.Match` which is not LSP-compliant.

#### PR #82 — Fix Windows file URIs
- **Author:** LaelLuo
- **Files:** `internal/lsp/client.go` + 9 tool files (+84/-84 lines)
- **Conflict risk:** HIGH (touches many files virtuald also modified)
- **What it does:** Replaces all `fmt.Sprintf("file://%s", path)` with `protocol.URIFromPath(path)` throughout the codebase. On Windows, the former produces invalid URIs like `file://C:\path` instead of `file:///C:/path`.
- **Action:** Do NOT apply the patch. Instead, manually find-and-replace across virtuald's codebase:
  - Search for `fmt.Sprintf("file://%s"`
  - Replace with `protocol.URIFromPath(`
  - Check `client.go` (7 occurrences), `hover.go`, `references.go`, `lsp-utilities.go`, `utilities.go`, and all tool files
  - Verify `protocol.URIFromPath` is available in the imported LSP library version

### Phase 4 — Significant Effort (~2-4 hrs)

#### PR #116 — Headless mode (connect to existing IDE LSP)
- **Author:** rickbatka
- **Opened:** 2026-02-08 (very recent, may still be evolving)
- **Files:** `main.go`, `internal/lsp/client.go`, `internal/watcher/watcher.go` + tests (+853/-582 lines)
- **Conflict risk:** HIGH — every modified file is also modified by virtuald
- **What it does:** Adds `--lsp-connect=hostname:port` flag to connect to an already-running LSP server over TCP instead of spawning a new one. This:
  - Avoids duplicate LSP server instances (one in IDE, one in MCP)
  - Shares IDE configuration and plugins
  - Bypasses file watcher entirely (the IDE handles it)
  - Works around "too many open files" bug on large repos
- **Action:**
  1. `NewClientHeadless()` in client.go is a new function — cherry-picks cleanly
  2. `main.go` changes (new flag, connection logic) need manual integration with virtuald's config struct
  3. `watcher.go` changes (skip watcher in headless mode) need manual integration with virtuald's refactored watcher
  4. Test files should apply cleanly if placed correctly
- **Considerations:** This PR is from Feb 2026 and may still be in active development. Check for updates before applying.

#### PR #80 — mcp-go v0.39.1 + logging/setLevel fix
- **Author:** LaelLuo
- **Files:** `go.mod`, `go.sum`, `tools.go`
- **Conflict risk:** HIGH
- **What it does:** Upgrades mcp-go from v0.25.0 to v0.39.1, fixing the `logging/setLevel` capability that was declared but not implemented. Also migrates API calls from removed methods to new ones.
- **Action:** Do this LAST. Do NOT apply the patch (it's based on upstream's v0.25.0, not virtuald's v0.33.0). Instead:
  1. Run `go get github.com/mark3labs/mcp-go@v0.39.1` on the merged codebase
  2. Run `go mod tidy`
  3. Fix any compilation errors from API changes between v0.33.0 and v0.39.1
  4. The logging/setLevel fix should come for free with the upgrade
- **Note:** This subsumes PR #103 (v0.37.0). The Dependabot PR #110 targets v0.43.2 which is even newer — could go straight there instead, but may have more breaking changes.

### Defer Indefinitely

#### PR #97 — Multi-LSP session support (+3,361 lines)
- **Author:** steffen-heil-secforge
- **Why defer:** This is a near-complete architectural rewrite. It replaces `main.go` and `tools.go` with a multi-session management system (new tools: `lsp_start`, `lsp_stop`, `lsp_select`, `lsp_list`, `lsp_save`, `lsp_load`). Fundamentally incompatible with virtuald's additions.
- **If multi-LSP is needed:** Run multiple MCP server instances (one per language per workspace). This is how most users configure it today. Alternatively, re-implement multi-LSP from scratch on top of the merged codebase later.

---

## Inter-PR Conflicts

| Conflict | Resolution |
|----------|------------|
| #80 vs #103 vs #107 (all touch mcp-go) | Apply #107 first (annotations only, dep already met), then #80 last (subsumes #103) |
| #97 vs #116 (architectural) | Skip #97, apply #116 |
| #97 vs #71 (both modify tool registration) | Skip #97, apply #71 |
| #82 vs #116 (both modify client.go) | Apply #82 first, then #116 |
| #82 vs #49 (both modify client.go initialized area) | #49 already in virtuald, just apply #82 |
| #99 vs #116 (both modify watcher.go) | Apply #99 first, then #116 |

---

## Recommended Application Order

```
1. #70  — Nix flake (new files, zero risk)
2. #109 — Coq support (unmodified file)
3. #106 — Typo fix
4. #113 — Typo fix
5. #62  — Claude Code README
6. #107 — Tool annotations
7. #71  — Disable individual tools
8. #99  — LSP-compliant glob matching
9. #82  — Windows file URIs
10. #116 — Headless mode
11. #80  — mcp-go v0.39.1 upgrade (do last, touches everything)
```

---

## Estimated Effort

| Phase | PRs | Time |
|-------|-----|------|
| Phase 1 (trivial) | #70, #109, #106, #113 | ~15 min |
| Phase 2 (low effort) | #62, #107, #71 | ~1 hr |
| Phase 3 (moderate) | #99, #82 | ~2 hrs |
| Phase 4 (significant) | #116, #80 | ~2-4 hrs |
| **Total** | **11 PRs** | **~4-8 hrs** |

---

## Post-Merge TODO

- [ ] Fix CI (upstream issue #78)
- [ ] Verify rename_symbol works cross-file (upstream issue #104)
- [ ] Test with gopls, rust-analyzer, pyright, typescript-language-server, clangd
- [ ] Update README with complete setup instructions for Claude Code
- [ ] Consider adding `textDocument/codeAction` support (the big missing feature for extract/move refactoring)
- [ ] Consider connecting to Steffen about the multi-LSP work once the base is stable
- [ ] Set up Dependabot
- [ ] Create proper releases with Go binaries
