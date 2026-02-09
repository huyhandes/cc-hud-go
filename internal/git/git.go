package git

import (
	"bytes"
	"context"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// GetBranch returns the current git branch name
func GetBranch() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}

// Status holds git status information
type Status struct {
	DirtyFiles int
	Ahead      int
	Behind     int
	Added      int
	Modified   int
	Deleted    int
}

// GetStatus returns git status information
func GetStatus() (*Status, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	status := &Status{}

	// Get ahead/behind
	cmd := exec.CommandContext(ctx, "git", "rev-list", "--left-right", "--count", "HEAD...@{upstream}")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err == nil {
		parts := strings.Fields(out.String())
		if len(parts) == 2 {
			status.Ahead, _ = strconv.Atoi(parts[0])
			status.Behind, _ = strconv.Atoi(parts[1])
		}
	}

	// Get file stats
	cmd = exec.CommandContext(ctx, "git", "status", "--porcelain")
	out.Reset()
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if len(line) < 3 {
			continue
		}

		status.DirtyFiles++

		code := line[:2]
		if strings.Contains(code, "A") {
			status.Added++
		}
		if strings.Contains(code, "M") {
			status.Modified++
		}
		if strings.Contains(code, "D") {
			status.Deleted++
		}
	}

	return status, nil
}
