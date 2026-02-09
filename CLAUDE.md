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

**Build:**
```bash
go build -o cc-hud-go .
```

**Run:**
```bash
go run .
```

**Tests:**
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run a specific test
go test -run TestName ./path/to/package

# Run tests with verbose output
go test -v ./...
```

**Linting:**
```bash
# Using golangci-lint (recommended)
golangci-lint run

# Using go vet
go vet ./...

# Format code
go fmt ./...
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

### Claude Code Statusline Integration

The tool implements the Claude Code statusline protocol, which expects:
- JSON output written to stdout
- Specific data structure for statusline information
- Real-time updates based on Claude Code session state

### Key Components

**Charm Libraries Usage:**
- **Bubble Tea**: Used for building the TUI framework and managing application state
- **Gum**: Inspiration for CLI interaction patterns and user-friendly output styling
- **Lipgloss** (likely): For styling and formatting terminal output

### Statusline Data Structure

The tool should output JSON that conforms to the Claude Code statusline API specification. Refer to the official documentation for the exact schema: https://code.claude.com/docs/en/statusline

### Design Principles

Following the Charm ecosystem style:
- Elegant, minimal terminal UI
- Smooth animations and transitions
- Composable components
- Clear separation between model, view, and update logic (Bubble Tea pattern)

## Claude Code Statusline API

Key information to display may include:
- Current model being used
- Token usage statistics
- Session state
- Active tools or capabilities
- Cost tracking
- Other relevant session metrics
- Keep track of tasks
- And all features that https://github.com/jarrodwatts/claude-hud ship.


Refer to the official docs for the complete API: https://code.claude.com/docs/en/statusline
