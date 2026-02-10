# Code Map - cc-hud-go

## Project Structure

```
cc-hud-go/
├── main.go                    # Entry point, CLI flags, stdin reading
├── go.mod                     # Go module dependencies
├── Makefile                   # Build automation
│
├── config/                    # Configuration management
│   ├── config.go             # Config struct, presets (Full/Essential/Minimal)
│   └── config_test.go        # Configuration tests
│
├── state/                     # Session state tracking
│   ├── state.go              # State struct, derived field calculation
│   └── state_test.go         # State tests
│
├── parser/                    # Input parsing
│   ├── parser.go             # Stdin JSON & transcript JSONL parsing
│   ├── stdin_test.go         # Stdin parser tests
│   ├── transcript_test.go    # Transcript parser tests
│   └── tasks_test.go         # Task tracking tests
│
├── segment/                   # Display segments (modular components)
│   ├── segment.go            # Segment interface & registry
│   ├── model.go              # 🤖 Model name & plan type
│   ├── context.go            # Token usage & gradient bar (NO 🧠 prefix)
│   ├── git.go                # 🌿 Git branch, status, file stats
│   ├── cost.go               # 💰 Cost tracking, ⏱ duration, 📝 file changes
│   ├── tools.go              # 📦 App, 🔌 MCP, ⚡ Skills, 🎨 Custom tools
│   ├── tasks.go              # Task progress dashboard
│   ├── agent.go              # 🤖 Active agent display
│   ├── ratelimit.go          # API rate limit tracking
│   └── *_test.go             # Segment tests
│
├── output/                    # JSON output formatting
│   ├── renderer.go           # Multi-line & single-line layouts
│   └── renderer_test.go      # Renderer tests
│
├── style/                     # Lipgloss styling system
│   ├── style.go              # Theme integration, gradient bars, tables
│   ├── gradient_test.go      # Gradient rendering tests
│   └── table_test.go         # Table rendering tests
│
├── theme/                     # Theme system
│   ├── theme.go              # Theme interface & loader
│   ├── catppuccin.go         # 4 Catppuccin themes (Macchiato/Mocha/Frappe/Latte)
│   └── theme_test.go         # Theme tests
│
├── internal/
│   ├── git/                  # Git command integration
│   │   ├── git.go           # Branch, status, diff stats
│   │   └── git_test.go
│   └── watcher/              # File watching utilities
│       └── watcher.go
│
├── version/                   # Version management
│   ├── version.go            # Git-based version detection
│   └── version_test.go
│
├── docs/                      # Documentation
│   ├── RELEASE_NOTES_v0.2.0.md
│   ├── BUG_FIXES.md
│   ├── CI_FIXES.md
│   └── CODEMAP.md            # This file
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
└── .github/workflows/         # CI/CD
    ├── ci.yml                # Tests, lint, build (skips docs)
    └── release.yml           # Multi-platform release builds
```

## Key Files by Function

### Entry Point
- `main.go` (208 lines)
  - CLI flag parsing (--help, --version)
  - Stdin reading and parsing
  - Configuration loading
  - Theme initialization
  - Output rendering

### Configuration System
- `config/config.go` (165 lines)
  - Config struct with display options
  - Three presets: Full, Essential, Minimal
  - Validation and defaults
  - Theme & color override support

### State Management
- `state/state.go` (120 lines)
  - Centralized session state
  - Automatic derived field calculation (percentages, totals)
  - Context, Git, Tools, Tasks, Cost tracking

### Parsing
- `parser/parser.go` (515 lines)
  - **StdinData struct** - Claude Code API format
  - **ParseStdin()** - Session metadata parsing
  - **ParseTranscript()** - Tool & task tracking from JSONL
  - **Tool categorization** - App/Internal/Custom/MCP/Skill
  - **Task tracking** - TodoWrite, TaskCreate, TaskUpdate processing

### Display Segments

Each segment implements:
```go
type Segment interface {
    ID() string
    Render(s *state.State, cfg *config.Config) (string, error)
    Enabled(cfg *config.Config) bool
}
```

**Segments:**
1. `model.go` (40 lines) - 🤖 model name
2. `context.go` (90 lines) - █▓▒░ gradient bar + 📥📤💾⚡ tokens
3. `git.go` (110 lines) - 🌿 branch + 📊 stats
4. `cost.go` (70 lines) - 💰 cost + ⏱ duration
5. `tools.go` (210 lines) - 📦🔌⚡🎨 categorized tools
6. `tasks.go` (200 lines) - Task dashboard or table
7. `agent.go` (45 lines) - 🤖 active agent
8. `ratelimit.go` (75 lines) - Rate limit tracking

### Output Rendering
- `output/renderer.go` (280 lines)
  - **renderMultiLine()** - Custom 4-line layout
  - **renderSingleLine()** - Compact horizontal layout
  - **renderContextBar()** - Gradient bar with percentage
  - **renderTokenDetails()** - Colored token breakdown
  - **renderFileChanges()** - +/- line changes

### Styling
- `style/style.go` (200 lines)
  - **Init()** - Theme color loading
  - **RenderGradientBar()** - Progress bars with smooth transitions
  - **RenderTable()** - Box-drawing table rendering
  - **13 semantic colors** - success, warning, danger, input, output, etc.

### Themes
- `theme/catppuccin.go` (180 lines)
  - **4 themes:** Macchiato (default), Mocha, Frappe, Latte
  - **LoadThemeFromConfig()** - Theme selection + color overrides
  - **ThemeWrapper** - Custom color override support

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

## Output Format

### Actual Output (Multi-line Layout)
```
Line 1: 🤖 Sonnet 4.5 │ █▓▒▒░░░░░░ 59%
Line 2: 📥 89k 📤 12k 💾 R:45k/W:23k ⚡ 200k │ 💰$0.0234  │  ⏱ 2m34s
Line 3: 🌿 main (dirty:2) │ 📝 +45/-12
Line 4: ╭─ 📦 App 23  🔌 MCP 2  ⚡ Skills 1 ─╮
```

**Note:** The README example showing "Context: 🧠" is **incorrect**. The actual output does NOT include a "Context:" label or 🧠 emoji.

## Important Implementation Details

### 1. Context Display (NO 🧠 emoji)
The context segment renders as:
```
█▓▒░░░░░░░ 59% 📥 89k 📤 12k 💾 R:45k/W:23k ⚡ 200k
```
**NOT:**
```
Context: 🧠 █▓▒░░░░░░░ 59%
```

### 2. Gradient Bar Colors
- Green (0-69%): Healthy usage
- Yellow (70-89%): Warning
- Red (90-100%): Danger

### 3. Task ID Indexing
- Claude Code uses **1-based task IDs** ("1", "2", "3")
- Parser converts to **0-based array indices**
- Bug fixed in v0.2.0 (commit f6f5fc4)

### 4. Skill Tracking
- Skills tracked by full name (e.g., "superpowers:using-git-worktrees")
- Extracted from Skill tool's `input.skill` parameter
- Bug fixed in v0.2.0 (commit f6f5fc4)

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
3. Run `make check` before committing
4. Ensure CI passes on all platforms

---

**Last Updated:** 2026-02-10
**Version:** v0.2.0
