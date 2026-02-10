package style

import (
	"strings"
	"testing"
)

func TestRenderTable(t *testing.T) {
	headers := []string{"Column1", "Column2"}
	rows := [][]string{
		{"Value1", "Value2"},
		{"Value3", "Value4"},
	}

	result := RenderTable(headers, rows)

	// Should contain table border characters
	if !strings.Contains(result, "┌") || !strings.Contains(result, "┐") {
		t.Error("Expected table top border characters")
	}

	if !strings.Contains(result, "└") || !strings.Contains(result, "┘") {
		t.Error("Expected table bottom border characters")
	}

	// Should contain headers
	if !strings.Contains(result, "Column1") || !strings.Contains(result, "Column2") {
		t.Error("Expected table headers in output")
	}

	// Should contain data
	if !strings.Contains(result, "Value1") || !strings.Contains(result, "Value4") {
		t.Error("Expected table data in output")
	}
}
