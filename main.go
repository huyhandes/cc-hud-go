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
)

// version is set at build time via -ldflags
var version = "dev"

func main() {
	// Parse flags
	versionFlag := flag.Bool("version", false, "Print version and exit")
	flag.BoolVar(versionFlag, "v", false, "Print version and exit (shorthand)")
	flag.Parse()

	// Handle version flag
	if *versionFlag {
		fmt.Printf("cc-hud-go %s\n", version)
		os.Exit(0)
	}
	// Load config
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".claude", "cc-hud-go", "config.json")
	cfg, err := config.LoadFromFile(configPath)
	if err != nil {
		cfg = config.Default()
	}

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
