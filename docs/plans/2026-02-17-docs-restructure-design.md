# Documentation Restructure Design

**Date:** 2026-02-17
**Status:** Approved

## Goal

Shorten README.md to a quick-start guide and extract detailed content into focused reference documents under `docs/`. Update CLAUDE.md workflow and CHANGELOG.md.

## Document Structure

### Root (keep here)

| File | Action | Target length |
|------|--------|---------------|
| `README.md` | Shorten | ~150 lines |
| `CHANGELOG.md` | Add docs entry, fix v0.1.0 date | minor update |
| `CLAUDE.md` | Add lint/fmt workflow, remove bare `go` commands | minor update |

### New/Updated in `docs/`

| File | Action | Source content |
|------|--------|---------------|
| `docs/BUILD_GUIDE.md` | **Create** | Install from source, `just` commands, platform builds, creating releases, running tests |
| `docs/CONFIG.md` | **Create** | Full config reference: all options, presets, display/git/tools/table options, example configs |
| `docs/COLOR_SCHEME.md` | **Update** | Merge themes section from README into existing file |
| `docs/CODEMAP.md` | **Update** | Regenerate with `codemap .` and `codemap -deps` output |

## README.md Outline (shortened)

1. Title + badges
2. Preview image
3. Features (concise bullet list, no subsections)
4. Quick Install (pre-built binaries only)
5. Quick Usage (Claude Code integration + standalone)
6. Configuration (3-line minimal example + link to `docs/CONFIG.md`)
7. Themes (one-liner + link to `docs/COLOR_SCHEME.md`)
8. Reference links section → BUILD_GUIDE, CONFIG, COLOR_SCHEME, CODEMAP
9. Contributing (brief) + License + Links

**Removed from README:** full config tables, all theme details, from-source build steps, just commands list, architecture section, running tests section, creating releases section.

## CLAUDE.md Changes

- Add under a "Workflow" section:
  ```
  After finishing any task, always run:
  just fmt
  just lint
  ```
- Remove all bare `go` commands (go build, go test, go fmt, go vet) from examples — just commands only
- Keep only `just` commands throughout

## CHANGELOG.md Changes

- Fix `v0.1.0` date from `2025-01-XX` to `2025-01-01`
- Add `v0.3.0` or `[Unreleased]` entry for docs restructure

## Constraints

- `CLAUDE.md` stays at root
- All new reference docs go under `docs/`
- No content is deleted — only moved to appropriate reference files
- README links to all reference docs
