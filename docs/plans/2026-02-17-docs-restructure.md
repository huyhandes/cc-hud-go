# Documentation Restructure Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Shorten README.md to a quick-start guide and extract detailed content into focused reference documents under `docs/`.

**Architecture:** Split the monolithic README into focused references: `docs/BUILD_GUIDE.md`, `docs/CONFIG.md`, updated `docs/COLOR_SCHEME.md`, and updated `docs/CODEMAP.md`. Update `CLAUDE.md` workflow and `CHANGELOG.md`.

**Tech Stack:** Markdown only — no code changes.

---

### Codebase Context

Key facts from codemap for CODEMAP.md update:
- **52 Go files**, 192 functions, 73 deps
- **Hub packages:** `state` (27 imports), `config` (22), `style` (12), `theme` (4), `format` (3)
- Dependency chain: `main → config, git, oauth, state, style, theme, version`
- `output/renderer` → `config, format, segment, state, style`
- All segments → `config, state, style`
- Notable packages not in current docs: `internal/oauth`, `cmd/test-gradient`, `cmd/test-oauth`

---

### Task 1: Create `docs/BUILD_GUIDE.md`

**Files:**
- Create: `docs/BUILD_GUIDE.md`

**Step 1: Create the file**

```markdown
# Build Guide

## Prerequisites

- Go 1.24+
- [Just](https://github.com/casey/just) (task runner)
- Git

## Quick Install (Pre-built Binaries)

Download from the [latest release](https://github.com/huyhandes/cc-hud-go/releases/latest):

\`\`\`bash
# macOS (Apple Silicon)
curl -L https://github.com/huyhandes/cc-hud-go/releases/latest/download/cc-hud-go-darwin-arm64.tar.gz | tar xz
sudo mv cc-hud-go-darwin-arm64 /usr/local/bin/cc-hud-go

# macOS (Intel)
curl -L https://github.com/huyhandes/cc-hud-go/releases/latest/download/cc-hud-go-darwin-amd64.tar.gz | tar xz
sudo mv cc-hud-go-darwin-amd64 /usr/local/bin/cc-hud-go

# Linux (amd64)
curl -L https://github.com/huyhandes/cc-hud-go/releases/latest/download/cc-hud-go-linux-amd64.tar.gz | tar xz
sudo mv cc-hud-go-linux-amd64 /usr/local/bin/cc-hud-go
\`\`\`

Available: Linux (`amd64`, `arm64`), macOS (`amd64`, `arm64`), Windows (`amd64`, `arm64`)

## Build from Source

\`\`\`bash
git clone git@github.com:huyhandes/cc-hud-go.git
cd cc-hud-go
just build          # builds with version info from git tags
just install        # installs to ~/.local/bin
\`\`\`

Or with `go install`:
\`\`\`bash
go install github.com/huyhandes/cc-hud-go@latest
\`\`\`

## Just Commands

| Command | Description |
|---------|-------------|
| `just build` | Build with version from git tags |
| `just install` | Build and install to `~/.local/bin` |
| `just test` | Run all tests |
| `just test-coverage` | Run tests with coverage report |
| `just check` | Format, vet, and test |
| `just fmt` | Format code (`go fmt ./...`) |
| `just vet` | Run `go vet ./...` |
| `just lint` | Run `golangci-lint run` |
| `just clean` | Remove build artifacts |
| `just build-all` | Build for all platforms |

## Running Tests

\`\`\`bash
just test                                    # all tests
just test-coverage                           # with coverage report
just test -- -run TestModelSegment ./segment # specific test
just test -- -v ./segment                    # verbose, one package
\`\`\`

## Creating a Release

\`\`\`bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
\`\`\`

GitHub Actions will automatically:
1. Build binaries for all platforms (Linux/macOS/Windows × amd64/arm64)
2. Create compressed archives (`.tar.gz` / `.zip`)
3. Generate SHA256 checksums
4. Publish the GitHub release with auto-generated notes

You can also trigger a release manually from the Actions tab.

## Version Information

- Release builds: tagged version (e.g. `v0.2.0`)
- Dev builds: `git describe` output (e.g. `v0.2.0-3-gabc1234-dirty`)
- Without git: falls back to `dev`
```

**Step 2: Commit**

