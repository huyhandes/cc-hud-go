package segment

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// ansiRe matches CSI escape sequences emitted by lipgloss (color/bold codes).
// We strip them so substring assertions don't depend on terminal styling.
var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// stripANSI removes lipgloss's ANSI styling from s.
func stripANSI(s string) string { return ansiRe.ReplaceAllString(s, "") }

// writeFlag writes contents to <dir>/<flagFile>, failing the test on error.
func writeFlag(t *testing.T, dir, flagFile, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(filepath.Join(dir, flagFile)), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, flagFile), []byte(contents), 0o644); err != nil {
		t.Fatalf("write flag: %v", err)
	}
}

// withConfigDir points CLAUDE_CONFIG_DIR at a temp dir for the test.
func withConfigDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("CLAUDE_CONFIG_DIR", dir)
	return dir
}

func TestRenderModeBadge(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T)
		want  string
		exact bool // true → want==got; false → strings.Contains
	}{
		{
			name:  "absent flag file",
			setup: func(t *testing.T) { withConfigDir(t) },
			want:  "",
			exact: true,
		},
		{
			name: "lite",
			setup: func(t *testing.T) {
				writeFlag(t, withConfigDir(t), ".caveman-active", "lite")
			},
			want:  "🦴 CAVEMAN:LITE",
			exact: false,
		},
		{
			name: "full",
			setup: func(t *testing.T) {
				writeFlag(t, withConfigDir(t), ".caveman-active", "full")
			},
			want:  "🦴 CAVEMAN:FULL",
			exact: false,
		},
		{
			name: "ultra",
			setup: func(t *testing.T) {
				writeFlag(t, withConfigDir(t), ".caveman-active", "ultra")
			},
			want:  "🦴 CAVEMAN:ULTRA",
			exact: false,
		},
		{
			name: "whitespace and uppercase normalized",
			setup: func(t *testing.T) {
				writeFlag(t, withConfigDir(t), ".caveman-active", "  ULTRA\n")
			},
			want:  "🦴 CAVEMAN:ULTRA",
			exact: false,
		},
		{
			name: "garbage contents",
			setup: func(t *testing.T) {
				writeFlag(t, withConfigDir(t), ".caveman-active", "FOOBAR")
			},
			want:  "",
			exact: true,
		},
		{
			name: "empty contents",
			setup: func(t *testing.T) {
				writeFlag(t, withConfigDir(t), ".caveman-active", "")
			},
			want:  "",
			exact: true,
		},
		{
			name: "flag path is a directory (read error)",
			setup: func(t *testing.T) {
				dir := withConfigDir(t)
				// Directory where the flag file would be → ReadFile fails.
				if err := os.MkdirAll(filepath.Join(dir, ".caveman-active"), 0o755); err != nil {
					t.Fatalf("mkdir flag: %v", err)
				}
			},
			want:  "",
			exact: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)
			got := stripANSI(renderModeBadge(".caveman-active", "🦴", "CAVEMAN"))
			switch {
			case tt.exact && got != tt.want:
				t.Errorf("renderModeBadge = %q, want %q", got, tt.want)
			case !tt.exact && !strings.Contains(got, tt.want):
				t.Errorf("renderModeBadge = %q, want it to contain %q", got, tt.want)
			}
		})
	}
}

// TestRenderModeBadgeConfigDirOverride confirms the override is honored AND
// that when the override dir has no flag, no badge is rendered (override must
// win over $HOME/.claude). This matches ponytail's own statusline script.
func TestRenderModeBadgeConfigDirOverride(t *testing.T) {
	overrideDir := t.TempDir()
	t.Setenv("CLAUDE_CONFIG_DIR", overrideDir)

	// Flag under override dir → badge rendered.
	writeFlag(t, overrideDir, ".caveman-active", "ultra")
	got := stripANSI(renderModeBadge(".caveman-active", "🦴", "CAVEMAN"))
	if !strings.Contains(got, "🦴 CAVEMAN:ULTRA") {
		t.Errorf("override dir with flag: got %q, want substring %q", got, "🦴 CAVEMAN:ULTRA")
	}

	// Fresh override dir with no flag → empty, regardless of $HOME/.claude.
	overrideDir2 := t.TempDir()
	t.Setenv("CLAUDE_CONFIG_DIR", overrideDir2)
	if got := renderModeBadge(".caveman-active", "🦴", "CAVEMAN"); got != "" {
		t.Errorf("override dir without flag: got %q, want empty (override must win over $HOME/.claude)", got)
	}
}

// TestClaudeConfigDirUnsetFallback verifies that with CLAUDE_CONFIG_DIR unset,
// claudeConfigDir() returns $HOME/.claude exactly (default behavior preserved).
// Setting the env var to "" is equivalent to unsetting it: claudeConfigDir
// treats empty the same as absent (dir != "" check), and t.Setenv auto-restores.
func TestClaudeConfigDirUnsetFallback(t *testing.T) {
	t.Setenv("CLAUDE_CONFIG_DIR", "")

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("UserHomeDir: %v", err)
	}
	if got := claudeConfigDir(); got != filepath.Join(home, ".claude") {
		t.Errorf("claudeConfigDir (unset env) = %q, want %q", got, filepath.Join(home, ".claude"))
	}
}

func TestCavemanSegmentDelegates(t *testing.T) {
	// Smoke test: CavemanSegment.Render is a one-line delegation to
	// renderModeBadge. No flag file → empty, no error.
	withConfigDir(t)
	seg := &CavemanSegment{}
	out, err := seg.Render(nil, nil)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}
	if out != "" {
		t.Errorf("Render with no flag = %q, want empty", out)
	}
}
