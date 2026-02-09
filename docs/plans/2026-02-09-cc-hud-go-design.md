# cc-hud-go Design Document

**Date**: 2026-02-09
**Status**: Approved
**Goal**: Feature-complete Claude Code HUD replicating claude-hud functionality in Go

## Overview

`cc-hud-go` is a Go-based statusline tool for Claude Code that displays comprehensive session information. It integrates with the Claude Code statusline API to provide rich, real-time metrics through pure JSON output.

## Architecture

### High-Level Design

The tool follows a pipeline pattern with three main layers:

**Data Input Layer**: Two concurrent data sources feed the system:
1. **Stdin Reader**: Continuously reads JSON from stdin (Claude Code's real-time statusline API), parsing session state, model info, and context metrics
2. **Transcript Watcher**: Watches and tails the transcript JSONL file, extracting tool operations, agent activity, and task progress

**Data Aggregation**: A central `State` struct maintains the current view of all tracked metrics, updated atomically as new data arrives from either source. Git information is obtained by shelling out to `git` commands. Configuration is loaded once at startup from `~/.claude/cc-hud-go/config.json`.

**Output Generation**: When state changes, a segment builder constructs individual JSON objects for each configured display element. Each segment includes ANSI-styled text using lipgloss. Segments are written to stdout as a JSON array.

**Error Handling**: Missing data sources don't crash the program—graceful degradation ensures the tool always outputs valid JSON with whatever data is available.

The main loop is event-driven using Bubbletea's message-based architecture: messages arrive from stdin channel, file watcher channel, and config reload signal, outputting JSON whenever state changes.

## Core Data Structures

### State Management

The central `State` struct holds all current session data:

```go
type State struct {
    Model        ModelInfo
    Context      ContextInfo
    RateLimits   RateLimitInfo
    Git          GitInfo
    Tools        ToolsState
    Agents       AgentInfo
    Tasks        TaskInfo
    Session      SessionInfo
}

type ModelInfo struct {
    Name     string
    PlanType string // Pro/Max/Team
}

type ContextInfo struct {
    UsedTokens  int
    TotalTokens int
    Percentage  float64
}

type RateLimitInfo struct {
    HourlyUsed      int
    HourlyTotal     int
    SevenDayUsed    int
    SevenDayTotal   int
}

type GitInfo struct {
    Branch      string
    DirtyFiles  int
    Ahead       int
    Behind      int
    Added       int
    Modified    int
    Deleted     int
}

type ToolsState struct {
    AppTools      map[string]int  // Built-in Claude Code tools
    InternalTools map[string]int  // Internal system operations
    CustomTools   map[string]int  // User-defined tools
    MCPTools      map[MCPServer]map[string]int  // MCP server tools
    Skills        map[string]SkillUsage  // Skills invoked
}

type SkillUsage struct {
    Count    int
    LastUsed time.Time
    Duration time.Duration
}

type MCPServer struct {
    Name string
    Type string
}

type AgentInfo struct {
    ActiveAgent   string
    TaskDesc      string
    ElapsedTime   time.Duration
}

type TaskInfo struct {
    Pending     int
    InProgress  int
    Completed   int
}

type SessionInfo struct {
    StartTime  time.Time
    Duration   time.Duration
    TokenSpeed float64  // tokens/sec
}
```

### Configuration

The `Config` struct mirrors claude-hud's options:

```go
type Config struct {
    Preset             string  // full/essential/minimal
    LineLayout         string  // expanded/compact
    PathLevels         int     // 1-3
    ContextValue       string  // percentage/tokens
    SevenDayThreshold  int     // 0-100

    Display DisplayConfig
    Git     GitConfig
    Tools   ToolsConfig
}

type DisplayConfig struct {
    Model      bool
    Path       bool
    Context    bool
    Git        bool
    Tools      bool
    Agents     bool
    Tasks      bool
    RateLimits bool
    Duration   bool
    Speed      bool
}

type GitConfig struct {
    ShowBranch      bool
    ShowDirty       bool
    ShowAheadBehind bool
    ShowFileStats   bool
}

type ToolsConfig struct {
    GroupByCategory bool
    ShowTopN        int
    ShowSkills      bool
    ShowMCP         bool
}
```

### Segments

Each output segment implements a `Segment` interface:

```go
type Segment interface {
    ID() string
    Render(state *State, config *Config) (string, error)
    Enabled(config *Config) bool
}
```

Concrete implementations:
- `ModelSegment`: Model name with plan badge
- `PathSegment`: Current directory path (truncated)
- `ContextSegment`: Visual progress bar with color coding
- `GitSegment`: Branch with status indicators
- `ToolsSegment`: Categorized tool counts
- `AgentSegment`: Active agent info
- `TasksSegment`: Task progress tracker
- `RateLimitSegment`: Usage bars for Pro/Max/Team

### Input Parsers

Separate types for data ingestion:
- `StdinParser`: Parses Claude Code JSON from stdin
- `TranscriptParser`: Parses JSONL lines, categorizes tools, tracks timestamps

## Concurrent Processing

### Bubbletea Event Loop

Using Bubbletea's message-based architecture:

```go
type Model struct {
    state  *State
    config *Config
}

// Message types
type StdinMsg StdinUpdate
type TranscriptMsg TranscriptUpdate
type TickMsg time.Time
type ConfigReloadMsg Config

func (m Model) Init() tea.Cmd {
    return tea.Batch(
        readStdinCmd(),
        watchTranscriptCmd(),
        tea.Tick(time.Second, func(t time.Time) tea.Msg {
            return TickMsg(t)
        }),
    )
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case StdinMsg:
        m.state.ApplyStdinUpdate(msg)
        return m, tea.Batch(outputJSONCmd(m.state, m.config), readStdinCmd())

    case TranscriptMsg:
        m.state.ApplyTranscriptUpdate(msg)
        return m, tea.Batch(outputJSONCmd(m.state, m.config), watchTranscriptCmd())

    case TickMsg:
        m.state.UpdateDerived() // Update duration, etc.
        return m, tea.Batch(
            outputJSONCmd(m.state, m.config),
            tea.Tick(time.Second, func(t time.Time) tea.Msg {
                return TickMsg(t)
            }),
        )

    case ConfigReloadMsg:
        m.config = &msg
        return m, outputJSONCmd(m.state, m.config)
    }
    return m, nil
}

func (m Model) View() string {
    return "" // No TUI output
}
```

### Synchronization

State is owned by the Bubbletea model, eliminating race conditions. Updates flow through messages, ensuring thread-safe operations.

### Graceful Shutdown

SIGINT/SIGTERM handlers flush final output and clean up resources before exit.

## Enhanced Tool Tracking

The `TranscriptParser` categorizes tools by analyzing JSONL structure:

**Tool Classification**:
- **App Tools**: Built-in Claude Code tools (Read, Edit, Bash, etc.)
- **Internal Tools**: System operations (git commands, search)
- **Custom Tools**: User-defined tools
- **MCP Tools**: Pattern `mcp__<server>__<tool>` (e.g., `mcp__claude_ai_Atlassian__getConfluencePage`)
- **Skills**: Detected via `Skill` tool calls with skill names

**Parsing Logic**:
- Extract `tool_use` blocks with `name` field
- Apply pattern matching for categorization
- Track timestamps for usage metrics
- Calculate duration where possible

**Segment Display Options**:
- Compact: `"Tools: 26 (App:15 MCP:8 Skills:3)"`
- Expanded: Breakdown by category with top tools listed
- Configurable grouping and display depth

## Segment Rendering

### Segment Types

1. **ModelSegment**:
   - Example: `"claude-sonnet-4.5 [Pro]"`
   - Color based on model tier using lipgloss

2. **PathSegment**:
   - Example: `"~/code/cc-hud-go"`
   - Home directory shortening, configurable depth

3. **ContextSegment**:
   - Example: `"[████████░░] 45%"` (green) or `"[██████████] 95%"` (red)
   - Thresholds: <70% green, 70-90% yellow, >90% red

4. **GitSegment**:
   - Example: `"main ✗3 ↑2 ↓1"`
   - 3 dirty files, 2 commits ahead, 1 behind

5. **ToolsSegment**:
   - Compact: `"Tools: 26 (App:15 MCP:8 Skills:3)"`
   - Expanded: Top tools by category

6. **AgentSegment**:
   - Example: `"Agent: explore [2m15s] - Searching codebase"`

7. **TasksSegment**:
   - Example: `"Tasks: 2/5 ✓ (3 pending)"`

8. **RateLimitSegment**:
   - Example: `"Hourly: [███░░] 60% | 7-day: [██░░░] 45%"`

### JSON Output Format

```json
{
  "segments": [
    {
      "id": "model",
      "text": "\u001b[34mclaude-sonnet-4.5\u001b[0m [Pro]"
    },
    {
      "id": "context",
      "text": "[████░░] 67%"
    }
  ]
}
```

Disabled segments are omitted entirely from the output.

### Lipgloss Styling

```go
var (
    greenStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
    yellowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
    redStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
    blueStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
)
```

## Configuration System

### File Location

`~/.claude/cc-hud-go/config.json` (created with defaults on first run if missing)

### Structure

```json
{
  "preset": "full",
  "layout": "expanded",
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
  }
}
```

### Presets

Three built-in templates:
- **Full**: All segments enabled, expanded layout, maximum detail
- **Essential**: Core metrics only (model, context, git, tasks), compact layout
- **Minimal**: Just model and context bar, single line

### Loading & Validation

- Read config file on startup
- Parse failures trigger warning to stderr, fall back to defaults
- Watch config file for changes, reload and trigger re-render
- Validate ranges (pathLevels 1-3, threshold 0-100)
- Invalid values fall back to defaults with warning

## Error Handling

### Data Source Failures

- **No stdin**: Exit gracefully with error to stderr (statusline mode requires Claude Code input)
- **Missing transcript**: Continue with stdin-only data, log warning once, retry finding transcript every 30s
- **Malformed JSON/JSONL**: Skip bad lines, log to stderr, continue processing (never crash)
- **Transcript rotation**: Detect file replacement (new session), reset activity counters, reopen file handle

### Git Edge Cases

- **Not a git repo**: Omit git segment entirely
- **Detached HEAD**: Show commit SHA instead of branch name
- **Git command timeout**: Kill if >1s, skip git segment for 10s before retrying
- **Permission errors**: Skip git segment, log once to stderr

### Resource Management

- **File descriptor leaks**: Ensure transcript file is closed/reopened properly on errors
- **Memory growth**: Limit tool counter maps (cap at 1000 unique tools)
- **CPU usage**: Throttle transcript parsing if file grows extremely fast (>1000 lines/sec), batch updates

### Startup Race Conditions

- **Config file being written**: If parse fails, retry after 100ms, then fall back to defaults
- **Transcript not yet created**: Normal for new sessions, wait silently

### Output Guarantees

Always output valid JSON, even if all data sources fail: `{"segments": []}`

## Testing Strategy (TDD)

### Test Structure

Write tests before implementation for each layer:

**1. Parser Tests** (`*_test.go` files):
- `TestStdinParser`: Parse valid Claude Code JSON, handle malformed input, extract all fields
- `TestTranscriptParser`: Parse JSONL lines, identify tool types, track timestamps, handle partial lines
- `TestConfigParser`: Load valid configs, handle missing files, validate preset application

**2. State Tests**:
- `TestStateUpdate`: Apply stdin/transcript updates, verify state consistency
- `TestToolsCategorization`: Correctly classify tools (app/internal/custom/MCP/skills)
- `TestGracefulDegradation`: Missing data sources, state remains valid

**3. Segment Tests**:
- `TestSegmentRendering`: Each segment type renders correct output for various states
- `TestContextColorThresholds`: Green <70%, yellow 70-90%, red >90%
- `TestSegmentFiltering`: Disabled segments omitted, enabled ones appear

**4. Integration Tests**:
- `TestEndToEnd`: Feed sample stdin + transcript data, verify JSON output
- `TestConcurrency`: Multiple rapid updates don't race or corrupt state
- `TestFileWatching`: Transcript file changes trigger updates

### Test Data

Create `testdata/` directory with:
- Sample stdin JSON files
- Transcript JSONL excerpts
- Config files for realistic testing

### Mocking

Use interfaces for testability:
- `io.Reader` for stdin
- `fs.FS` for config
- Mock git commands via interface

## Implementation Phases

### Project Structure

```
cc-hud-go/
├── main.go                 # Entry point, Bubbletea program
├── config/
│   ├── config.go          # Config loading, validation, presets
│   └── config_test.go
├── parser/
│   ├── stdin.go           # Claude Code JSON parser
│   ├── transcript.go      # JSONL parser, tool categorization
│   └── parser_test.go
├── state/
│   ├── state.go           # State struct, update methods
│   └── state_test.go
├── segment/
│   ├── segment.go         # Segment interface
│   ├── model.go           # ModelSegment
│   ├── context.go         # ContextSegment
│   ├── git.go             # GitSegment
│   ├── tools.go           # ToolsSegment
│   ├── agent.go           # AgentSegment
│   ├── tasks.go           # TasksSegment
│   ├── ratelimit.go       # RateLimitSegment
│   └── segment_test.go
├── output/
│   ├── renderer.go        # JSON output generation
│   └── renderer_test.go
├── internal/
│   ├── git/               # Git command wrappers
│   ├── watcher/           # File watching logic
│   └── ansi/              # ANSI color helpers (if needed)
└── testdata/              # Sample inputs for testing
```

### TDD Build Order

Write tests first for each phase:

1. **Phase 1**: Config system (load, validate, presets)
2. **Phase 2**: State management (struct, updates)
3. **Phase 3**: Stdin parser (JSON → state)
4. **Phase 4**: Basic segments (model, path, context)
5. **Phase 5**: Git integration (commands, parsing, segment)
6. **Phase 6**: Transcript parser (JSONL, tool categorization)
7. **Phase 7**: Advanced segments (tools, agents, tasks, rate limits)
8. **Phase 8**: Bubbletea orchestration (event loop, messages)
9. **Phase 9**: Integration tests & polish

## Dependencies

### Core Libraries

- **bubbletea**: Event loop and message-based architecture
- **lipgloss**: ANSI styling and colors
- **bubbles**: Optional components (progress bars, spinners)

### Standard Library

- `encoding/json`: Config and output JSON marshaling
- `bufio`: Stdin scanning, file reading
- `os/exec`: Git command execution
- `path/filepath`: Path operations
- `os`: File I/O, environment, signals
- `time`: Timestamps, durations
- `strings`, `fmt`: String manipulation
- `regexp`: Tool name pattern matching
- `io`, `io/fs`: File abstractions for testability

### Development Tools

- `golangci-lint`: Linting (as specified in CLAUDE.md)
- `go test -cover`: Test coverage (aim for >80% on core logic)
- `go mod`: Dependency management

### Build

- Single static binary: `go build -o cc-hud-go .`
- No external runtime dependencies
- Cross-compile for multiple platforms if needed

## Success Criteria

- [ ] Feature parity with claude-hud (all metrics tracked)
- [ ] Enhanced tool tracking (app/internal/custom/MCP/skills)
- [ ] Pure JSON output for Claude Code statusline
- [ ] Graceful degradation (no crashes on missing data)
- [ ] Event-driven updates (real-time responsiveness)
- [ ] Config file support with presets
- [ ] >80% test coverage on core logic
- [ ] Clean, idiomatic Go following Charm ecosystem style

## References

- Claude Code statusline API: https://code.claude.com/docs/en/statusline
- Original claude-hud: https://github.com/jarrodwatts/claude-hud
- Oh My Posh Claude segment: https://ohmyposh.dev/docs/segments/cli/claude
- Bubbletea: https://github.com/charmbracelet/bubbletea
- Lipgloss: https://github.com/charmbracelet/lipgloss
- Bubbles: https://github.com/charmbracelet/bubbles