```bash
git add docs/BUILD_GUIDE.md
git commit -m "docs: add BUILD_GUIDE.md"
```

---

### Task 2: Create `docs/CONFIG.md`

**Files:**
- Create: `docs/CONFIG.md`

**Step 1: Create the file**

Extract the full configuration reference from README.md into this file. Content:

```markdown
# Configuration Reference

Config file location: `~/.claude/cc-hud-go/config.json`

## Full Example

\`\`\`json
{
  "theme": "macchiato",
  "colors": {},
  "preset": "full",
  "lineLayout": "expanded",
  "pathLevels": 2,
  "contextValue": "percentage",
  "sevenDayThreshold": 80,
  "display": {
    "model": true,
    "path": true,
    "context": true,
    "git": true,
    "tools": true,
    "agents": true,
    "tasks": true,
    "rateLimits": true,
    "duration": true,
    "speed": true
  },
  "git": {
    "showBranch": true,
    "showDirty": true,
    "showAheadBehind": true,
    "showFileStats": true
  },
  "tools": {
    "groupByCategory": true,
    "showTopN": 5,
    "showSkills": true,
    "showMCP": true
  },
  "tables": {
    "toolsTableThreshold": 999,
    "tasksTableThreshold": 999,
    "contextTableThreshold": 999
  }
}
\`\`\`

## Top-Level Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `theme` | string | `"macchiato"` | Color theme: `macchiato`, `mocha`, `frappe`, `latte` |
| `colors` | object | `{}` | Custom semantic color overrides (hex) |
| `preset` | string | `"full"` | Preset: `full`, `essential`, `minimal` |
| `lineLayout` | string | `"expanded"` | Layout: `expanded` or `compact` |
| `pathLevels` | int | `2` | Directory levels to show in path (1–3) |
| `contextValue` | string | `"percentage"` | Context display format |
| `sevenDayThreshold` | int | `80` | Warning threshold for 7-day rate limit (0–100) |

## Presets

| Preset | Description |
|--------|-------------|
| `full` | All segments enabled |
| `essential` | Core metrics only (model, context, git, cost) |
| `minimal` | Minimal info (model, context) |

## Display Options

Boolean flags to enable/disable individual segments:

| Key | Default | Description |
|-----|---------|-------------|
| `model` | `true` | Model name and plan type |
| `path` | `true` | Current working directory |
| `context` | `true` | Token usage with color thresholds |
| `git` | `true` | Git branch, status, file stats |
| `tools` | `true` | Tool usage by category |
| `agents` | `true` | Active agent name and task |
| `tasks` | `true` | Task completion progress |
| `rateLimits` | `true` | 7-day API usage |
| `duration` | `true` | Session duration |
| `speed` | `true` | Token processing speed |

## Git Options

| Key | Default | Description |
|-----|---------|-------------|
| `showBranch` | `true` | Current branch name |
| `showDirty` | `true` | Count of uncommitted files |
| `showAheadBehind` | `true` | Commits ahead/behind remote |
| `showFileStats` | `true` | Added/modified/deleted counts |

## Tools Options

| Key | Default | Description |
|-----|---------|-------------|
| `groupByCategory` | `true` | Group by App/MCP/Skills/Custom |
| `showTopN` | `5` | Number of top tools to show (0 = all) |
| `showSkills` | `true` | Include skill usage |
| `showMCP` | `true` | Include MCP tool usage |

## Table Options

Smart threshold: when item count exceeds the threshold, switches from styled lipgloss boxes (╭╮) to plain table view (┌┬┐).

| Key | Default | Description |
|-----|---------|-------------|
| `toolsTableThreshold` | `999` | Tool count for table view |
| `tasksTableThreshold` | `999` | Task count for table view |
| `contextTableThreshold` | `999` | Context size for table view |

Default of 999 always uses the styled lipgloss box. Lower to 3–5 for earlier table switching on large datasets.

## Custom Colors

Override any semantic color while keeping the base theme:

\`\`\`json
{
  "theme": "macchiato",
  "colors": {
    "success": "#00ff00",
    "warning": "#ffaa00",
    "danger": "#ff0000",
    "primary": "#00aaff"
  }
}
\`\`\`

Available semantic color keys:

| Key | Role |
|-----|------|
| `success` | Positive states, low usage (green) |
| `warning` | Caution, medium usage (yellow/orange) |
| `danger` | Critical states, high usage (red) |
| `input` | Input tokens (blue) |
| `output` | Output tokens (emerald) |
| `cacheRead` | Cache read tokens (purple) |
| `cacheWrite` | Cache write tokens (pink) |
| `primary` | Main brand color (purple) |
| `highlight` | Highlights and accents (cyan) |
| `accent` | Secondary accents (orange) |
| `muted` | Borders, subtle elements (gray) |
| `bright` | Bright text (white/cream) |
| `info` | Informational elements (teal) |

## Example Configs

Pre-configured examples in `examples/`:

\`\`\`bash
cp examples/config-macchiato.json ~/.claude/cc-hud-go/config.json
cp examples/config-mocha.json ~/.claude/cc-hud-go/config.json
cp examples/config-frappe.json ~/.claude/cc-hud-go/config.json
cp examples/config-latte.json ~/.claude/cc-hud-go/config.json
cp examples/config-custom-colors.json ~/.claude/cc-hud-go/config.json
\`\`\`

See `examples/README.md` for the full theme guide.
```

