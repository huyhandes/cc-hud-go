package config

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg := Default()

	if !cfg.Display.Model {
		t.Error("expected Display.Model to be true")
	}
}

func TestLoadFromFile(t *testing.T) {
	// Test missing file (should return defaults)
	cfg, err := LoadFromFile("nonexistent.json")
	if err != nil {
		t.Fatalf("should not error on missing file: %v", err)
	}

	if cfg.Theme != "macchiato" {
		t.Error("expected default theme on missing file")
	}

	// Test invalid JSON (should return defaults)
	cfg, err = LoadFromFile("../testdata/config_invalid.json")
	if err != nil {
		t.Fatalf("should not error on invalid JSON: %v", err)
	}

	if cfg.Theme != "macchiato" {
		t.Error("expected default theme on invalid JSON")
	}
}

func TestConfigTheme(t *testing.T) {
	cfg := Default()

	if cfg.Theme == "" {
		t.Error("Expected default theme to be set")
	}

	if cfg.Theme != "macchiato" {
		t.Errorf("Expected default theme 'macchiato', got %s", cfg.Theme)
	}
}

func TestConfigColorOverrides(t *testing.T) {
	cfg := Default()

	if cfg.Colors == nil {
		t.Error("Expected Colors map to be initialized")
	}
}
