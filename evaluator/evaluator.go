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
	"strings"
)

type Evaluator struct {
	// walkers map environment names AND absolute paths to walker instances
	walkers    map[string]*walker.Walker
	walkerList []*walker.Walker
	files      []core.FileInformation
	programs   map[string][]ast.Node
	printer    alerts.Printer
}

func NewEvaluator(files []core.FileInformation) *Evaluator {
	evaluator := &Evaluator{
		walkers:    make(map[string]*walker.Walker),
		walkerList: make([]*walker.Walker, 0),
		files:      files,
		programs:   make(map[string][]ast.Node),
		printer:    alerts.NewPrinter(),
	}

	for _, file := range evaluator.files {
		w := walker.NewWalker(file.Path(), file.NewPath("/dynamic", ".lua"))
		evaluator.walkerList = append(evaluator.walkerList, w)
		// Index by path initially
		abs, err := filepath.Abs(file.Path())
		if err == nil {
			evaluator.walkers[abs] = w
		} else {
			evaluator.walkers[file.Path()] = w
		}
	}

	return evaluator
}

func (e *Evaluator) GetAlerts(sourcePath string) []alerts.Alert {
	return e.printer.GetAlerts(sourcePath)
}

// ParseAll reads and parses all files in the evaluator's list from disk.
func (e *Evaluator) ParseAll(cwd string) error {
	for i, w := range e.walkerList {
		sourcePath := e.files[i].Path()
		sourceFile, err := os.OpenFile(filepath.Join(cwd, sourcePath), os.O_RDONLY, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to open source file: %v", err)
		}

		lex := lexer.NewLexer(sourceFile)
		tokens, tokenizeErr := lex.Tokenize()
		sourceFile.Close()

		e.printer.StageAlerts(sourcePath, lex.GetAlerts())
		if tokenizeErr != nil {
			continue
		}

		p := parser.NewParser(tokens)
		program := p.Parse()
		e.printer.StageAlerts(sourcePath, p.GetAlerts())

		w.SetProgram(program)
		e.programs[sourcePath] = program
	}
	return nil
}

// RunAnalysis performs the PreWalk and Walk phases across all files.
func (e *Evaluator) RunAnalysis() {
	walker.SetupLibraryEnvironments()
	e.printer = alerts.NewPrinter() // Clear previous alerts

	// Pass 0: Reset all walkers and remove old environment names from the map
	// We keep the absolute paths in e.walkers
	newWalkers := make(map[string]*walker.Walker)
	
	for _, w := range e.walkerList {
		w.Reset()
		abs, err := filepath.Abs(w.Env().HybroidPath())
		if err == nil {
			newWalkers[abs] = w
		} else {
			newWalkers[w.Env().HybroidPath()] = w
		}
	}
	e.walkers = newWalkers

	// Pass 1: PreWalk (Registers environment names in e.walkers)
	for _, w := range e.walkerList {
		w.PreWalk(e.walkers)
		if w.Env().Name != "" {
			e.walkers[w.Env().Name] = w
		}
	}

	// Pass 2: Walk
	for _, w := range e.walkerList {
		if !w.Walked {
			w.Walk()
		}
	}

	// Pass 3: PostWalk
	for i, w := range e.walkerList {
		w.PostWalk()
		e.printer.StageAlerts(e.files[i].Path(), w.GetAlerts())
	}
}

// Action maintains the exact same build process as before, but uses the refactored phases.
func (e *Evaluator) Action(cwd, outputDir string) error {
	err := e.ParseAll(cwd)
	if err != nil {
		return err
	}

	// Check for errors before walking
	if e.hasErrors() {
		e.printer.PrintAlerts()
		return nil
	}

	e.RunAnalysis()

	if e.hasErrors() {
		e.printer.PrintAlerts()
		return nil
	}

	return e.EmitLua(cwd, outputDir)
}

func (e *Evaluator) hasErrors() bool {
	for _, fileAlerts := range e.printer.AllAlerts() {
		for _, a := range fileAlerts {
			if a.AlertType() == alerts.Error {
				return true
			}
		}
	}
	return false
}

// EmitLua handles the Lua code generation and file writing.
func (e *Evaluator) EmitLua(cwd, outputDir string) error {
	outputPath := filepath.Join(cwd, outputDir)
	if outputDir != "" {
		if stat, err := os.Lstat(outputPath); err == nil && stat.IsDir() {
			os.RemoveAll(outputPath)
		}
	}

	gen := generator.NewGenerator()
	for _, w := range e.walkerList {
		gen.SetUniqueEnvName(w.Env().Name)
	}

	for i, w := range e.walkerList {
		gen.SetEnv(w.Env().Name, w.Env().Type)
		gen.GenerateUsedLibraries(w.Env().UsedLibraries)
		
		if e.files[i].FileName == "level" {
			gen.GenerateWithBuiltins(w.Program())
		} else if w.Env().Type != ast.LevelEnv {
			gen.Generate(w.Program(), w.Env().UsedBuiltinVars)
		} else {
			gen.Generate(w.Program(), []string{})
		}

		e.printer.StageAlerts(e.files[i].Path(), gen.GetAlerts())

		err := os.MkdirAll(filepath.Join(outputPath, e.files[i].DirectoryPath), os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to write transpiled file to destination: %v", err)
		}
		
		// Fix: .lua extension logic from original
		luaPath := e.files[i].NewPath(outputPath, ".lua")
		err = os.WriteFile(luaPath, []byte(gen.GetSrc()), os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to write transpiled file to destination: %v", err)
		}

		gen = generator.NewGenerator()
	}

	e.printer.PrintAlerts()
	generator.ResetGlobalGeneratorValues()
	return nil
}

// UpdateFileContent parses a specific file from a string (in-memory) instead of disk.
func (e *Evaluator) UpdateFileContent(path string, content string) error {
	lex := lexer.NewLexer(strings.NewReader(content))
	tokens, tokenizeErr := lex.Tokenize()
	e.printer.StageAlerts(path, lex.GetAlerts())
	if tokenizeErr != nil {
		return tokenizeErr
	}

	p := parser.NewParser(tokens)
	program := p.Parse()
	e.printer.StageAlerts(path, p.GetAlerts())

	// Find the walker for this path
	abs, _ := filepath.Abs(path)
	if w, ok := e.walkers[abs]; ok {
		w.SetProgram(program)
	} else if w, ok := e.walkers[path]; ok {
		w.SetProgram(program)
	}
	e.programs[path] = program
	return nil
}

// AnalyzeFile re-runs analysis for a specific file and returns its walker.
func (e *Evaluator) AnalyzeFile(path string) *walker.Walker {
	// For now, we re-run full project analysis to ensure cross-file consistency.
	// This can be optimized later to be incremental.
	e.RunAnalysis()

	abs, _ := filepath.Abs(path)
	if w, ok := e.walkers[abs]; ok {
		return w
	}
	return e.walkers[path]
}

func (e *Evaluator) Walkers() map[string]*walker.Walker {
	return e.walkers
}
