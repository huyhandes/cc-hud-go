package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/internal/git"
	"github.com/huybui/cc-hud-go/output"
	"github.com/huybui/cc-hud-go/parser"
	"github.com/huybui/cc-hud-go/state"
	"github.com/huybui/cc-hud-go/style"
	"github.com/huybui/cc-hud-go/theme"
	"github.com/huybui/cc-hud-go/version"
)

func printUsage() {
	fmt.Fprintf(os.Stderr, `cc-hud-go - Claude Code statusline tool

USAGE:
    cc-hud-go [OPTIONS]

DESCRIPTION:
    A Go-based statusline tool for Claude Code that displays rich, real-time
    information about your current Claude Code session.

    Reads session data from stdin (provided by Claude Code) and outputs
    formatted JSON to stdout for the Claude Code statusline.

OPTIONS:
    -h, --help     Show this help message and exit
    -v, --version  Print version information and exit

CONFIGURATION:
    Config file: ~/.claude/cc-hud-go/config.json

    Available presets:
        full       - All features enabled (default)
        essential  - Core metrics only
        minimal    - Minimal information

    See https://github.com/huyhandes/cc-hud-go#configuration for full
    configuration options.

INTEGRATION:
    Add to your Claude Code config (~/.claude/config.json):

        {
          "statusline": {
            "command": "cc-hud-go"
          }
        }

EXAMPLES:
    # Test with sample data
    echo '{"model":"claude-sonnet-4.5"}' | cc-hud-go

    # Check version
    cc-hud-go --version

    # Show help
    cc-hud-go --help

MORE INFO:
    Documentation: https://github.com/huyhandes/cc-hud-go
    Report issues: https://github.com/huyhandes/cc-hud-go/issues

`)
}

func main() {
	// Customize usage message
	flag.Usage = printUsage

	// Define flags
	var (
		versionFlag bool
		helpFlag    bool
	)

	flag.BoolVar(&versionFlag, "version", false, "Print version and exit")
	flag.BoolVar(&versionFlag, "v", false, "Print version and exit (shorthand)")
	flag.BoolVar(&helpFlag, "help", false, "Show help message and exit")
	flag.BoolVar(&helpFlag, "h", false, "Show help message and exit (shorthand)")

	// Parse flags
	flag.Parse()

	// Handle help flag
	if helpFlag {
		printUsage()
		os.Exit(0)
	}

	// Handle version flag
	if versionFlag {
		fmt.Println(version.Get())
		os.Exit(0)
	}
	// Load config
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".claude", "cc-hud-go", "config.json")
	cfg, err := config.LoadFromFile(configPath)
	if err != nil {
		cfg = config.Default()
	}

	// Initialize theme and style system
	themeInstance := theme.LoadThemeFromConfig(cfg.Theme, cfg.Colors)
	style.Init(themeInstance)

	// Initialize state
	s := state.New()

	// Read JSON from stdin (Claude Code sends one JSON object)
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		os.Exit(1)
	}

	// Parse stdin JSON and update state
	if err := parser.ParseStdin(data, s); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing stdin: %v\n", err)
		os.Exit(1)
	}

	// Parse transcript file for tool usage if available
	if s.Session.TranscriptPath != "" {
		if err := parser.ParseTranscript(s.Session.TranscriptPath, s); err != nil {
			// Don't fail on transcript errors, just log
			fmt.Fprintf(os.Stderr, "Warning: failed to parse transcript: %v\n", err)
		}
	}

	// Update git information
	if branch, err := git.GetBranch(); err == nil {
		s.Git.Branch = branch
	}

	if status, err := git.GetStatus(); err == nil {
		s.Git.DirtyFiles = status.DirtyFiles
		s.Git.Ahead = status.Ahead
		s.Git.Behind = status.Behind
		s.Git.Added = status.Added
		s.Git.Modified = status.Modified
		s.Git.Deleted = status.Deleted
	}

	// Render and output statusline
	result, err := output.Render(s, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering output: %v\n", err)
		os.Exit(1)
	}

	// Output to stdout and exit
	fmt.Println(result)
}
