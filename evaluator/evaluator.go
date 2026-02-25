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
	"sync"
)

type Evaluator struct {
	mu sync.Mutex
	// walkers map environment names AND absolute paths to walker instances
	walkers     map[string]*walker.Walker
	walkerList  []*walker.Walker
	files       []core.FileInformation
	programs    map[string][]ast.Node
	parseAlerts map[string][]alerts.Alert
	printer     alerts.Printer
}

func NewEvaluator(files []core.FileInformation) *Evaluator {
	evaluator := &Evaluator{
		walkers:     make(map[string]*walker.Walker),
		walkerList:  make([]*walker.Walker, 0),
		files:       files,
		programs:    make(map[string][]ast.Node),
		parseAlerts: make(map[string][]alerts.Alert),
		printer:     alerts.NewPrinter(),
	}

	for _, file := range evaluator.files {
		w := walker.NewWalker(file.Path(), file.NewPath("/dynamic", ".lua"))
		evaluator.walkerList = append(evaluator.walkerList, w)
		// Index by path initially
		abs, err := filepath.Abs(file.Path())
		if err == nil {
			evaluator.walkers[abs] = w
		}
		evaluator.walkers[file.Path()] = w
	}

	return evaluator
}

func (e *Evaluator) GetAlerts(sourcePath string) []alerts.Alert {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.printer.GetAlerts(e.canonicalPath(sourcePath))
}

func (e *Evaluator) canonicalPath(path string) string {
	path = filepath.ToSlash(filepath.Clean(path))
	for _, file := range e.files {
		sourcePath := filepath.ToSlash(filepath.Clean(file.Path()))
		if sourcePath == path {
			return sourcePath
		}
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	absPath = filepath.ToSlash(filepath.Clean(absPath))

	for _, file := range e.files {
		sourcePath := filepath.ToSlash(filepath.Clean(file.Path()))
		fileAbs, err := filepath.Abs(sourcePath)
		if err != nil {
			continue
		}
		if filepath.ToSlash(filepath.Clean(fileAbs)) == absPath {
			return sourcePath
		}
	}

	matchCount := 0
	matchPath := ""
	for _, file := range e.files {
		sourcePath := filepath.ToSlash(filepath.Clean(file.Path()))
		if filepath.Base(sourcePath) == filepath.Base(absPath) {
			matchCount++
			matchPath = sourcePath
		}
	}

	if matchCount == 1 {
		return matchPath
	}

	return path
}

// ParseAll reads and parses all files in the evaluator's list from disk.
func (e *Evaluator) ParseAll(cwd string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.parseAll(cwd)
}

func (e *Evaluator) parseAll(cwd string) error {
	for i, w := range e.walkerList {
		sourcePath := e.files[i].Path()
		sourceFile, err := os.OpenFile(filepath.Join(cwd, sourcePath), os.O_RDONLY, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to open source file: %v", err)
		}

		lex := lexer.NewLexer(sourceFile)
		tokens, tokenizeErr := lex.Tokenize()
		sourceFile.Close()

		fileAlerts := make([]alerts.Alert, 0)
		fileAlerts = append(fileAlerts, lex.GetAlerts()...)
		if tokenizeErr != nil {
			e.parseAlerts[sourcePath] = fileAlerts
			e.printer.StageAlerts(sourcePath, fileAlerts)
			continue
		}

		p := parser.NewParser(tokens)
		program := p.Parse()
		fileAlerts = append(fileAlerts, p.GetAlerts()...)
		e.parseAlerts[sourcePath] = fileAlerts
		e.printer.StageAlerts(sourcePath, fileAlerts)

		w.SetProgram(program)
		e.programs[sourcePath] = program
	}
	return nil
}

// RunAnalysis performs the PreWalk and Walk phases across all files.
func (e *Evaluator) RunAnalysis() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.runAnalysis()
}