**Step 2: Commit**

```bash
git add docs/CONFIG.md
git commit -m "docs: add CONFIG.md reference"
```

---

### Task 3: Update `docs/COLOR_SCHEME.md`

**Files:**
- Modify: `docs/COLOR_SCHEME.md`

**Step 1: Read the existing file**

Read `docs/COLOR_SCHEME.md` to understand what's already there, then append the Themes section from README.

**Step 2: Append themes section**

Add a "## Themes" section at the end of the existing `docs/COLOR_SCHEME.md` with the 4 Catppuccin variants (Macchiato, Mocha, Frappe, Latte), the visual features descriptions (gradient bars, multi-line layout), and the smart adaptive layouts explanation from the README Themes section.

**Step 3: Commit**

```bash
git add docs/COLOR_SCHEME.md
git commit -m "docs: merge themes section into COLOR_SCHEME.md"
```

---

### Task 4: Update `docs/CODEMAP.md`

**Files:**
- Modify: `docs/CODEMAP.md`

**Step 1: Read existing file**

Read `docs/CODEMAP.md` to understand what to preserve and what to update.

**Step 2: Overwrite with updated content**

Replace with a fresh codemap based on the `codemap .` and `codemap -deps` output gathered during planning.

Key additions vs current file:
- Add `internal/oauth` package
- Add `cmd/test-gradient` and `cmd/test-oauth` utilities
- Add dependency flow section from `codemap -deps` (hub packages: state←27, config←22, style←12)
- Add external dependencies list (lipgloss, termenv, go-colorful, etc.)

**Step 3: Commit**

```bash
git add docs/CODEMAP.md
git commit -m "docs: regenerate CODEMAP.md with codemap output"
```

---

### Task 5: Update `CLAUDE.md`

**Files:**
- Modify: `CLAUDE.md`

**Step 1: Add Workflow section**

After the "Linting" section, add:

```markdown
## Workflow

After finishing any task, always run:
```bash
just fmt
just lint
```
```

**Step 2: Remove bare `go` commands**

- In "Manual Build" section: keep it (it's for users who don't have `just`) — but scope it clearly as "without `just`"
- In "Tests" section: remove the `go test ./...` / `go test -cover ./...` examples; keep only `just test` and `just test-coverage`, plus the specific-test example using `go test -run` (no equivalent in just)
- In "Linting" section: remove `go fmt ./...` and `go vet ./...` manual commands at bottom; keep only `just` commands
- In "Dependencies" section: `go get` / `go mod tidy` are fine to keep — no `just` equivalent

**Step 3: Commit**

```bash
git add CLAUDE.md
git commit -m "docs: add lint/fmt workflow to CLAUDE.md, prefer just over go commands"
```

---

### Task 6: Update `CHANGELOG.md`

**Files:**
- Modify: `CHANGELOG.md`

**Step 1: Add unreleased entry and fix v0.1.0 date**

- Fix `## [0.1.0] - 2025-01-XX` → `## [0.1.0] - 2025-01-01`
- Add at top:

