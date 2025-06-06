package evaluator

import (
	"hybroid/alerts"
	"hybroid/core"
	"os"
	"strings"
	"testing"
)

var testsFolder = "/test/"
var cwd = ""
var testFolderName = ""

func newEval(t *testing.T) {

	cwd, _ = os.Getwd()

	files, err := core.CollectFiles(cwd + testsFolder + testFolderName)

	if err != nil {
		t.Errorf("Test case falied: file collection error: %v", err)
		t.FailNow()
	}

	eval := NewEvaluator(files)

	evalErr := eval.Action(cwd+testsFolder+testFolderName, "")
	if evalErr != nil {
		t.Errorf("Test case falied: evaluation error: %v", err)
		t.FailNow()
	}

	alrts := eval.GetAlerts("test.hyb")

	if len(alrts) != 0 {
		t.Errorf("Unexpected alerts found: ")
		alerts.PrintAlerts(t, "Unexpected", alrts...)
		t.FailNow()
	}
}

func readFile(t *testing.T) string {
	source, err := os.ReadFile(cwd + testsFolder + testFolderName + "/test.lua")
	if err != nil {
		t.Errorf("failed to write transpiled file to destination: %v", err)
		t.FailNow()
	}

	return string(source)
}

func readExpectedFile(t *testing.T) string {
	source, err := os.ReadFile(cwd + testsFolder + testFolderName + "/expected.lua")
	if err != nil {
		t.Errorf("failed to write transpiled file to destination: %v", err)
		t.FailNow()
	}

	return string(source)
}

// useful minifier code
func minify(str string) string {
	str = strings.Trim(str, "\t \n\r")

	str = strings.Replace(str, "\r", " ", -1)
	str = strings.Replace(str, "\n", " ", -1)
	str = strings.Replace(str, "\t", " ", -1)

	for i := 0; i < len(str); i++ {
		switch rune(str[i]) {
		case ' ':
			if i == 0 {
				continue
			}
			switch str[i-1] {
			case ' ':
				if i == len(str)-1 {
					str = str[:i-1]
				} else if i == 0 {
					str = str[i+1:]
				} else {
					str = str[:i] + str[i+1:]
				}
				i--
			}
		}
	}

	return str
}

func check(t *testing.T) {
	expected := readExpectedFile(t)
	generated := readFile(t)
	minifiedGenerated := minify(generated)
	minifiedExpected := minify(expected)
	if minifiedGenerated != minifiedExpected {
		t.Errorf("Test case failed: expected\n%v\ngot\n%v", minifiedExpected, minifiedGenerated)
		t.FailNow()
	}
}

func TestStatements(t *testing.T) {
	testFolderName = "statements"

	newEval(t)
	check(t)
}

func TestExpressions(t *testing.T) {
	testFolderName = "expressions"

	newEval(t)
	check(t)
}
