# Code Map - cc-hud-go

**Stats:** 52 Go files · 192 functions · 73 deps

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
├── config/                    # Configuration management
│   ├── config.go             # Config struct, presets (Full/Essential/Minimal)
│   └── config_test.go        # Configuration tests
│
├── state/                     # Session state tracking
│   ├── state.go              # State struct, derived field calculation
│   └── state_test.go         # State tests
│
├── parser/                    # Input parsing (split into 4 files)
│   ├── stdin.go              # StdinData type, ParseStdin()
│   ├── transcript.go         # TranscriptLine types, ParseTranscript*()
│   ├── task.go               # TaskItem, TaskTracker, task processing
│   ├── tool.go               # ToolCategory, CategorizeTool(), appTools map
│   ├── stdin_test.go         # Stdin parser tests
│   ├── transcript_test.go    # Transcript parser tests
│   └── tasks_test.go         # Task tracking tests
│
├── segment/                   # Display segments (modular components)
│   ├── segment.go            # Segment interface, All(), ByID() registry
│   ├── model.go              # Model name display
│   ├── context.go            # Token usage & gradient bar
│   ├── git.go                # Git branch, status, file stats
│   ├── cost.go               # Cost tracking & duration
│   ├── tools.go              # Tool usage categorization
│   ├── tasks.go              # Task progress dashboard
│   ├── agent.go              # Active agent display
│   ├── ratelimit.go          # API rate limit tracking (5h + 7d)
│   └── *_test.go             # Segment tests (incl. fivehour_test.go)
│
├── output/                    # Output formatting
│   ├── renderer.go           # Multi-line & single-line layouts
│   └── renderer_test.go      # Renderer tests
│
├── format/                    # Shared formatting helpers (DRY)
│   ├── format.go             # Tokens(), Duration(), Cost()
│   └── format_test.go        # Formatter tests
│
├── style/                     # Lipgloss styling system (split into 3 files)
│   ├── style.go              # Colors, styles, Init(), ThresholdColor()
│   ├── gradient.go           # RenderGradientBar(), color interpolation
│   ├── table.go              # RenderTable() box-drawing tables
│   ├── gradient_test.go      # Gradient rendering tests
│   └── table_test.go         # Table rendering tests
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
│   ├── CODEMAP.md            # This file
│   ├── COLOR_SCHEME.md       # Color palette reference
│   ├── CONFIG.md             # Configuration reference
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
├── examples/                  # Example configurations
│   ├── README.md
│   ├── config-macchiato.json
│   ├── config-mocha.json
│   ├── config-frappe.json
│   ├── config-latte.json
│   └── config-custom-colors.json
│
├── testdata/                  # Test fixtures
│   ├── config_valid.json
│   └── config_invalid.json
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
| `main` | Entry point: CLI flags, stdin reading, config loading, theme init, rendering |
| `config` | Config struct, three presets (Full/Essential/Minimal), validation and defaults |
| `state` | Centralized session state with automatic derived field calculation |
| `parser` | Dual input parsing: stdin JSON from Claude Code and JSONL transcript files |
| `segment` | Modular display components implementing the Segment interface |
| `output` | Output renderer producing multi-line and single-line statusline layouts |
| `format` | Shared formatting helpers: Tokens(), Duration(), Cost() |
| `style` | Lipgloss styling: semantic colors, gradient bars, box-drawing tables |
| `theme` | Theme system with 4 Catppuccin variants (Macchiato/Mocha/Frappe/Latte) |
| `version` | Git-based version detection and build info |
| `internal/git` | Git integration via command execution: branch, status, diff stats |
| `internal/oauth` | OAuth authentication helpers for Claude Code API access |
| `cmd/test-gradient` | Developer utility for visual testing of gradient bars |
| `cmd/test-oauth` | Developer utility for testing the OAuth authentication flow |

## Dependency Graph

### Hub Packages (most imported)

```
state    ← 27 importers  (session data flows through here)
config   ← 22 importers  (all segments + renderer read config)
style    ← 12 importers  (all rendering code uses style)
theme    ←  4 importers
format   ←  3 importers
oauth    ←  2 importers
```

### Top-level import relationships

