package parser_test

import (
	"fmt"
	"hybroid/alerts"
	"hybroid/helpers"
	"hybroid/lexer"
	"hybroid/parser"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func printAlerts[A any](t *testing.T, kind string, alertsToPrint ...A) {
	t.Logf("%s %d alert(s):", kind, len(alertsToPrint))
	for i, alert := range alertsToPrint {
		if alert, ok := any(alert).(reflect.Type); ok {
			t.Logf("%d. %s", i+1, alert.Name())
			continue
		}

		actualAlert := any(alert).(alerts.Alert)
		token := actualAlert.GetSpecifier().GetTokens()[0]
		loc := token.Location
		name := reflect.ValueOf(alert).Elem().Type().Name()
		msg := strings.TrimSpace(actualAlert.GetMessage())
		t.Logf("%d. %s (%s) at line %d, column %d on token '%s'", i+1, name, msg, loc.Line, loc.Column.Start, token.Lexeme)
	}
}

type parseResults struct {
	alerts    []alerts.Alert
	hasAlerts bool
}

func performParsing(path, subtest string) (parseResults, error) {
	results := parseResults{}

	files, err := helpers.CollectFiles(path)

	if err != nil {
		return results, err
	}
	if len(files) == 0 {
		return results, fmt.Errorf("found no files in '%s'", path)
	}
	file := slices.IndexFunc(files, func(file helpers.FileInformation) bool {
		return file.FileName == subtest
	})
	if file == -1 {
		return results, fmt.Errorf("missing subtest '%s' file", subtest)
	}

	sourcePath := files[file].Path()
	sourceFile, err := os.OpenFile(filepath.Join(path, sourcePath), os.O_RDONLY, os.ModePerm)
	if err != nil {
		return results, fmt.Errorf("failed to open file '%s': %v", sourcePath, err)
	}

	lexer := lexer.NewLexer(sourceFile)
	tokens, _ := lexer.Tokenize()

	parser := parser.NewParser(tokens)
	parser.Parse()

	results.alerts = parser.GetAlerts()
	results.hasAlerts = len(results.alerts) != 0

	return results, nil
}

func performTest(t *testing.T, testName string, expectedAlerts []reflect.Type) {
	previousDir, _ := os.Getwd()

	path := previousDir + "/tests/" + testName
	os.Chdir(path)

	t.Run("Valid", func(t *testing.T) {
		results, err := performParsing(path, "valid")
		if err != nil {
			t.Error(err)
			return
		}
		if results.hasAlerts {
			t.Errorf("[Valid] Found alerts in *valid* input")
			printAlerts(t, "Unexpected", results.alerts...)
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		results, err := performParsing(path, "invalid")
		if err != nil {
			t.Error(err)
			return
		}
		if !results.hasAlerts {
			t.Errorf("[Invalid] No alerts in *invalid* input")
		}

		alertTypes := make([]reflect.Type, 0)
		for _, alert := range results.alerts {
			alertTypes = append(alertTypes, reflect.ValueOf(alert).Elem().Type())
		}
		if !helpers.ListsAreSame(alertTypes, expectedAlerts) {
			t.Errorf("[Invalid] Mismatch in *expected* and *received* alerts")

			printAlerts(t, "Expected", expectedAlerts...)
			printAlerts(t, "Received", results.alerts...)
		}
	})

	os.Chdir(previousDir)
}

func TestExpressions(t *testing.T) {
	expectedAlerts := []reflect.Type{
		reflect.TypeFor[alerts.ExpectedIdentifier](),
	}
	performTest(t, "expressions", expectedAlerts)
}
