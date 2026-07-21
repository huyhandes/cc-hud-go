# Code Map - cc-hud-go

**Stats:** ~30 Go files · default-only mode · 4 Catppuccin themes

## Project Structure

```
cc-hud-go/
├── main.go                    # Entry point, CLI flags, stdin reading
├── main_test.go               # Main package tests
├── integration_test.go        # End-to-end integration tests
├── go.mod                     # Go module dependencies
├── go.sum                     # Dependency checksums
├── justfile                   # Build automation (Just)
├── CLAUDE.md                  # Claude Code project instructions
├── README.md                  # User-facing documentation
├── CHANGELOG.md               # Version history
├── LICENSE
│
├── state/                     # Session state tracking
│   ├── state.go              # State struct, derived field calculation
│   └── state_test.go         # State tests
│
├── parser/                    # Stdin JSON parsing (single concern)
│   ├── stdin.go              # StdinData type, ParseStdin()
│   └── stdin_test.go         # Stdin parser tests
│
├── segment/                   # Display segments (modular components)
│   ├── segment.go            # Segment interface, All(), ByID() registry
│   ├── model.go              # Model name display
│   ├── caveman.go            # Caveman mode badge (self-gating)
│   ├── ponytail.go           # Ponytail mode badge (self-gating)
│   ├── git.go                # Git branch, status, file stats
│   ├── fivehour.go           # 5-hour rate limit tracking
│   ├── ratelimit.go          # 7-day API rate limit tracking
│   └── *_test.go             # Segment tests
│
├── output/                    # Output formatting
│   ├── renderer.go           # Multi-line & single-line layouts
│   └── renderer_test.go      # Renderer tests
│
├── format/                    # Shared formatting helpers (DRY)
│   ├── format.go             # Tokens(), Duration(), Cost()
│   └── format_test.go        # Formatter tests
│
├── style/                     # Lipgloss styling system
│   ├── style.go              # Colors, styles, Init(), ThresholdColor()
│   ├── gradient.go           # RenderGradientBar(), color interpolation
│   └── gradient_test.go      # Gradient rendering tests
│
├── theme/                     # Theme system
│   ├── theme.go              # Theme interface & loader
│   ├── catppuccin.go         # 4 Catppuccin themes (Macchiato/Mocha/Frappe/Latte)
│   └── theme_test.go         # Theme tests
│
├── version/                   # Version management
│   ├── version.go            # Git-based version detection
│   └── version_test.go
│
├── internal/
│   ├── git/                  # Git command integration
│   │   ├── git.go            # Branch, status, diff stats
│   │   └── git_test.go
│   └── oauth/                # OAuth authentication helpers
│       ├── oauth.go          # OAuth flow implementation
│       └── oauth_test.go
│
├── cmd/                       # Developer utility binaries
│   ├── test-gradient/        # Gradient bar visual tester
│   │   └── main.go
│   └── test-oauth/           # OAuth flow tester
│       └── main.go
│
├── docs/                      # Documentation
│   ├── adr/                  # Architectural Decision Records
│   ├── CODEMAP.md            # This file
│   ├── MANUAL_TEST.md        # Manual testing guide
│   ├── TEST_RESULTS.md       # Test results
│   ├── BUILD_GUIDE.md        # Build instructions
│   ├── PROJECT_STATUS.md     # Current project status
│   ├── README_FIXES.md       # README correction notes
│   ├── BUG_FIXES.md          # Bug fix documentation
│   ├── CI_FIXES.md           # CI fix documentation
│   ├── RELEASE_NOTES_v0.2.0.md
│   └── plans/                # Design and implementation plans
│
├── testdata/                  # Test fixtures
│
├── assets/                    # Screenshots and preview images
│   └── preview.jpeg
│
└── .github/                   # CI/CD
    ├── workflows/
    │   ├── ci.yml             # Tests, lint, build (skips docs)
    │   └── release.yml        # Multi-platform release builds
    └── WORKFLOWS.md
```

## Package Descriptions

| Package | Purpose |
|---------|---------|
| `main` | Entry point: CLI flags, stdin reading, theme init, rendering |
| `state` | Centralized session state with automatic derived field calculation |
| `parser` | Stdin JSON parsing from Claude Code |
| `segment` | Modular display components implementing the Segment interface |
| `output` | Output renderer producing multi-line and single-line statusline layouts |
| `format` | Shared formatting helpers: Tokens(), Duration(), Cost() |
| `style` | Lipgloss styling: semantic colors, gradient bars |
| `theme` | Theme system with 4 Catppuccin variants (Macchiato/Mocha/Frappe/Latte) |
| `version` | Git-based version detection and build info |
| `internal/git` | Git integration via command execution: branch, status, diff stats |
| `internal/oauth` | OAuth authentication helpers for Claude Code API access (always fetched) |
| `cmd/test-gradient` | Developer utility for visual testing of gradient bars |
| `cmd/test-oauth` | Developer utility for testing the OAuth authentication flow |

## Dependency Graph

### Hub Packages (most imported)

```
state    ← session data flows through here
style    ← all rendering code uses style
theme    ← 4 importers
format   ← 3 importers
oauth    ← 2 importers
```

### Top-level import relationships

```
main
  ├── internal/git
  ├── internal/oauth
  ├── state
  ├── style
  ├── theme
  └── version

output/renderer
  ├── format
  ├── segment  (all)
  ├── state
  └── style

parser/stdin      → state

segment/*         → state, style

style/style       → theme
```

## Data Flow

