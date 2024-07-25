package evaluator

import (
	"hybroid/generator"
	"hybroid/helpers"
	"os"
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
func check(fileName string, t *testing.T) {
	expected := readExpectedFile(fileName, t)
	generated := readFile(fileName, t)

	if generated != expected {
		t.Errorf("Test case failed: expected %v\ngot\n%v", generated, expected)
		t.Fail()
	}
}

// sample test
func Test1(t *testing.T) {
	testFolderName = "test1"
	newEval(t)

	check("test1", t)
}
