package version

import (
	"os/exec"
	"strings"
)

// Version is set at build time via -ldflags "-X github.com/huyhandes/cc-hud-go/version.Version=x.y.z"
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
