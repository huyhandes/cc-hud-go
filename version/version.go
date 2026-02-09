package version

import (
	"fmt"
	"os/exec"
	"strings"
)

// Version is set at build time via -ldflags "-X github.com/huybui/cc-hud-go/version.Version=x.y.z"
var Version = ""

// Get returns the version string
// If Version is not set at build time, it attempts to get the version from git
func Get() string {
	if Version != "" && Version != "dev" {
		return Version
	}

	// Try to get version from git
	if gitVersion := getGitVersion(); gitVersion != "" {
		return gitVersion
	}

	// Fallback to dev
	return "dev"
}

// getGitVersion attempts to get version from git describe
func getGitVersion() string {
	cmd := exec.Command("git", "describe", "--tags", "--always", "--dirty")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	version := strings.TrimSpace(string(output))
	if version == "" {
		return ""
	}

	return version
}

// GetDetailed returns detailed version information
func GetDetailed() string {
	v := Get()
	gitCommit := getGitCommit()
	gitDirty := isGitDirty()

	info := fmt.Sprintf("cc-hud-go %s", v)

	if gitCommit != "" {
		info += fmt.Sprintf(" (commit: %s", gitCommit)
		if gitDirty {
			info += ", dirty"
		}
		info += ")"
	}

	return info
}

// getGitCommit returns the current git commit hash
func getGitCommit() string {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// isGitDirty checks if the working directory has uncommitted changes
func isGitDirty() bool {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(strings.TrimSpace(string(output))) > 0
}
