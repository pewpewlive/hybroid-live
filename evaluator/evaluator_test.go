package evaluator

import (
	"hybroid/generator"
	"hybroid/helpers"
	"os"
	"strings"
	"testing"
)

var testsFolder = "/test/"
var cwd = ""
var testFolderName = ""

func newEval(t *testing.T) {
	eval := NewEvaluator(generator.Generator{})

	cwd, _ = os.Getwd()

	files, err := helpers.CollectFiles(cwd + testsFolder + testFolderName)

	if err != nil {
		t.Errorf("Test case falied: file collection error: %v", err)
		t.FailNow()
	}

	for _, file := range files {
		eval.AssignFile(file)
	}

	evalErr := eval.Action(cwd+testsFolder+testFolderName, "/gen")
	if evalErr != nil {
		t.Errorf("Test case falied: evaluation error: %v", err)
		t.FailNow()
	}
}

func readFile(fileName string, t *testing.T) string {
	source, err := os.ReadFile(cwd + testsFolder + testFolderName + "/gen/" + fileName + ".lua")
	if err != nil {
		t.Errorf("failed to write transpiled file to destination: %v", err)
	}

	return string(source)
}

func readExpectedFile(fileName string, t *testing.T) string {
	source, err := os.ReadFile(cwd + testsFolder + testFolderName + "/expected/" + fileName + ".lua")
	if err != nil {
		t.Errorf("failed to write transpiled file to destination: %v", err)
	}

	return string(source)
}

func checkFiles(generated, expected string, t *testing.T) {
	if generated != expected {
		t.Errorf("Test case failed: expected %v\ngot\n%v", generated, expected)
		t.Fail()
	}
}

// useful minifier code
func minify(str string) string {
	for i := 0; i < len(str); i++ {
		switch rune(str[i]) {
		case '\t', '\n', ' ', '\r':
			if i == 0 {
				continue
			}
			switch str[i-1] {
			case '\t', '\n', ' ', '\r':
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

	str = strings.Trim(str, "\t \n\r")

	str = strings.Replace(str, "\r", " ", -1)
	str = strings.Replace(str, "\n", " ", -1)
	str = strings.Replace(str, "\t", " ", -1)

	return str
}

func check(fileName string, t *testing.T) {
	expected := readExpectedFile(fileName, t)
	generated := readFile(fileName, t)
	minifiedGenerated := minify(generated)
	minifiedExpected := minify(expected)
	if minifiedGenerated != minifiedExpected {
		t.Errorf("Test case failed: expected %v\ngot\n%v", minifiedGenerated, minifiedExpected)
		t.Fail()
	}
}

var tests = []string{
	"test1",
	"test_statements",
	"test_expressions",
	"test_string",
}

// sample test
func TestAll(t *testing.T) {
	for _, test := range tests {
		testFolderName = test
		t.Run(test, func(t *testing.T) {
			newEval(t)
			check(test, t)
		})
	}
}
