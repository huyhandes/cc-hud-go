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

	if !strings.Contains(result, "┌") || !strings.Contains(result, "┐") {
		t.Error("Expected table top border characters")
	}

	if !strings.Contains(result, "└") || !strings.Contains(result, "┘") {
		t.Error("Expected table bottom border characters")
	}

	if !strings.Contains(result, "Column1") || !strings.Contains(result, "Column2") {
		t.Error("Expected table headers in output")
	}

	if !strings.Contains(result, "Value1") || !strings.Contains(result, "Value4") {
		t.Error("Expected table data in output")
	}
}

func TestRenderTableSingleColumn(t *testing.T) {
	headers := []string{"Name"}
	rows := [][]string{{"Alice"}, {"Bob"}}

	result := RenderTable(headers, rows)

	if !strings.Contains(result, "Name") {
		t.Error("expected header")
	}
	if !strings.Contains(result, "Alice") || !strings.Contains(result, "Bob") {
		t.Error("expected row data")
	}
}

func TestRenderTableEmptyRows(t *testing.T) {
	headers := []string{"H1", "H2"}
	rows := [][]string{}

	result := RenderTable(headers, rows)

	if !strings.Contains(result, "H1") {
		t.Error("expected headers even with no rows")
	}
}

func TestRenderTableManyColumns(t *testing.T) {
	headers := []string{"A", "B", "C", "D"}
	rows := [][]string{{"1", "2", "3", "4"}}

	result := RenderTable(headers, rows)

	for _, h := range headers {
		if !strings.Contains(result, h) {
			t.Errorf("expected header %q in output", h)
		}
	}
}
