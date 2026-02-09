package version

import (
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	// Save original version
	originalVersion := Version
	defer func() { Version = originalVersion }()

	t.Run("returns set version", func(t *testing.T) {
		Version = "v1.2.3"
		got := Get()
		if got != "v1.2.3" {
			t.Errorf("Get() = %q, want %q", got, "v1.2.3")
		}
	})

	t.Run("falls back to git or dev", func(t *testing.T) {
		Version = ""
		got := Get()
		// Should return either git version or "dev"
		if got == "" {
			t.Error("Get() should not return empty string")
		}
	})
}

func TestGetDetailed(t *testing.T) {
	// Save original version
	originalVersion := Version
	defer func() { Version = originalVersion }()

	Version = "v1.0.0"
	got := GetDetailed()

	// Should contain version
	if !strings.Contains(got, "cc-hud-go") {
		t.Errorf("GetDetailed() should contain 'cc-hud-go', got: %s", got)
	}

	if !strings.Contains(got, "v1.0.0") {
		t.Errorf("GetDetailed() should contain version 'v1.0.0', got: %s", got)
	}
}

func TestGetGitVersion(t *testing.T) {
	version := getGitVersion()
	// In a git repo, this should return something
	// If not in a git repo, it returns empty string
	// We just verify it doesn't panic
	t.Logf("Git version: %s", version)
}

func TestGetGitCommit(t *testing.T) {
	commit := getGitCommit()
	// Should return short hash or empty if not in git repo
	t.Logf("Git commit: %s", commit)

	if commit != "" && len(commit) < 7 {
		t.Errorf("Git commit hash seems too short: %s", commit)
	}
}

func TestIsGitDirty(t *testing.T) {
	// Just verify it doesn't panic
	isDirty := isGitDirty()
	t.Logf("Git dirty: %v", isDirty)
}
