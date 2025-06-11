package evaluator

import (
	"fmt"
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/core"
	"hybroid/generator"
	"hybroid/lexer"
	"hybroid/parser"
	"hybroid/walker"
	"os"
	"path/filepath"
	"time"

	color "github.com/mitchellh/colorstring"
)

type Evaluator struct {
	walkers    map[string]*walker.Walker
	walkerList []*walker.Walker
	files      []core.FileInformation
	printer    alerts.Printer
}

func NewEvaluator(files []core.FileInformation) Evaluator {
	evaluator := Evaluator{
		walkers:    make(map[string]*walker.Walker),
		walkerList: make([]*walker.Walker, 0),
		files:      files,
		printer:    alerts.NewPrinter(),
	}

	for _, file := range evaluator.files {
		evaluator.walkerList = append(evaluator.walkerList, walker.NewWalker(file.Path(), file.NewPath("/dynamic", ".lua")))
	}

	return evaluator
}

func (e *Evaluator) GetAlerts(sourcePath string) []alerts.Alert {
	return e.printer.GetAlerts(sourcePath)
}

func (e *Evaluator) Action(cwd, outputDir string) error {
	generate := true

	walker.SetupLibraryEnvironments()

	for i := range e.walkerList {
		sourcePath := e.files[i].Path()
		sourceFile, err := os.OpenFile(filepath.Join(cwd, sourcePath), os.O_RDONLY, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to open source file: %v", err)
		}
		defer sourceFile.Close()

		color.Printf("[dark_gray]-->File: %s\n", sourcePath)

		start := time.Now()

		lexer := lexer.NewLexer(sourceFile)
		tokens, tokenizeErr := lexer.Tokenize()

		fmt.Printf("Tokenizing time: %f seconds\n\n", time.Since(start).Seconds())
		e.printer.StageAlerts(sourcePath, lexer.GetAlerts())
		start = time.Now()

		if tokenizeErr != nil {
			generate = false
			continue
		}

		fmt.Printf("Parsing %d tokens\n", len(tokens))

		parser := parser.NewParser(tokens)
		program := parser.Parse()
		fmt.Printf("Parsing time: %f seconds\n\n", time.Since(start).Seconds())
		e.printer.StageAlerts(sourcePath, parser.GetAlerts())

		for _, v := range parser.GetAlerts() {
			if v.AlertType() == alerts.Error {
				generate = false
				break
			}
		}

		// ast.DrawNodes(prog)

		color.Printf("[dark_gray]-->File: %s\n", sourcePath)

		start = time.Now()

		e.walkerList[i].SetProgram(program)
		fmt.Println("Prewalking environments...")
		e.walkerList[i].PreWalk(e.walkers)
		fmt.Printf("Prewalking time: %f seconds\n\n", time.Since(start).Seconds())
	}

	for i, walker := range e.walkerList {
		sourcePath := e.files[i].Path()
		color.Printf("[dark_gray]-->File: %s\n", sourcePath)

		start := time.Now()

		fmt.Println("Walking through the nodes...")
		if !walker.Walked {
			walker.Walk()
		}
		fmt.Printf("Walking time: %f seconds\n\n", time.Since(start).Seconds())

		e.printer.StageAlerts(sourcePath, walker.GetAlerts())
		for _, v := range walker.GetAlerts() {
			if v.AlertType() == alerts.Error {
				generate = false
				break
			}
		}
	}

	if !generate {
		e.printer.PrintAlerts()
		return nil
	}

	fmt.Printf("-Preparing values for generation...\n")
	gen := generator.NewGenerator()
	for _, walker := range e.walkerList {
		gen.SetUniqueEnvName(walker.Env().Name)
	}

	for i, walker := range e.walkerList {
		sourcePath := e.files[i].Path()
		color.Printf("[dark_gray]-->File: %s\n", sourcePath)

		start := time.Now()
		fmt.Println("Generating the lua code...")

		gen.SetEnv(walker.Env().Name, walker.Env().Type)
		if e.files[i].FileName == "level" {
			gen.GenerateWithBuiltins(walker.Program())
		} else if e.walkerList[i].Env().Type != ast.LevelEnv {
			gen.Generate(walker.Program(), e.walkerList[i].Env().UsedBuiltinVars)
		} else {
			gen.Generate(walker.Program(), []string{})
		}

		e.printer.StageAlerts(e.files[i].Path(), gen.GetAlerts())

		fmt.Printf("Generating time: %f seconds\n\n", time.Since(start).Seconds())

		err := os.MkdirAll(filepath.Join(cwd, outputDir, e.files[i].DirectoryPath), os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to write transpiled file to destination: %v", err)
		}
		err = os.WriteFile(e.files[i].NewPath(filepath.Join(cwd, outputDir), ".lua"), []byte(gen.GetSrc()), os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to write transpiled file to destination: %v", err)
		}

		gen = generator.NewGenerator()
	}

	e.printer.PrintAlerts()
	generator.ResetGlobalGeneratorValues()

	return nil
}
