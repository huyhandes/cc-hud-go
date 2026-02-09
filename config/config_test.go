package config

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg := Default()

	if cfg.Preset != "full" {
		t.Errorf("expected preset 'full', got '%s'", cfg.Preset)
	}

	if cfg.PathLevels != 2 {
		t.Errorf("expected pathLevels 2, got %d", cfg.PathLevels)
	}

	if !cfg.Display.Model {
		t.Error("expected Display.Model to be true")
	}
}

func TestEssentialPreset(t *testing.T) {
	cfg := Essential()

	if cfg.Preset != "essential" {
		t.Errorf("expected preset 'essential', got '%s'", cfg.Preset)
	}

	if cfg.LineLayout != "compact" {
		t.Errorf("expected layout 'compact', got '%s'", cfg.LineLayout)
	}

	if cfg.Display.Tools {
		t.Error("expected Display.Tools to be false in essential preset")
	}

	if !cfg.Display.Model || !cfg.Display.Context {
		t.Error("expected core displays (Model, Context) to be true")
	}
}

func TestMinimalPreset(t *testing.T) {
	cfg := Minimal()

	if cfg.Preset != "minimal" {
		t.Errorf("expected preset 'minimal', got '%s'", cfg.Preset)
	}

	if cfg.Display.Git || cfg.Display.Tasks {
		t.Error("expected Git and Tasks to be false in minimal preset")
	}

	if !cfg.Display.Model || !cfg.Display.Context {
		t.Error("expected core displays (Model, Context) to be true")
	}
}
