//go:build integration

package main

import (
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

func buildBinary(t *testing.T) (string, func()) {
	t.Helper()
	if err := exec.Command("go", "build", "-o", "cc-hud-go-test", ".").Run(); err != nil {
		t.Fatalf("failed to build: %v", err)
	}
	return "cc-hud-go-test", func() { exec.Command("rm", "cc-hud-go-test").Run() }
}

// runBinary pipes stdin into the built binary, optionally appending extra
// argv to the invocation. CLICOLOR_FORCE=1 makes lipgloss emit ANSI codes
// even though stdout is a pipe (otherwise theme differences are invisible
// to the theme-flag integration test).
func runBinary(t *testing.T, stdin string, args ...string) string {
	t.Helper()
	argv := []string{"-c", "CLICOLOR_FORCE=1 /bin/sh -c 'echo " + strconv.Quote(stdin) + "' | CLICOLOR_FORCE=1 ./cc-hud-go-test " +
		strings.Join(args, " ") + " 2>&1"}
	out, err := exec.Command("/bin/sh", argv...).Output()
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			t.Fatalf("failed to run: %v", err)
		}
	}
	return strings.TrimSpace(string(out))
}

func TestIntegration(t *testing.T) {
	_, cleanup := buildBinary(t)
	defer cleanup()

	stdin := `{"session_id":"test123","cwd":"/test","model":{"id":"claude-sonnet-4-5","display_name":"Sonnet 4.5"},"workspace":{"current_dir":"/test","project_dir":"/test"},"context_window":{"total_input_tokens":50000,"total_output_tokens":10000,"context_window_size":200000,"used_percentage":30.0,"remaining_percentage":70.0,"current_usage":{"input_tokens":40000,"output_tokens":10000,"cache_creation_input_tokens":5000,"cache_read_input_tokens":5000}}}`

	outputStr := runBinary(t, stdin)
	if outputStr == "" {
		t.Fatal("no output received")
	}

	// Should contain model name and context bar (stdin-driven, still live).
	if !strings.Contains(outputStr, "Sonnet") {
		t.Errorf("expected output to contain 'Sonnet', got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "🧠") {
		t.Errorf("expected output to contain context bar, got: %s", outputStr)
	}

	// Primary external-behavior seam: deleted segments must not render.
	for _, glyph := range []string{"📦", "👤", "⏳", "🔄", "✅"} {
		if strings.Contains(outputStr, glyph) {
			t.Errorf("output should NOT contain deleted-segment glyph %s, got: %s", glyph, outputStr)
		}
	}

	t.Logf("Success! Output: %s", outputStr)
}

func TestIntegrationWithRateLimits(t *testing.T) {
	_, cleanup := buildBinary(t)
	defer cleanup()

	// OAuth is now the single source of truth for rate limits; piped stdin
	// no longer drives the 📊 segment. Focus on what stdin still controls:
	// model and context.
	stdin := `{"session_id":"test123","cwd":"/test","model":{"id":"claude-sonnet-4-5","display_name":"Sonnet 4.5"},"workspace":{"current_dir":"/test","project_dir":"/test"},"context_window":{"total_input_tokens":50000,"total_output_tokens":10000,"context_window_size":200000,"used_percentage":30.0,"remaining_percentage":70.0}}`

	outputStr := runBinary(t, stdin)
	if outputStr == "" {
		t.Fatal("no output received")
	}

	if !strings.Contains(outputStr, "Sonnet") {
		t.Errorf("expected output to contain 'Sonnet', got: %s", outputStr)
	}

	// Rate-limit segments render empty without OAuth data — verify no
	// half-rendered rate-limit line leaks through.
	for _, glyph := range []string{"⏳", "🔄", "✅", "📦", "👤"} {
		if strings.Contains(outputStr, glyph) {
			t.Errorf("output should NOT contain deleted/empty-segment glyph %s, got: %s", glyph, outputStr)
		}
	}

	t.Logf("Success with rate limits path! Output: %s", outputStr)
}

// TestIntegrationThemeFlag proves --theme wires through end-to-end: mocha
// and the default (macchiato) must render different ANSI color codes.
// Also confirms an unknown value falls back gracefully (no error, same
// output as default).
//
// Comparison ignores line 1 because the rate-limit segments pull live
// OAuth data that varies call-to-call; lines 2+ are colored and stable.
func TestIntegrationThemeFlag(t *testing.T) {
	_, cleanup := buildBinary(t)
	defer cleanup()

	stdin := `{"session_id":"test123","cwd":"/test","model":{"id":"claude-sonnet-4-5","display_name":"Sonnet 4.5"},"workspace":{"current_dir":"/test","project_dir":"/test"},"context_window":{"total_input_tokens":50000,"total_output_tokens":10000,"context_window_size":200000,"used_percentage":30.0,"remaining_percentage":70.0,"current_usage":{"input_tokens":40000,"output_tokens":10000,"cache_creation_input_tokens":5000,"cache_read_input_tokens":5000}}}`

	// dropLine1 removes the model+context+ratelimit line (volatile OAuth data).
	dropLine1 := func(out string) string {
		lines := strings.Split(out, "\n")
		if len(lines) <= 1 {
			return out
		}
		return strings.Join(lines[1:], "\n")
	}

	defaultOut := dropLine1(runBinary(t, stdin))
	mochaOut := dropLine1(runBinary(t, stdin, "--theme", "mocha"))
	bogusOut := dropLine1(runBinary(t, stdin, "--theme", "bogus"))

	if defaultOut == "" || mochaOut == "" {
		t.Fatal("expected non-empty output for both default and mocha")
	}

	// mocha must differ in styling from default macchiato.
	if defaultOut == mochaOut {
		t.Errorf("expected --theme mocha to produce different ANSI codes than default macchiato\n"+
			"default: %q\nmocha:   %q", defaultOut, mochaOut)
	}

	// Unknown value falls back to macchiato — output must match default.
	if bogusOut != defaultOut {
		t.Errorf("expected --theme bogus to fall back to macchiato (same as default)\n"+
			"default: %q\nbogus:   %q", defaultOut, bogusOut)
	}

	t.Logf("default: %s", defaultOut)
	t.Logf("mocha:   %s", mochaOut)
}
