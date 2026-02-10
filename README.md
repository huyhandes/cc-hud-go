# cc-hud-go

A Go-based statusline tool for [Claude Code](https://code.claude.com) that displays rich, real-time information about your current Claude Code session.

[![CI](https://github.com/huyhandes/cc-hud-go/actions/workflows/ci.yml/badge.svg)](https://github.com/huyhandes/cc-hud-go/actions/workflows/ci.yml)
[![Release](https://github.com/huyhandes/cc-hud-go/actions/workflows/release.yml/badge.svg)](https://github.com/huyhandes/cc-hud-go/actions/workflows/release.yml)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.24-blue.svg)](https://golang.org)
[![Go Report Card](https://goreportcard.com/badge/github.com/huyhandes/cc-hud-go)](https://goreportcard.com/report/github.com/huyhandes/cc-hud-go)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Latest Release](https://img.shields.io/github/v/release/huyhandes/cc-hud-go)](https://github.com/huyhandes/cc-hud-go/releases/latest)

![Preview](assets/preview.jpeg)

## Features

### 📊 Real-time Metrics
- **Model Information** - Current Claude model and plan type
- **Context Usage** - Token usage with color-coded thresholds (green/yellow/red)
- **Rate Limits** - 7-day API usage tracking with visual warnings
- **Cost Tracking** - Session cost (USD), duration, and code changes (lines added/removed)
- **Session Stats** - Duration and token processing speed

### 🔧 Development Insights
- **Git Integration** - Branch name, dirty files, ahead/behind status, file stats
- **Tool Tracking** - Categorized tool usage (App/Internal/Custom/MCP/Skills)
- **Task Progress** - Task completion tracking (completed/total)
- **Agent Activity** - Active agent name and current task description

### 🎨 Beautiful Themes & Visuals
- **Catppuccin Color Palettes** - 4 carefully curated themes (Macchiato, Mocha, Frappe, Latte)
- **Gradient Progress Bars** - Smooth color transitions with Unicode characters (█▓▒░)
- **Smart Adaptive Layouts** - Automatic table rendering for detailed views
- **Custom Color Overrides** - Override any semantic color with your preferences

### ⚙️ Flexible Configuration
- **Multiple Presets** - Full, Essential, and Minimal display modes
- **Granular Control** - Enable/disable individual segments
- **Customizable Thresholds** - Configure warning levels for context and rate limits
- **JSON Configuration** - Easy configuration via `~/.claude/cc-hud-go/config.json`

## Installation

### Pre-built Binaries

Download pre-built binaries from the [latest release](https://github.com/huyhandes/cc-hud-go/releases/latest):

```bash
# Linux (amd64)
wget https://github.com/huyhandes/cc-hud-go/releases/latest/download/cc-hud-go-linux-amd64.tar.gz
tar -xzf cc-hud-go-linux-amd64.tar.gz
sudo mv cc-hud-go-linux-amd64 /usr/local/bin/cc-hud-go

# macOS (Apple Silicon)
wget https://github.com/huyhandes/cc-hud-go/releases/latest/download/cc-hud-go-darwin-arm64.tar.gz
tar -xzf cc-hud-go-darwin-arm64.tar.gz
sudo mv cc-hud-go-darwin-arm64 /usr/local/bin/cc-hud-go

# macOS (Intel)
wget https://github.com/huyhandes/cc-hud-go/releases/latest/download/cc-hud-go-darwin-amd64.tar.gz
tar -xzf cc-hud-go-darwin-amd64.tar.gz
sudo mv cc-hud-go-darwin-amd64 /usr/local/bin/cc-hud-go
```

Available builds:
- Linux: `amd64`, `arm64`
- macOS: `amd64`, `arm64`
- Windows: `amd64`, `arm64`

### From Source

```bash
# Clone the repository
git clone git@github.com:huyhandes/cc-hud-go.git
cd cc-hud-go

# Build with version info (using Make - recommended)
make build

# Or build manually
go build -o cc-hud-go .

# Install to GOPATH/bin (optional)
make install
# Or move to PATH manually
sudo mv cc-hud-go /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/huyhandes/cc-hud-go@latest
```

## Themes & Visual Customization

cc-hud-go features beautiful [Catppuccin](https://github.com/catppuccin/catppuccin) color palettes with gradient progress bars and smart adaptive layouts.

### Available Themes

**🌙 Macchiato** (default) - Dark theme with purple accents
```json
{ "theme": "macchiato" }
```

**🌑 Mocha** - Darkest variant with rich, deep colors
```json
{ "theme": "mocha" }
```

**🌆 Frappe** - Medium-dark with warmer tones
```json
{ "theme": "frappe" }
```

**☀️ Latte** - Light theme for bright environments
```json
{ "theme": "latte" }
```

### Visual Features

**Gradient Progress Bars** - Smooth color transitions using Unicode block characters:
```
█▓▒░░░░░░░ 35% 📥 89k 📤 12k   ← Green (healthy)
█▓▓▓▒▒▒░░░ 75% 📥 150k 📤 25k  ← Yellow (warning)
█▓▓▓▓▓▒▒▒░ 95% 📥 190k 📤 38k  ← Red (danger)
```

**Multi-line Layout** - Clean 4-line display grouping related metrics:
```
Line 1: 🤖 Sonnet 4.5 │ █▓▒▒░░░░░░ 59%
Line 2: 📥 89k 📤 12k 💾 R:45k/W:23k ⚡ 200k │ 💰$0.0234  │  ⏱ 2m34s
Line 3: 🌿 main (dirty:2) │ 📝 +45/-12
Line 4: ╭─ 📦 App 23  🔌 MCP 2  ⚡ Skills 1 ─╮
```

**Smart Adaptive Layouts** - Automatic switching between inline lipgloss boxes and table views:
- Below threshold: Compact inline display with styled boxes (╭╮╰╯)
- Above threshold: Detailed table view with box-drawing characters (┌┬┐)
- Configurable thresholds per segment type (default: 999 for lipgloss boxes)

**Semantic Colors** - Each element uses meaningful color associations:
- Input tokens: Blue 📥 | Output tokens: Emerald 📤
- Cache reads: Purple 💾 | Cache writes: Pink
- Success: Green | Warning: Yellow | Danger: Red

### Custom Colors

Override any semantic color while keeping the base theme:

```json
{
  "theme": "macchiato",
  "colors": {
    "success": "#00ff00",
    "warning": "#ffaa00",
    "danger": "#ff0000",
    "primary": "#00aaff"
  }
}
```

**Available semantic colors:**
- `success` - Completed states, positive indicators (green)
- `warning` - Warnings, medium thresholds (yellow/orange)
- `danger` - Errors, high thresholds (red)
- `input` - Input tokens (blue)
- `output` - Output tokens (emerald)
- `cacheRead` - Cache read tokens (purple)
- `cacheWrite` - Cache write tokens (pink)
- `primary` - Main brand color (purple)
- `highlight` - Highlights and accents (cyan)
- `accent` - Secondary accents (orange)
- `muted` - Borders, subtle elements (gray)
- `bright` - Bright text (white/cream)
- `info` - Informational elements (teal)

### Example Configs

Pre-configured examples are available in the [`examples/`](examples/) directory:

```bash
# Copy Macchiato theme (default)
cp examples/config-macchiato.json ~/.claude/cc-hud-go/config.json

# Copy Mocha theme (darkest)
cp examples/config-mocha.json ~/.claude/cc-hud-go/config.json

# Copy custom colors example
cp examples/config-custom-colors.json ~/.claude/cc-hud-go/config.json
```

See [`examples/README.md`](examples/README.md) for detailed theme documentation and customization guide.

## Usage

### Integration with Claude Code

Add to your Claude Code statusline configuration:

```json
{
  "statusline": {
    "command": "cc-hud-go"
  }
}
```

The tool reads session data from stdin (provided by Claude Code) and outputs formatted JSON to stdout.

### Standalone Usage

For testing or development:

```bash
# View help (shows usage, configuration, examples)
cc-hud-go --help
cc-hud-go -h

# Check version (auto-detects from git tags or shows release version)
cc-hud-go --version
cc-hud-go -v

# Test with sample stdin data
echo '{"model":"claude-sonnet-4.5","context":{"used":5000,"total":10000}}' | cc-hud-go
```

**Version Information:**
- Release builds: Shows the tagged version (e.g., `v0.1.0`)
- Development builds: Auto-detects from `git describe` (e.g., `v0.1.0-dirty`)
- Without git: Falls back to `dev`

The `--help` flag displays comprehensive usage information including:
- Command syntax and description
- Available command-line options
- Configuration file location and presets
- Integration instructions for Claude Code
- Usage examples
- Links to documentation and issue tracker

## Configuration

### Configuration File

Create `~/.claude/cc-hud-go/config.json`:

```json
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
```

### Presets

**Full** (default) - All features enabled
```json
{
  "preset": "full"
}
```

**Essential** - Core metrics only
```json
{
  "preset": "essential"
}
```

**Minimal** - Minimal information
```json
{
  "preset": "minimal"
}
```

### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `theme` | string | `"macchiato"` | Color theme: `macchiato`, `mocha`, `frappe`, or `latte` |
| `colors` | object | `{}` | Custom color overrides (hex codes) |
| `preset` | string | `"full"` | Preset configuration: `full`, `essential`, or `minimal` |
| `lineLayout` | string | `"expanded"` | Layout style: `expanded` or `compact` |
| `pathLevels` | int | `2` | Number of directory levels to show (1-3) |
| `contextValue` | string | `"percentage"` | Context display format |
| `sevenDayThreshold` | int | `80` | Warning threshold for 7-day rate limit (0-100) |

#### Display Options

All boolean flags to enable/disable segments:
- `model` - Show model name and plan type
- `path` - Show current working directory
- `context` - Show token usage
- `git` - Show git information
- `tools` - Show tool usage statistics
- `agents` - Show active agent information
- `tasks` - Show task progress
- `rateLimits` - Show API rate limit usage
- `duration` - Show session duration
- `speed` - Show token processing speed

#### Git Options

- `showBranch` - Display current git branch
- `showDirty` - Show count of dirty files
- `showAheadBehind` - Show commits ahead/behind remote
- `showFileStats` - Show added/modified/deleted file counts

#### Tools Options

- `groupByCategory` - Group tools by category (App/MCP/Skills/Custom)
- `showTopN` - Number of top tools to display (0 = all)
- `showSkills` - Include skill usage in tool counts
- `showMCP` - Include MCP tool usage in tool counts

#### Table Options

Smart adaptive rendering thresholds (switches from inline lipgloss boxes to table view):

- `toolsTableThreshold` - Tool count threshold for table view (default: 999)
- `tasksTableThreshold` - Task count threshold for table view (default: 999)
- `contextTableThreshold` - Context size threshold for table view (default: 999)

**Default behavior (threshold: 999):** Always use styled lipgloss boxes with colored borders (╭╮╰╯)

**Table view (when count > threshold):** Switch to plain table format with box-drawing characters (┌┬┐) featuring:
- Sortable columns
- Filtered display (e.g., last 3 completed tasks)
- Truncated long names for better readability

Lower thresholds (e.g., 5 or 3) to switch to table view earlier for very large datasets.

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
│   └── transcript_test.go
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
│   ├── style.go
│   └── table_test.go
├── theme/           # Theme system with Catppuccin palettes
│   ├── theme.go
│   ├── catppuccin.go
│   └── theme_test.go
├── examples/        # Example configuration files
│   ├── README.md
│   └── config-*.json
├── internal/
│   ├── git/         # Git integration via command execution
│   │   ├── git.go
│   │   └── git_test.go
│   └── watcher/     # File watching utilities
│       └── watcher.go
├── testdata/        # Test fixtures and sample data
├── docs/            # Documentation and planning
│   └── COLOR_SCHEME.md
├── assets/          # Screenshots and preview images
├── main.go          # Application entry point
└── go.mod
```

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
- `ModelSegment` - Current Claude model and plan type
- `ContextSegment` - Token usage with color-coded thresholds
- `GitSegment` - Branch, dirty files, ahead/behind, file stats
- `CostSegment` - Cost tracking, duration, lines changed
- `ToolsSegment` - Tool usage categorized by type (App/MCP/Skills/Custom)
- `TasksSegment` - Task completion progress
- `AgentSegment` - Active agent name and current task
- `RateLimitSegment` - 7-day API usage tracking

**State** - Centralized session state with automatic derived field calculation:
- Context percentage calculation
- Session duration tracking
- Tool usage categorization
- Task progress aggregation

**Parser** - Dual parsing system:
- Stdin parser for Claude Code session data (JSON)
- Transcript parser for tool usage tracking (JSONL)

**Theme System** - Catppuccin color palettes with customization:
- 4 beautiful themes: Macchiato, Mocha, Frappe, Latte
- Custom color override support
- Theme interface for extensibility

**Style System** - Semantic color palette using Lipgloss:
- Status colors (green/yellow/red for thresholds)
- Flow colors (blue for input, emerald for output)
- Cache colors (purple for reads, pink for writes)
- Primary UI colors (purple, cyan, orange)
- Gradient progress bars with smooth color transitions
- Smart table rendering with box-drawing characters
- TrueColor support with forced color output

**Renderer** - JSON output formatter for Claude Code statusline API

### Design Principles

Built with the [Charm](https://charm.sh) ecosystem and [Catppuccin](https://github.com/catppuccin/catppuccin) color palettes:
- Clean, elegant terminal styling with [Lipgloss](https://github.com/charmbracelet/lipgloss)
- Beautiful themes with Catppuccin color palettes
- Gradient progress bars and smart adaptive layouts
- Composable segment architecture
- Clean separation between state, rendering, and configuration
- Graceful degradation (missing config → defaults)
- Comprehensive test coverage with TDD approach
- Semantic color system with meaningful associations
- Modular design inspired by [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Gum](https://github.com/charmbracelet/gum) patterns

## Development

### Prerequisites

- Go 1.24 or higher
- Git

### Setup

```bash
# Clone repository
git clone git@github.com:huyhandes/cc-hud-go.git
cd cc-hud-go

# Install dependencies
go mod download

# Run tests
make test

# Build with version info
make build

# Or build manually
go build -o cc-hud-go .
```

### Make Commands

The project includes a Makefile for common development tasks:

```bash
make help            # Show all available commands
make build           # Build with version from git tags
make test            # Run all tests
make test-coverage   # Run tests with coverage
make check           # Run fmt, vet, and test
make fmt             # Format code
make vet             # Run go vet
make lint            # Run golangci-lint
make clean           # Remove build artifacts
make install         # Install to GOPATH/bin
make build-all       # Build for all platforms
```

### Creating a Release

The project uses GitHub Actions for automated releases:

```bash
# Create and push a version tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

This will trigger the release workflow which:
1. Builds binaries for all supported platforms (Linux, macOS, Windows)
2. Supports both `amd64` and `arm64` architectures
3. Creates compressed archives (`.tar.gz` for Unix, `.zip` for Windows)
4. Generates SHA256 checksums
5. Creates a GitHub release with auto-generated release notes
6. Uploads all artifacts to the release

You can also trigger a release manually from the Actions tab on GitHub.

### Running Tests

```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# Integration tests
go test -tags=integration -v .

# Specific package
go test ./segment -v

# Specific test
go test -run TestModelSegment ./segment
```

### Code Quality

```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Lint (requires golangci-lint)
golangci-lint run
```

## Inspiration

This project draws inspiration from:
- [claude-hud](https://github.com/jarrodwatts/claude-hud) - Original HUD implementation
- [Oh My Posh Claude Segment](https://ohmyposh.dev/docs/segments/cli/claude) - Claude segment for Oh My Posh
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework philosophy
- [Gum](https://github.com/charmbracelet/gum) - CLI interaction patterns

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`go test ./...`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Adding New Segments

1. Create `segment/<name>.go` implementing the `Segment` interface
2. Add corresponding test file `segment/<name>_test.go`
3. Register in `segment/segment.go` `All()` function
4. Add configuration option in `config/config.go` if needed
5. Update README with new segment documentation

## License

MIT License - see [LICENSE](LICENSE) file for details

## Links

- [Claude Code Documentation](https://code.claude.com/docs)
- [Claude Code Statusline API](https://code.claude.com/docs/en/statusline)
- [Issue Tracker](https://github.com/huyhandes/cc-hud-go/issues)

---

Built with ❤️ using [Go](https://golang.org) and [Charm](https://charm.sh)
