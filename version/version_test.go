package version

import (
	"testing"
)

func TestGet(t *testing.T) {
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
		if got == "" {
			t.Error("Get() should not return empty string")
		}
	})

	t.Run("dev version falls back to git or dev", func(t *testing.T) {
		Version = "dev"
		got := Get()
		// "dev" is treated same as empty — tries git first
		if got == "" {
			t.Error("Get() should not return empty string")
		}
	})
}

func TestGetGitVersion(t *testing.T) {
	version := getGitVersion()
	// In a git repo, this should return something
	// If not in a git repo, it returns empty string
	// We just verify it doesn't panic
	t.Logf("Git version: %s", version)
}
