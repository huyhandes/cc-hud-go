package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/internal/git"
	"github.com/huybui/cc-hud-go/internal/watcher"
	"github.com/huybui/cc-hud-go/output"
	"github.com/huybui/cc-hud-go/parser"
	"github.com/huybui/cc-hud-go/segment"
	"github.com/huybui/cc-hud-go/state"
)

type model struct {
	state    *state.State
	config   *config.Config
	segments []segment.Segment
}

type stdinMsg struct {
	data []byte
	err  error
}

type transcriptMsg struct {
	line string
}

type tickMsg time.Time

func initialModel() model {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".claude", "cc-hud-go", "config.json")

	cfg, err := config.LoadFromFile(configPath)
	if err != nil {
		cfg = config.Default()
	}

	return model{
		state:    state.New(),
		config:   cfg,
		segments: segment.All(),
	}
}

func readStdinCmd() tea.Cmd {
	return func() tea.Msg {
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadBytes('\n')
		return stdinMsg{data: line, err: err}
	}
}

func watchTranscriptCmd() tea.Cmd {
	return func() tea.Msg {
		// Find transcript file (simplified - in production would query Claude Code)
		home, _ := os.UserHomeDir()
		transcriptPath := filepath.Join(home, ".claude", "transcript.jsonl")

		lines := make(chan string, 10)
		stop := make(chan struct{})

		go func() {
			watcher.Watch(transcriptPath, lines, stop)
		}()

		// Read one line
		line := <-lines
		return transcriptMsg{line: line}
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func outputCmd(m model) tea.Cmd {
	return func() tea.Msg {
		json, err := output.Render(m.state, m.config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "render error: %v\n", err)
			return nil
		}

		fmt.Println(json)
		return nil
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(readStdinCmd(), watchTranscriptCmd(), tickCmd())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case stdinMsg:
		if msg.err != nil {
			if msg.err == io.EOF {
				return m, tea.Quit
			}
			// Log error but continue
			fmt.Fprintf(os.Stderr, "stdin error: %v\n", msg.err)
			return m, readStdinCmd()
		}

		// Parse and update state
		if err := parser.ParseStdin(msg.data, m.state); err != nil {
			fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		}

		// Output JSON and read next line
		return m, tea.Batch(outputCmd(m), readStdinCmd())

	case transcriptMsg:
		if err := parser.ParseTranscriptLine([]byte(msg.line), m.state); err != nil {
			fmt.Fprintf(os.Stderr, "transcript parse error: %v\n", err)
		}
		return m, tea.Batch(outputCmd(m), watchTranscriptCmd())

	case tickMsg:
		// Update git info
		if branch, err := git.GetBranch(); err == nil {
			m.state.Git.Branch = branch
		}

		if status, err := git.GetStatus(); err == nil {
			m.state.Git.DirtyFiles = status.DirtyFiles
			m.state.Git.Ahead = status.Ahead
			m.state.Git.Behind = status.Behind
			m.state.Git.Added = status.Added
			m.state.Git.Modified = status.Modified
			m.state.Git.Deleted = status.Deleted
		}

		return m, tea.Batch(outputCmd(m), tickCmd())
	}

	return m, nil
}

func (m model) View() string {
	return ""
}

func main() {
	// Run without alternate screen and without mouse since we're a statusline tool
	p := tea.NewProgram(
		initialModel(),
		tea.WithInput(nil),
		tea.WithoutRenderer(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