```
┌─────────────────────────────────────────────────────────────┐
│                     Claude Code                              │
│  (sends JSON via stdin)                                      │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        ▼
              ┌─────────────────┐
              │   main.go       │
              │  - Read stdin   │
              │  - Init theme   │
              │  - Fetch OAuth  │
              └────────┬────────┘
                       │
         ┌─────────────┼─────────────┐
         ▼             ▼             ▼
   ┌─────────┐  ┌──────────┐  ┌─────────┐
   │ Parser  │  │  OAuth   │  │  Theme  │
   │ - Stdin │  │ - Rates  │  │ - Colors│
   └────┬────┘  └─────┬────┘  └────┬────┘
        │             │             │
        ▼             ▼             ▼
   ┌────────────────────────────────────┐
   │           State                    │
   │  - Session data                    │
   │  - Rate limits                     │
   │  - Derived fields                  │
   └───────────────┬────────────────────┘
                   │
        ┌──────────┼──────────┐
        ▼          ▼          ▼
   ┌────────┐ ┌────────┐ ┌────────┐
   │Segment │ │Segment │ │Segment │
   │Model   │ │  Git   │ │  ...   │
   └───┬────┘ └───┬────┘ └───┬────┘
       │          │          │
       └──────────┼──────────┘
                  ▼
         ┌─────────────────┐
         │  Output/Renderer│
         │  - Multi-line   │
         │  - Single-line  │
         └────────┬────────┘
                  │
                  ▼
         ┌─────────────────┐
         │  JSON output    │
         │  (to stdout)    │
         └─────────────────┘
                  │
                  ▼
         ┌─────────────────┐
         │  Claude Code    │
         │  (statusline)   │
         └─────────────────┘
```

## Key Files by Function

### Entry Point
- `main.go`
  - CLI flag parsing (--help, --version, --theme)
  - Stdin reading and parsing
  - Theme initialization
  - OAuth rate-limit fetch (always, no stdin fallback)
  - Output rendering

### State Management
- `state/state.go`
  - Centralized session state
  - Automatic derived field calculation (percentages, totals)
  - Git, rate-limit tracking

### Parsing
- `parser/stdin.go` - StdinData struct, ParseStdin()

### Display Segments

Each segment self-gates by returning `""` when it has nothing to render.

Available segments (via `All()` or `ByID()` registry):

| File | ID | Display |
|------|----|---------|
| `model.go` | `model` | Current Claude model name |
| `caveman.go` | `caveman` | Caveman mode badge (when active) |
| `ponytail.go` | `ponytail` | Ponytail mode badge (when active) |
| `git.go` | `git` | Branch, dirty files, ahead/behind, file stats |
| `fivehour.go` | `fivehour` | 5-hour API rate limit tracking |
| `ratelimit.go` | `ratelimit` | 7-day API rate limit tracking |

### Formatting
- `format/format.go` - Shared helpers: Tokens(), Duration(), Cost()

### Output Rendering
- `output/renderer.go`
  - **renderMultiLine()** - Custom multi-line layout (uses ByID() map)
  - **renderSingleLine()** - Compact horizontal layout
  - **renderFileChanges()** - +/- line changes

### Styling
- `style/style.go` - Colors, Init(), ThresholdColor(), GetRenderer()
- `style/gradient.go` - RenderGradientBar(), color interpolation
- `style/table.go` - RenderTable() box-drawing tables

### Themes
- `theme/catppuccin.go`
  - **4 themes:** Macchiato (default), Mocha, Frappe, Latte
  - **LoadThemeFromConfig()** - Theme selection + color overrides
  - **ThemeWrapper** - Custom color override support

## External Dependencies

From the Charm ecosystem (via go.mod):

| Dependency | Purpose |
|-----------|---------|
| `charmbracelet/lipgloss` | Terminal styling and layout |
| `charmbracelet/colorprofile` | Color profile detection |
| `charmbracelet/x/ansi` | ANSI escape sequence handling |
| `charmbracelet/x/cellbuf` | Terminal cell buffer |
| `charmbracelet/x/term` | Terminal utilities |
| `muesli/termenv` | Terminal environment detection |
| `lucasb-eyer/go-colorful` | Color manipulation and interpolation |
| `mattn/go-isatty` | TTY detection |
| `mattn/go-runewidth` | Unicode rune width calculation |
| `rivo/uniseg` | Unicode segmentation |

## Output Format

### Multi-line Layout Example
```
Line 1: model name
Line 2: git branch + dirty files + file change stats
Line 3: rate-limit segments (when OAuth data is available)
```

## Important Implementation Details

### 1. Gradient Bar Colors
- Green (0-69%): Healthy usage
- Yellow (70-89%): Warning
- Red (90-100%): Danger

### 2. CI Behavior
CI **runs** on:
- Go source changes (`*.go`)
- Module changes (`go.mod`, `go.sum`)
- Workflow changes (`.github/workflows/*`)

CI **skips** on:
- Documentation (`*.md`, `docs/**`)
- Assets (`assets/**`)

## Test Coverage

```
Package          Coverage
─────────────────────────
main             62%
state            100%
parser           85%
segment          92%
output           78%
style            85%
theme            100%
version          88%
internal/git     75%
internal/oauth   (new)
─────────────────────────
Overall          ~85%
```

## Version History

- **v0.1.0** - Initial release (basic segments)
- **v0.2.0** - Visual enhancements (themes, gradients, bug fixes)
- **default-only** - Config file removed; `--theme` CLI flag; transcript parsing removed; OAuth always fetched

## Contributing

When adding features:
1. Follow TDD approach (test first)
2. Update both code and documentation
3. Run `just check` before committing
4. Ensure CI passes on all platforms

---

**Last Updated:** 2026-02-17
**Version:** v0.2.0