func (e *Evaluator) runAnalysis() {
	walker.SetupLibraryEnvironments()
	e.printer = alerts.NewPrinter() // Clear previous alerts

	for _, file := range e.files {
		sourcePath := file.Path()
		if parseAlerts, ok := e.parseAlerts[sourcePath]; ok {
			e.printer.StageAlerts(sourcePath, parseAlerts)
		}
	}

	// Pass 0: Reset all walkers and rebuild the mapping from absolute paths
	// This clears any stale environment names from previous runs.
	newWalkers := make(map[string]*walker.Walker)
	for _, w := range e.walkerList {
		w.Reset()
		abs, err := filepath.Abs(w.Env().HybroidPath())
		if err == nil {
			newWalkers[abs] = w
		}
		newWalkers[w.Env().HybroidPath()] = w
	}
	e.walkers = newWalkers

	// Pass 1: PreWalk (Registers environment names in e.walkers)
	for _, w := range e.walkerList {
		w.PreWalk(e.walkers)
		// After PreWalk, the walker has its environment name set if it had an 'env' statement.
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
	e.mu.Lock()
	defer e.mu.Unlock()

	err := e.parseAll(cwd)
	if err != nil {
		return err
	}

	e.runAnalysis()

	if e.hasErrors() {
		e.printer.PrintAlerts()
		return nil
	}

	return e.emitLua(cwd, outputDir)
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
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.emitLua(cwd, outputDir)
}

func (e *Evaluator) emitLua(cwd, outputDir string) error {
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
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.updateFileContent(path, content)
}

func (e *Evaluator) updateFileContent(path string, content string) error {
	path = e.canonicalPath(path)
	lex := lexer.NewLexer(strings.NewReader(content))
	tokens, tokenizeErr := lex.Tokenize()
	fileAlerts := make([]alerts.Alert, 0)
	fileAlerts = append(fileAlerts, lex.GetAlerts()...)
	e.parseAlerts[path] = fileAlerts
	e.printer.StageAlerts(path, fileAlerts)
	if tokenizeErr != nil {
		return tokenizeErr
	}

	p := parser.NewParser(tokens)
	program := p.Parse()
	fileAlerts = append(fileAlerts, p.GetAlerts()...)
	e.parseAlerts[path] = fileAlerts
	e.printer.StageAlerts(path, fileAlerts)

	// Find the walker for this path
	abs, _ := filepath.Abs(path)
	if w, ok := e.walkers[abs]; ok {
		w.SetProgram(program)
	} else if w, ok := e.walkers[path]; ok {
		w.SetProgram(program)
	} else {
		// Fallback: try to match by hybroid path
		for _, w := range e.walkerList {
			wAbs, _ := filepath.Abs(w.Env().HybroidPath())
			if wAbs == abs || w.Env().HybroidPath() == path {
				w.SetProgram(program)
				break
			}
		}
	}
	e.programs[path] = program
	return nil
}

// AnalyzeFile re-runs analysis for a specific file and returns its walker.
func (e *Evaluator) AnalyzeFile(path string) *walker.Walker {
	e.mu.Lock()
	defer e.mu.Unlock()
	// For now, we re-run full project analysis to ensure cross-file consistency.
	// This can be optimized later to be incremental.
	e.runAnalysis()

	canonical := e.canonicalPath(path)
	if w, ok := e.walkers[canonical]; ok {
		return w
	}

	abs, _ := filepath.Abs(path)
	if w, ok := e.walkers[abs]; ok {
		return w
	}
	return e.walkers[path]
}

func (e *Evaluator) Walkers() map[string]*walker.Walker {
	e.mu.Lock()
	defer e.mu.Unlock()
	copyMap := make(map[string]*walker.Walker, len(e.walkers))
	for k, v := range e.walkers {
		copyMap[k] = v
	}
	return copyMap
}

func (e *Evaluator) WalkerList() []*walker.Walker {
	e.mu.Lock()
	defer e.mu.Unlock()
	return append([]*walker.Walker{}, e.walkerList...)
}