```
main
  ├── config
  ├── internal/git
  ├── internal/oauth
  ├── state
  ├── style
  ├── theme
  └── version

output/renderer
  ├── config
  ├── format
  ├── segment  (all)
  ├── state
  └── style

parser/stdin      → state
parser/transcript → state
parser/task       → state

segment/*         → config, state, style  (all segments import these three)

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
              │  - Load config  │
              │  - Init theme   │
              └────────┬────────┘
                       │
         ┌─────────────┼─────────────┐
         ▼             ▼             ▼
   ┌─────────┐  ┌──────────┐  ┌─────────┐
   │ Parser  │  │  Config  │  │  Theme  │
   │ - Stdin │  │ - Preset │  │ - Colors│
   │ - JSONL │  │ - Options│  │ - Styles│
   └────┬────┘  └─────┬────┘  └────┬────┘
        │             │             │
        ▼             ▼             ▼
   ┌────────────────────────────────────┐
   │           State                    │
   │  - Session data                    │
   │  - Tool usage                      │
   │  - Task tracking                   │
   │  - Derived fields                  │
   └───────────────┬────────────────────┘
                   │
        ┌──────────┼──────────┐
        ▼          ▼          ▼
   ┌────────┐ ┌────────┐ ┌────────┐
   │Segment │ │Segment │ │Segment │
   │Model   │ │Context │ │  ...   │
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
  - CLI flag parsing (--help, --version)
  - Stdin reading and parsing
  - Configuration loading
  - Theme initialization
  - Output rendering

### Configuration System
- `config/config.go`
  - Config struct with display options
  - Three presets: Full, Essential, Minimal
  - Validation and defaults
  - Theme & color override support

### State Management
- `state/state.go`
  - Centralized session state
  - Automatic derived field calculation (percentages, totals)
  - Context, Git, Tools, Tasks, Cost tracking

### Parsing (split into 4 files)
- `parser/stdin.go` - StdinData struct, ParseStdin()
- `parser/transcript.go` - TranscriptLine types, ParseTranscript*()
- `parser/task.go` - TaskTracker, task tool processing
- `parser/tool.go` - ToolCategory, CategorizeTool()

### Display Segments

Each segment implements:
```go
type Segment interface {
    ID() string
    Render(s *state.State, cfg *config.Config) (string, error)
    Enabled(cfg *config.Config) bool
}
```

Available segments (via `All()` or `ByID()` registry):

| File | ID | Display |
|------|----|---------|
| `model.go` | `model` | Current Claude model name |
| `context.go` | `context` | Gradient bar + token counts |
| `git.go` | `git` | Branch, dirty files, ahead/behind, file stats |
| `cost.go` | `cost` | USD cost + session duration |
| `tools.go` | `tools` | Tool usage by category (App/MCP/Skills/Custom) |
| `tasks.go` | `tasks` | Task completion progress |
| `agent.go` | `agent` | Active agent name and current task |
| `ratelimit.go` | `ratelimit` | 5-hour and 7-day API rate limit tracking |

### Formatting
- `format/format.go` - Shared helpers: Tokens(), Duration(), Cost()

### Output Rendering
- `output/renderer.go`
  - **renderMultiLine()** - Custom 4-line layout (uses ByID() map)
  - **renderSingleLine()** - Compact horizontal layout
  - **renderContextBar()** - Gradient bar with percentage
  - **renderFileChanges()** - +/- line changes

### Styling (split into 3 files)
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
Line 1: model name + context gradient bar + percentage
Line 2: token breakdown (input/output/cache) + cost + duration
Line 3: git branch + dirty files + file change stats
Line 4: tool usage table (App / MCP / Skills / Custom counts)
```

## Important Implementation Details

### 1. Context Display (NO label emoji prefix)
The context segment renders as:
```
█▓▒░░░░░░░ 59% input:89k output:12k cache_r:45k cache_w:23k limit:200k
```

### 2. Gradient Bar Colors
- Green (0-69%): Healthy usage
- Yellow (70-89%): Warning
- Red (90-100%): Danger

### 3. Task ID Indexing
- Claude Code uses **1-based task IDs** ("1", "2", "3")
- Parser converts to **0-based array indices**
- Bug fixed in v0.2.0

### 4. Skill Tracking
- Skills tracked by full name (e.g., "superpowers:using-git-worktrees")
- Extracted from Skill tool's `input.skill` parameter
- Bug fixed in v0.2.0

### 5. CI Behavior
CI **runs** on:
- Go source changes (`*.go`)
- Module changes (`go.mod`, `go.sum`)
- Workflow changes (`.github/workflows/*`)

CI **skips** on:
- Documentation (`*.md`, `docs/**`)
- Examples (`examples/**`)
- Assets (`assets/**`)

## Test Coverage

```
Package          Coverage
─────────────────────────
main             62%
config           100%
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

- **v0.1.0** - Initial release (basic segments, presets)
- **v0.2.0** - Visual enhancements (themes, gradients, bug fixes)

## Contributing

When adding features:
1. Follow TDD approach (test first)
2. Update both code and documentation
3. Run `just check` before committing
4. Ensure CI passes on all platforms

---

**Last Updated:** 2026-02-17
**Version:** v0.2.0
