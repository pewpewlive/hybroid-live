package evaluator

import (
	"hybroid/core"
	"strings"
	"testing"
)

func TestParserAlertsPersistence(t *testing.T) {
	// Setup a minimal evaluator
	files := []core.FileInformation{
		{
			DirectoryPath: ".",
			FileName:      "test",
			FileExtension: ".hyb",
		},
	}
	eval := NewEvaluator(files)

	code := `
env L as Level
use Pewpew
Pewpew:Print(
` // Missing ) -> Parser error 'expected )'

	// Mock the update flow from LSP
	// 1. Update content (Runs parser, generates alerts)
	eval.UpdateFileContent("test.hyb", code)

	// Check alerts immediately after update
	alertsBefore := eval.GetAlerts("test.hyb")
	if len(alertsBefore) == 0 {
		t.Fatalf("Expected parser alerts after UpdateFileContent, got 0")
	}
	t.Logf("Alerts before RunAnalysis: %d", len(alertsBefore))

	// 2. Run Analysis (This is where we suspect alerts are wiped)
	eval.RunAnalysis()

	// Check alerts after analysis
	alertsAfter := eval.GetAlerts("test.hyb")
	t.Logf("Alerts after RunAnalysis: %d", len(alertsAfter))
	foundParserError := false
	for _, a := range alertsAfter {
		t.Logf("Alert: %s", a.Message())
		if strings.Contains(strings.ToLower(a.Message()), "expected") {
			foundParserError = true
		}
	}

	if !foundParserError {
		t.Fatalf("Parser diagnostic was wiped out by RunAnalysis!")
	}
}
