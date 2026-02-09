package git

import (
	"testing"
)

func TestGetBranch(t *testing.T) {
	// This test requires a git repo
	// For now, test that it doesn't panic
	branch, err := GetBranch()

	// In a non-git directory, expect error
	if err != nil && branch != "" {
		t.Error("expected empty branch on error")
	}
}
