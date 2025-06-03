package parser_test

import (
	"fmt"
	"hybroid/alerts"
	"hybroid/core"
	"hybroid/lexer"
	"hybroid/parser"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"
	"time"
)

type parseResults struct {
	alerts    []alerts.Alert
	hasAlerts bool
}

func performParsing(t *testing.T, path, subtest string) (parseResults, error) {
	results := parseResults{}

	files, err := core.CollectFiles(path)

	if err != nil {
		return results, err
	}
	if len(files) == 0 {
		return results, fmt.Errorf("found no files in '%s'", path)
	}
	file := slices.IndexFunc(files, func(file core.FileInformation) bool {
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

	// Two functions that work as callbacks to finishing parsing and checking for timeout
	parseFunc := func(parser *parser.Parser, succeeded chan<- bool) {
		parser.Parse()
		succeeded <- true
	}
	hangCheck := func(succeeded <-chan bool, hangFree chan<- bool) {
		timeout := time.After(2 * time.Second)
		select {
		case <-succeeded:
			hangFree <- true // The parser successfully finished parsing, send true
		case <-timeout:
			hangFree <- false // The parser reached timeout, send false
		}
	}

	// Create a channel to receive successful parsing and to see if it hasn't hung
	succeeded := make(chan bool)
	hangFree := make(chan bool, 1)

	parser := parser.NewParser(tokens)

	// Initiate the parsing and the hang check
	go parseFunc(&parser, succeeded)
	go hangCheck(succeeded, hangFree)

	// Will return true if success, false if timeout
	if !<-hangFree {
		t.Errorf("Parser hung on %s", sourcePath)
		alrts := parser.GetAlerts()
		alerts.PrintAlerts(t, "Hung with", alrts[:min(len(alrts), 100)]...)
		t.FailNow()
	}

	results.alerts = parser.GetAlerts()
	results.hasAlerts = len(results.alerts) != 0

	return results, nil
}

func performTest(t *testing.T, testName string, expectedAlerts []reflect.Type) {
	previousDir, _ := os.Getwd()

	path := previousDir + "/tests/" + testName
	os.Chdir(path)

	t.Run("Valid", func(t *testing.T) {
		results, err := performParsing(t, path, "valid")
		if err != nil {
			t.Error(err)
			return
		}
		if results.hasAlerts {
			t.Errorf("[Valid] Found alerts in *valid* input")
			alerts.PrintAlerts(t, "Unexpected", results.alerts...)
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		results, err := performParsing(t, path, "invalid")
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
		if !core.ListsAreSame(alertTypes, expectedAlerts) {
			t.Errorf("[Invalid] Mismatch in *expected* and *received* alerts")

			alerts.PrintAlerts(t, "Expected", expectedAlerts...)
			alerts.PrintAlerts(t, "Received", results.alerts...)
		}
	})

	os.Chdir(previousDir)
}

func TestExpressions(t *testing.T) {
	expectedAlerts := []reflect.Type{
		reflect.TypeFor[alerts.ExpectedExpression](),
		reflect.TypeFor[alerts.UnexpectedKeyword](),
		reflect.TypeFor[alerts.ExpectedExpression](),
		reflect.TypeFor[alerts.UnknownStatement](),
		reflect.TypeFor[alerts.ExpectedExpression](),
		reflect.TypeFor[alerts.ExpectedSymbol](),
	}
	performTest(t, "expressions", expectedAlerts)
}
