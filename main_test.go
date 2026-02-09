package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestVersionFlag(t *testing.T) {
	tests := []struct {
		name string
		flag string
	}{
		{"long form", "--version"},
		{"short form", "-v"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("go", "run", ".", tt.flag)
			output, err := cmd.Output()
			if err != nil {
				t.Fatalf("Command failed: %v", err)
			}

			got := strings.TrimSpace(string(output))
			// Version should be non-empty and either be "dev", start with "v", or be a git hash
			if got == "" {
				t.Error("Version output should not be empty")
			}

			// Should be a valid version format (either "dev", "vX.Y.Z", or git describe format)
			if got != "dev" && !strings.HasPrefix(got, "v") && len(got) < 7 {
				t.Errorf("Version output has unexpected format: %s", got)
			}
		})
	}
}

func TestHelpFlag(t *testing.T) {
	tests := []struct {
		name string
		flag string
	}{
		{"long form", "--help"},
		{"short form", "-h"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("go", "run", ".", tt.flag)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			// Help flag exits with 0, so we don't check error
			_ = cmd.Run()

			// Help message goes to stderr
			output := stderr.String()

			// Check for key sections in help text
			expectedSections := []string{
				"USAGE:",
				"DESCRIPTION:",
				"OPTIONS:",
				"CONFIGURATION:",
				"INTEGRATION:",
				"EXAMPLES:",
				"MORE INFO:",
			}

			for _, section := range expectedSections {
				if !strings.Contains(output, section) {
					t.Errorf("Help output missing section '%s'", section)
				}
			}

			// Check that both flags are documented
			if !strings.Contains(output, "-h, --help") {
				t.Error("Help output should document -h and --help flags")
			}
			if !strings.Contains(output, "-v, --version") {
				t.Error("Help output should document -v and --version flags")
			}
		})
	}
}

func TestInvalidFlag(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "--invalid-flag")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Invalid flag should exit with non-zero
	err := cmd.Run()
	if err == nil {
		t.Error("Invalid flag should cause non-zero exit")
	}

	output := stderr.String()

	// Should show error message
	if !strings.Contains(output, "flag provided but not defined") {
		t.Error("Should show error for invalid flag")
	}

	// Should still show usage information
	if !strings.Contains(output, "USAGE:") {
		t.Error("Should show usage on invalid flag")
	}
}

func TestUsageFunction(t *testing.T) {
	// Test that printUsage doesn't panic and contains expected content
	// We can't easily test the output directly since it writes to stderr,
	// but we can at least verify the function exists and runs
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("printUsage() panicked: %v", r)
		}
	}()

	// The function is tested indirectly via the help flag tests above
	// This test just ensures no panic occurs
	printUsage()
}
