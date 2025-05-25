package evaluator

import (
	"fmt"
	"hybroid/alerts"
	"hybroid/core"
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

func (e *Evaluator) Action(cwd, outputDir string) error {
	evalFailed := make([]bool, 0)

	walker.SetupLibraryEnvironments()

	for i := range e.walkerList {
		evalFailed = append(evalFailed, false)

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
			evalFailed[i] = true
			continue
		}

		fmt.Printf("Parsing %d tokens\n", len(tokens))

		parser := parser.NewParser(tokens)
		program := parser.Parse()
		for _, v := range parser.GetAlerts() {
			if v.GetAlertType() == alerts.Error {
				evalFailed[i] = true
				break
			}
		}
		fmt.Printf("Parsing time: %f seconds\n\n", time.Since(start).Seconds())
		e.printer.StageAlerts(sourcePath, parser.GetAlerts())

		// ast.DrawNodes(prog)

		color.Printf("[dark_gray]-->File: %s\n", sourcePath)

		start = time.Now()

		e.walkerList[i].SetProgram(program)
		fmt.Println("[Pass 1] Walking environments...")
		e.walkerList[i].Pass1(e.walkers)
		fmt.Printf("Pass 1 time: %f seconds\n\n", time.Since(start).Seconds())
	}

	for i, walker := range e.walkerList {
		sourcePath := e.files[i].Path()
		if walker.Walked {
			e.printer.StageAlerts(sourcePath, walker.GetAlerts())
			continue
		}
		color.Printf("[dark_gray]-->File: %s\n", sourcePath)

		start := time.Now()

		fmt.Println("[Pass 2] Walking through the nodes...")
		walker.Pass2()
		fmt.Printf("Pass 2 time: %f seconds\n\n", time.Since(start).Seconds())

		e.printer.StageAlerts(sourcePath, walker.GetAlerts())
	}

	// fmt.Printf("-Preparing values for generation...\n")
	// generator := generator.NewGenerator()
	// for _, walker := range e.walkerList {
	// 	generator.SetUniqueEnvName(walker.Environment.Name)
	// }

	// for i, walker := range e.walkerList {
	// 	cont := false
	// 	for _, v := range walker.GetAlerts() {
	// 		if v.GetAlertType() == alerts.Error {
	// 			cont = true
	// 			break
	// 		}
	// 	}
	// 	if evalFailed[i] || cont {
	// 		continue
	// 	}
	// 	sourcePath := e.files[i].Path()
	// 	color.Printf("[dark_gray]-->File: %s\n", sourcePath)

	// 	start := time.Now()
	// 	fmt.Println("Generating the lua code...")

	// 	generator.SetEnv(walker.Environment.Name, walker.Environment.Type)
	// 	if e.files[i].FileName == "level" {
	// 		//generator.GenerateWithBuiltins(walker.Nodes)
	// 	} else if e.walkerList[i].Environment.Type != ast.LevelEnv {
	// 		//generator.Generate(walker.Nodes, e.walkerList[i].Environment.UsedBuiltinVars)
	// 	} else {
	// 		//generator.Generate(walker.Nodes, []string{})
	// 	}

	// 	e.printer.StageAlerts(e.files[i].Path(), generator.GetAlerts())

	// 	fmt.Printf("Generating time: %f seconds\n\n", time.Since(start).Seconds())

	// 	err := os.MkdirAll(filepath.Join(cwd, outputDir, e.files[i].DirectoryPath), os.ModePerm)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to write transpiled file to destination: %v", err)
	// 	}
	// 	err = os.WriteFile(e.files[i].NewPath(filepath.Join(cwd, outputDir), ".lua"), []byte(generator.GetSrc()), os.ModePerm)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to write transpiled file to destination: %v", err)
	// 	}

	// 	generator.Clear()
	// }

	e.printer.PrintAlerts()

	return nil
}
