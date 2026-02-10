//go:build integration
// +build integration

package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestIntegration(t *testing.T) {
	// Build binary
	cmd := exec.Command("go", "build", "-o", "cc-hud-go-test", ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build: %v", err)
	}
	defer exec.Command("rm", "cc-hud-go-test").Run()

	// Create a dummy transcript file so the watcher doesn't block
	testHome := "/tmp/cc-hud-go-test-home"
	transcriptDir := testHome + "/.claude"
	exec.Command("mkdir", "-p", transcriptDir).Run()
	exec.Command("touch", transcriptDir+"/transcript.jsonl").Run()
	defer exec.Command("rm", "-rf", testHome).Run()

	// Create a test script that pipes input with correct Claude Code format
	testScript := `/bin/sh -c '
export HOME=` + testHome + `
echo "{\"session_id\":\"test123\",\"cwd\":\"/test\",\"model\":{\"id\":\"claude-sonnet-4-5\",\"display_name\":\"Sonnet 4.5\"},\"workspace\":{\"current_dir\":\"/test\",\"project_dir\":\"/test\"},\"context_window\":{\"total_input_tokens\":50000,\"total_output_tokens\":10000,\"context_window_size\":200000,\"used_percentage\":30.0,\"remaining_percentage\":70.0,\"current_usage\":{\"input_tokens\":40000,\"output_tokens\":10000,\"cache_creation_input_tokens\":5000,\"cache_read_input_tokens\":5000}}}"
' | ./cc-hud-go-test 2>&1`

	// Run the test script
	output, err := exec.Command("/bin/sh", "-c", testScript).Output()
	if err != nil {
		// Check if it's just a kill signal (expected)
		if _, ok := err.(*exec.ExitError); !ok {
			t.Fatalf("failed to run: %v", err)
		}
	}

	outputStr := strings.TrimSpace(string(output))
	if outputStr == "" {
		t.Fatal("no output received")
	}

	// Output should be plain text, not JSON
	// Should contain model name
	if !strings.Contains(outputStr, "Sonnet") {
		t.Errorf("expected output to contain 'Sonnet', got: %s", outputStr)
	}

	// Should contain context bar indicator
	if !strings.Contains(outputStr, "[") || !strings.Contains(outputStr, "]") {
		t.Errorf("expected output to contain progress bar, got: %s", outputStr)
	}

	// Should contain separator
	if !strings.Contains(outputStr, "|") {
		t.Errorf("expected output to contain separator '|', got: %s", outputStr)
	}

	t.Logf("Success! Output: %s", outputStr)
}

func TestIntegrationWithRateLimits(t *testing.T) {
	// Build binary
	cmd := exec.Command("go", "build", "-o", "cc-hud-go-test", ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build: %v", err)
	}
	defer exec.Command("rm", "cc-hud-go-test").Run()

	// Create a dummy transcript file so the watcher doesn't block
	testHome := "/tmp/cc-hud-go-test-home"
	transcriptDir := testHome + "/.claude"
	exec.Command("mkdir", "-p", transcriptDir).Run()
	exec.Command("touch", transcriptDir+"/transcript.jsonl").Run()
	defer exec.Command("rm", "-rf", testHome).Run()

	// Create a test script that pipes input with rate limits
	testScript := `/bin/sh -c '
export HOME=` + testHome + `
echo "{\"session_id\":\"test123\",\"cwd\":\"/test\",\"model\":{\"id\":\"claude-sonnet-4-5\",\"display_name\":\"Sonnet 4.5\"},\"workspace\":{\"current_dir\":\"/test\",\"project_dir\":\"/test\"},\"context_window\":{\"total_input_tokens\":50000,\"total_output_tokens\":10000,\"context_window_size\":200000,\"used_percentage\":30.0,\"remaining_percentage\":70.0},\"rate_limits\":{\"hourly_used\":25,\"hourly_total\":50,\"seven_day_used\":670,\"seven_day_total\":1000}}"
' | ./cc-hud-go-test 2>&1`

	// Run the test script
	output, err := exec.Command("/bin/sh", "-c", testScript).Output()
	if err != nil {
		// Check if it's just a kill signal (expected)
		if _, ok := err.(*exec.ExitError); !ok {
			t.Fatalf("failed to run: %v", err)
		}
	}

	outputStr := strings.TrimSpace(string(output))
	if outputStr == "" {
		t.Fatal("no output received")
	}

	// Should contain rate limit emoji indicator
	if !strings.Contains(outputStr, "📊") {
		t.Errorf("expected output to contain rate limit indicator 📊, got: %s", outputStr)
	}

	// Should contain percentage (67%)
	if !strings.Contains(outputStr, "67%") {
		t.Errorf("expected output to contain '67%%', got: %s", outputStr)
	}

	// Rate limit should be on line 1 (first line)
	lines := strings.Split(outputStr, "\n")
	if len(lines) < 1 {
		t.Fatal("expected at least 1 line of output")
	}
	firstLine := lines[0]
	if !strings.Contains(firstLine, "📊") {
		t.Errorf("expected rate limit on first line, got: %s", firstLine)
	}

	t.Logf("Success with rate limits! Output: %s", outputStr)
}
