# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

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

# Using go directly
go test ./...
go test -cover ./...

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

# Manual commands
go fmt ./...
go vet ./...
golangci-lint run
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

## Architecture

### Project Structure

```
cc-hud-go/
├── config/          # Configuration management with presets
│   ├── config.go
│   └── config_test.go
├── state/           # Session state tracking and derived fields
│   ├── state.go
│   └── state_test.go
├── parser/          # Dual input parsing (stdin JSON & transcript JSONL)
│   ├── parser.go
│   ├── stdin_test.go
│   ├── transcript_test.go
│   └── tasks_test.go
├── segment/         # Modular display segments
│   ├── segment.go   # Segment interface & registry
│   ├── model.go     # Model and plan type display
│   ├── context.go   # Token usage and context window
│   ├── git.go       # Git branch, status, file stats
│   ├── cost.go      # Cost tracking and code metrics
│   ├── tools.go     # Tool usage categorization
│   ├── tasks.go     # Task progress tracking
│   ├── agent.go     # Active agent and task info
│   ├── ratelimit.go # API rate limit monitoring
│   └── *_test.go
├── output/          # JSON output renderer for statusline API
│   ├── renderer.go
│   └── renderer_test.go
├── style/           # Lipgloss styling with semantic color system
│   └── style.go
├── version/         # Version detection and build info
│   ├── version.go
│   └── version_test.go
├── internal/
│   ├── git/         # Git integration via command execution
│   │   ├── git.go
│   │   └── git_test.go
│   └── watcher/     # File watching utilities
│       └── watcher.go
├── testdata/        # Test fixtures and sample data
├── docs/            # Documentation and planning
│   ├── plans/       # Design and implementation plans
│   └── COLOR_SCHEME.md
├── assets/          # Screenshots and preview images
├── main.go          # Application entry point
├── main_test.go     # Main package tests
├── integration_test.go  # Integration tests
├── justfile         # Build and development commands
└── go.mod
```

### Claude Code Statusline Integration

The tool implements the Claude Code statusline protocol, which expects:
- JSON output written to stdout (rendered by `output/renderer.go`)
- Specific data structure for statusline information
- Real-time updates based on Claude Code session state via stdin

**Data Flow:**
1. Read JSON session data from stdin (provided by Claude Code)
2. Parse transcript file for tool usage tracking
3. Fetch git information from current repository
4. Update state with derived fields (percentages, durations)
5. Render enabled segments based on configuration
6. Output formatted JSON to stdout

### Key Components

**Segments** - Modular display components implementing the `Segment` interface:
```go
type Segment interface {
    ID() string
    Render(s *state.State, cfg *config.Config) (string, error)
    Enabled(cfg *config.Config) bool
}
```

Available segments:
- **ModelSegment** - Current Claude model and plan type
- **ContextSegment** - Token usage with color-coded thresholds
- **GitSegment** - Branch, dirty files, ahead/behind, file stats
- **CostSegment** - Cost tracking, duration, lines changed
- **ToolsSegment** - Tool usage categorized by type (App/MCP/Skills/Custom)
- **TasksSegment** - Task completion progress
- **AgentSegment** - Active agent name and current task
- **RateLimitSegment** - 7-day API usage tracking

**State Management** - Centralized session state with automatic derived field calculation:
- Context percentage calculation
- Session duration tracking
- Tool usage categorization
- Task progress aggregation

**Parser System** - Dual parsing approach:
- **Stdin parser** - Session data from Claude Code (JSON)
- **Transcript parser** - Tool usage tracking from JSONL file

**Style System** - Semantic color palette using Lipgloss:
- Status colors (green/yellow/red for thresholds)
- Flow colors (blue for input, emerald for output)
- Cache colors (purple for reads, pink for writes)
- Primary UI colors (purple, cyan, orange)
- TrueColor support with forced color output

**Configuration** - Three preset modes with granular control:
- **Full** - All features enabled
- **Essential** - Core metrics only
- **Minimal** - Minimal information
- Per-segment enable/disable options
- Customizable thresholds and display formats

### Design Principles

Following the Charm ecosystem style:
- Elegant, minimal terminal UI with Lipgloss styling
- Composable segment architecture
- Clean separation between state, rendering, and configuration
- Graceful degradation (missing config → defaults)
- Comprehensive test coverage with TDD approach
- Semantic color system with meaningful associations

## Claude Code Statusline API

The tool displays comprehensive session information including:

**Core Metrics:**
- Current model being used (e.g., claude-sonnet-4.5)
- Plan type (Free, Pro, Team, Enterprise)
- Token usage statistics (used/total, percentage, color-coded thresholds)
- Context window tracking (input, output, cache read, cache create)

**Development Insights:**
- Git branch, dirty files, ahead/behind status
- File statistics (added, modified, deleted)
- Tool usage categorization (App tools, MCP tools, Skills, Custom)
- Task progress tracking (pending, in-progress, completed)
- Active agent name and current task description

**Session Tracking:**
- Cost tracking in USD
- Session duration and API duration
- Code changes (lines added/removed)
- Token processing speed
- API rate limits (hourly and 7-day)

Refer to the official docs for the complete API: https://code.claude.com/docs/en/statusline

## Testing & Quality

**Running Tests:**
- Unit tests for all segments (`segment/*_test.go`)
- Parser tests for stdin and transcript parsing (`parser/*_test.go`)
- State management tests (`state/state_test.go`)
- Config validation tests (`config/config_test.go`)
- Integration tests (`integration_test.go`)

**Test Data:**
- Sample session data in `testdata/`
- Fixture files for transcript parsing
- Mock git repositories for testing

**Code Quality:**
- Comprehensive test coverage with TDD approach
- Linting with golangci-lint
- Go vet for static analysis
- Consistent formatting with go fmt

## Contributing

When adding new features:

1. **New Segments** - Create in `segment/<name>.go` with tests
2. **State Fields** - Add to appropriate struct in `state/state.go`
3. **Configuration** - Update `config/config.go` if needed
4. **Styling** - Use semantic colors from `style/style.go`
5. **Tests** - Write tests before implementation (TDD)
6. **Documentation** - Update both CLAUDE.md and README.md
