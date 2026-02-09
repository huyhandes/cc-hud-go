//go:build integration
// +build integration

package main

import (
	"encoding/json"
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

	// Create a test script that pipes input
	testScript := `/bin/sh -c '
export HOME=` + testHome + `
echo "{\"model\":\"claude-sonnet-4.5\",\"planType\":\"Pro\",\"context\":{\"used\":50000,\"total\":200000}}"
sleep 1
' | ./cc-hud-go-test 2>&1`

	// Run the test script
	output, err := exec.Command("/bin/sh", "-c", testScript).Output()
	if err != nil {
		// Check if it's just a kill signal (expected)
		if _, ok := err.(*exec.ExitError); !ok {
			t.Fatalf("failed to run: %v", err)
		}
	}

	outputStr := string(output)
	if outputStr == "" {
		t.Fatal("no output received")
	}

	// Find the first line that looks like JSON
	lines := strings.Split(outputStr, "\n")
	var jsonLine string
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "{") {
			jsonLine = line
			break
		}
	}

	if jsonLine == "" {
		t.Fatalf("no JSON found in output: %s", outputStr)
	}

	// Parse and validate JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonLine), &result); err != nil {
		t.Fatalf("failed to parse JSON: %v\nLine: %s", err, jsonLine)
	}

	segments, ok := result["segments"].([]interface{})
	if !ok || len(segments) == 0 {
		t.Error("expected segments in output")
	}

	t.Logf("Success! Output: %s", jsonLine)
}