```markdown
## [Unreleased]

### Changed
- **Documentation restructure**: README shortened to quick-start; detailed content moved to `docs/BUILD_GUIDE.md`, `docs/CONFIG.md`, `docs/COLOR_SCHEME.md`, `docs/CODEMAP.md`
- **CLAUDE.md**: Added mandatory lint/fmt workflow; prefer `just` over bare `go` commands
```

**Step 2: Commit**

```bash
git add CHANGELOG.md
git commit -m "docs: update CHANGELOG with docs restructure entry"
```

---

### Task 7: Shorten `README.md`

**Files:**
- Modify: `README.md`

**Step 1: Rewrite to ~150 lines**

Replace the full README content with a shortened version:

```markdown
# cc-hud-go

A Go-based statusline tool for [Claude Code](https://code.claude.com) that displays rich, real-time information about your current session.

[badges...]

![Preview](assets/preview.jpeg)

## Features

- **Model & Context** — Model name, plan type, token usage with color-coded thresholds
- **Rate Limits** — 5-hour and 7-day API usage tracking
- **Cost & Duration** — Session cost (USD), duration, token speed
- **Git** — Branch, dirty files, ahead/behind, file stats
- **Tools** — Categorized usage (App / MCP / Skills / Custom)
- **Tasks** — Completion progress (pending / in-progress / done)
- **Agents** — Active agent name and current task

## Quick Install

\`\`\`bash
# macOS (Apple Silicon)
curl -L https://github.com/huyhandes/cc-hud-go/releases/latest/download/cc-hud-go-darwin-arm64.tar.gz | tar xz
sudo mv cc-hud-go-darwin-arm64 /usr/local/bin/cc-hud-go
\`\`\`

See [Build Guide](docs/BUILD_GUIDE.md) for all platforms, building from source, and `go install`.

## Usage

Add to your Claude Code statusline config:

\`\`\`json
{
  "statusline": {
    "command": "cc-hud-go"
  }
}
\`\`\`

Test standalone:
\`\`\`bash
echo '{"model":"claude-sonnet-4.5","context":{"used":5000,"total":10000}}' | cc-hud-go
cc-hud-go --help
cc-hud-go --version
\`\`\`

## Configuration

Create `~/.claude/cc-hud-go/config.json`:

\`\`\`json
{
  "theme": "macchiato",
  "preset": "full"
}
\`\`\`

See [Configuration Reference](docs/CONFIG.md) for all options.

## Themes

4 [Catppuccin](https://github.com/catppuccin/catppuccin) themes: `macchiato` (default), `mocha`, `frappe`, `latte`.

\`\`\`json
{ "theme": "mocha" }
\`\`\`

See [Color Scheme Guide](docs/COLOR_SCHEME.md) for theme details and custom color overrides.

## Reference

| Document | Description |
|----------|-------------|
| [Build Guide](docs/BUILD_GUIDE.md) | Install, build, test, release |
| [Configuration Reference](docs/CONFIG.md) | All config options |
| [Color Scheme Guide](docs/COLOR_SCHEME.md) | Themes and color customization |
| [Code Map](docs/CODEMAP.md) | Architecture and package dependencies |

## Architecture

Composable segment system: each display element is an independent `Segment` implementing `Render(state, config) (string, error)`. State flows: stdin JSON → parser → state → segments → JSON output.

See [Code Map](docs/CODEMAP.md) for full package structure.

## Contributing

1. Fork → feature branch → `just check` → PR
2. New segments: `segment/<name>.go` + test + register in `segment/segment.go`
3. See [CLAUDE.md](CLAUDE.md) for dev workflow

## License

MIT — see [LICENSE](LICENSE)

## Links

- [Statusline API docs](https://code.claude.com/docs/en/statusline)
- [Issue tracker](https://github.com/huyhandes/cc-hud-go/issues)
- [Releases](https://github.com/huyhandes/cc-hud-go/releases)
```

**Step 2: Commit**

```bash
git add README.md
git commit -m "docs: shorten README, link to reference docs"
```

---

### Task 8: Final check

**Step 1: Verify all links in README resolve**

Read `README.md` and confirm each `docs/*.md` link corresponds to a file that now exists.

**Step 2: Run lint/fmt (no Go changes, but verify habit)**

```bash
just fmt
just lint
```

Expected: no output / 0 issues (no Go files were changed).

**Step 3: Final commit if any fixups**

```bash
git add -A
git commit -m "docs: fixups from final review"
```
