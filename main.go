package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/huyhandes/cc-hud-go/internal/git"
	"github.com/huyhandes/cc-hud-go/internal/oauth"
	"github.com/huyhandes/cc-hud-go/output"
	"github.com/huyhandes/cc-hud-go/parser"
	"github.com/huyhandes/cc-hud-go/state"
	"github.com/huyhandes/cc-hud-go/style"
	"github.com/huyhandes/cc-hud-go/theme"
	"github.com/huyhandes/cc-hud-go/version"
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
    --theme NAME   Catppuccin variant to render (default macchiato).
                   Values: macchiato, mocha, frappe, latte.
                   Unknown values fall back to macchiato.

CONFIGURATION:
    The tool ships opinionated defaults; there is no config file. The only
    knob is --theme. Select a variant by passing the flag to your statusline
    command.

INTEGRATION:
    Add to your Claude Code config (~/.claude/config.json):

        {
          "statusline": {
            "command": "cc-hud-go --theme mocha"
          }
        }

EXAMPLES:
    # Test with sample data
    echo '{"model":"claude-sonnet-4.5"}' | cc-hud-go

    # Select a theme variant
    echo '{...}' | cc-hud-go --theme latte

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
		themeFlag   string
	)

	flag.BoolVar(&versionFlag, "version", false, "Print version and exit")
	flag.BoolVar(&versionFlag, "v", false, "Print version and exit (shorthand)")
	flag.BoolVar(&helpFlag, "help", false, "Show help message and exit")
	flag.BoolVar(&helpFlag, "h", false, "Show help message and exit (shorthand)")
	flag.StringVar(&themeFlag, "theme", "macchiato", "Catppuccin variant (macchiato, mocha, frappe, latte)")

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

	// Initialize theme and style system. Unknown values fall back to
	// macchiato via theme.GetTheme's default branch.
	style.Init(theme.GetTheme(themeFlag))

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

	// Fetch rate limit usage from OAuth API. OAuth is the single source of
	// truth; on failure the rate-limit segments render empty (the
	// FiveHourSegment empty-return pattern).
	if usage, err := oauth.FetchUsage(); err == nil {
		s.RateLimits.FiveHourPercent = usage.FiveHour.Utilization
		s.RateLimits.SevenDayPercent = usage.SevenDay.Utilization
		s.RateLimits.FiveHourResetsAt = usage.FiveHour.ResetsAt.Format("2006-01-02T15:04:05Z07:00")
		s.RateLimits.SevenDayResetsAt = usage.SevenDay.ResetsAt.Format("2006-01-02T15:04:05Z07:00")
	}

	// Render and output statusline
	result, err := output.Render(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering output: %v\n", err)
		os.Exit(1)
	}

	// Output to stdout and exit
	fmt.Println(result)
}
