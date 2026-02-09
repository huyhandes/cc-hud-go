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
