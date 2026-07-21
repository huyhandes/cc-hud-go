# AGENTS.md

This file provides guidance to AI agents working with code in this repository.

## Project Overview

`cc-hud-go` is a Go-based statusline tool for Claude Code that displays helpful information about the current Claude Code session. It integrates with the Claude Code statusline API (https://code.claude.com/docs/en/statusline) to provide rich, real-time information.

**Inspiration sources:**
- https://github.com/jarrodwatts/claude-hud (original HUD implementation)
- https://ohmyposh.dev/docs/segments/cli/claude (Oh My Posh Claude segment)
- https://github.com/charmbracelet/bubbletea (TUI framework style)
- https://github.com/charmbracelet/gum (CLI interaction style)

## Development Commands

**Using Just (Recommended):**
```bash
# Show all available commands
just

# Build with version info from git tags
just build

# Run tests
just test

# Run tests with coverage
just test-coverage

# Format, vet, and test
just check

# Clean build artifacts
just clean

# Build and install to ~/.local/bin
just install
```

**Manual Build:**
```bash
# Build without version info (shows git tag or "dev")
go build -o cc-hud-go .

# Build with specific version
go build -ldflags "-X github.com/huyhandes/cc-hud-go/version.Version=v1.0.0" -o cc-hud-go .
```

**Run:**
```bash
# Normal run (expects stdin)
go run .

# Show help
go run . --help
go run . -h

# Show version (auto-detects from git tags)
go run . --version
go run . -v
```

**Tests:**
```bash
# Run all tests (using just)
just test

# Run tests with coverage
just test-coverage

# Run a specific test
go test -run TestName ./path/to/package

# Run tests with verbose output
go test -v ./...
```

**Linting:**
```bash
# Format code
just fmt

# Run go vet
just vet

# Run golangci-lint (if installed)
just lint

# Run all checks (format, vet, test)
just check
```

**Dependencies:**
```bash
# Add a dependency
go get github.com/package/name

# Update dependencies
go get -u ./...

# Tidy dependencies
go mod tidy
```

## Workflow

After finishing any task, always run:

```bash
just fmt
just lint
```

## Architecture

### Project Structure

```
cc-hud-go/
├── state/           # Session state tracking and derived fields
│   ├── state.go
│   └── state_test.go
├── parser/          # Stdin JSON parsing (single concern)
│   ├── stdin.go              # StdinData type, ParseStdin()
│   └── stdin_test.go
├── segment/         # Modular display segments
│   ├── segment.go   # Segment interface, All(), ByID() registry
│   ├── model.go     # Model name display
│   ├── caveman.go   # Caveman mode badge (self-gating)
│   ├── ponytail.go  # Ponytail mode badge (self-gating)
│   ├── git.go       # Git branch, status, file stats
│   ├── fivehour.go  # 5-hour rate limit tracking
│   ├── ratelimit.go # 7-day API rate limit tracking
│   └── *_test.go
├── output/          # Output renderer for statusline API
│   ├── renderer.go
│   └── renderer_test.go
├── format/          # Shared formatting helpers (DRY)
│   ├── format.go    # Tokens(), Duration(), Cost()
│   └── format_test.go
├── style/           # Lipgloss styling (split per file)
│   ├── style.go     # Colors, Init(), ThresholdColor()
│   ├── gradient.go  # RenderGradientBar(), color interpolation
│   └── gradient_test.go
├── theme/           # Theme system (Catppuccin variants)
│   ├── theme.go
│   ├── catppuccin.go
│   └── theme_test.go
├── version/         # Version detection and build info
│   ├── version.go
│   └── version_test.go
├── internal/
│   ├── git/         # Git integration via command execution
│   │   ├── git.go
│   │   └── git_test.go
│   └── oauth/       # OAuth helpers (always fetched, no stdin fallback)
│       ├── oauth.go
│       └── oauth_test.go
├── cmd/             # Developer utility binaries
│   ├── test-gradient/
│   └── test-oauth/
├── testdata/        # Test fixtures and sample data
├── docs/            # Documentation
│   ├── adr/         # Architectural Decision Records
│   ├── CODEMAP.md
│   ├── PROJECT_STATUS.md
│   ├── BUILD_GUIDE.md
│   └── plans/
├── assets/          # Screenshots and preview images
├── main.go          # Application entry point
├── main_test.go     # Main package tests
├── integration_test.go  # Integration tests
├── justfile         # Build and development commands
└── go.mod
```

Notes:
- No `config/` directory. Defaults are hardcoded where they are used; the only user-facing knob is the `--theme` CLI flag. See [docs/adr/0001-default-only-mode.md](docs/adr/0001-default-only-mode.md) and [docs/adr/0002-theme-via-cli-flag.md](docs/adr/0002-theme-via-cli-flag.md).
- No `examples/` directory. There is no config file to demonstrate.
- `parser/` only contains `stdin.go` + `stdin_test.go`. Transcript JSONL parsing is gone. See [docs/adr/0004-no-transcript-parsing.md](docs/adr/0004-no-transcript-parsing.md).
- `style/` has no `table.go`. `segment/` has no `tools.go`/`agent.go`/`tasks.go`/`context.go`/`cost.go`.

### Claude Code Statusline Integration

The tool implements the Claude Code statusline protocol, which expects:
- JSON output written to stdout (rendered by `output/renderer.go`)
- Specific data structure for statusline information
- Real-time updates based on Claude Code session state via stdin

**Data Flow:**
1. Read JSON session data from stdin (provided by Claude Code)
2. Fetch git information from current repository
3. Fetch OAuth rate-limit data on every refresh (no stdin fallback). See [docs/adr/0003-always-fetch-oauth.md](docs/adr/0003-always-fetch-oauth.md).
4. Update state with derived fields (percentages, durations)
5. Render segments (self-gating: a segment returns `""` to hide itself)
6. Output formatted JSON to stdout

### Key Components

**Segments** - Modular display components implementing the `Segment` interface. Surviving segments:
- **ModelSegment** - Current Claude model name
- **CavemanSegment** - Caveman mode badge (renders only when `.caveman-active` flag file exists)
- **PonytailSegment** - Ponytail mode badge (renders only when `.ponytail-active` flag file exists)
- **GitSegment** - Branch, dirty files, ahead/behind, file stats
- **FiveHourSegment** - 5-hour rate limit tracking
- **RateLimitSegment** - 7-day API usage tracking

Segments self-gate by returning an empty string when they have nothing to show; that is the only on/off mechanism.

**Format Package** - Shared formatting helpers eliminating duplication:
- `format.Tokens()` - Token count formatting (e.g. 5000 -> "5k")
- `format.Duration()` - Milliseconds to human-readable duration
- `format.Cost()` - USD cost formatting

**State Management** - Centralized session state with automatic derived field calculation.

**Parser System** - Single concern:
- `stdin.go` - Session data from Claude Code (JSON)

**Style System** - Semantic color palette using Lipgloss:
- `style.go` - Theme integration, ThresholdColor() helper
- `gradient.go` - Gradient progress bars with smooth color transitions

**Theme System** - 4 Catppuccin variants. Source of truth: `theme/catppuccin.go`. Selected via the `--theme` CLI flag (default `macchiato`; values `macchiato`/`mocha`/`frappe`/`latte`; unknown values fall back to macchiato).

**Configuration** - No config file. The only user-facing knob is `--theme`. Defaults are hardcoded.

### Design Principles

Following the Charm ecosystem style:
- Elegant, minimal terminal UI with Lipgloss styling
- Composable segment architecture
- Clean separation between state, rendering, and output
- Opinionated defaults; customization via CLI flags only
- Self-gating segments: empty-string return is the on/off mechanism
- Semantic color system with meaningful associations

## Claude Code Statusline API

The tool displays session information including:

**Core:**
- Current model being used (e.g., claude-sonnet-4.5)
- Git branch, dirty files, ahead/behind status, file statistics
- API rate limits (5-hour and 7-day), fetched live from the OAuth API
- Mode badges (caveman / ponytail) when active

Refer to the official docs for the complete API: https://code.claude.com/docs/en/statusline

## Testing & Quality

**Running Tests:**
- Unit tests for all segments (`segment/*_test.go`)
- Parser tests for stdin parsing (`parser/*_test.go`)
- State management tests (`state/state_test.go`)
- Integration tests (`integration_test.go`)

**Test Data:**
- Sample session data in `testdata/`
- Mock git repositories for testing

**Code Quality:**
- Comprehensive test coverage with TDD approach
- Linting with golangci-lint
- Go vet for static analysis
- Consistent formatting with go fmt

## Contributing

When adding features:

1. **New Segments** - Create in `segment/<name>.go` with tests; self-gate via empty-string return
2. **State Fields** - Add to appropriate struct in `state/state.go`
3. **New user-facing knobs** - Prefer a CLI flag over a config key (see ADRs 0001 and 0002)
4. **Styling** - Use semantic colors from `style/style.go`
5. **Tests** - Write tests before implementation (TDD)
6. **Documentation** - Update both AGENTS.md and README.md
