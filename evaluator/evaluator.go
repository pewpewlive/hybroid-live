package evaluator

import (
	"fmt"
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/generator"
	"hybroid/helpers"
	"hybroid/lexer"
	"hybroid/parser"
	"hybroid/walker"
	"hybroid/walker/pass1"
	"hybroid/walker/pass2"
	"os"
	"path/filepath"
	"time"
)

type Evaluator struct {
	walkers    map[string]*walker.Walker
	walkerList []*walker.Walker
	files      []helpers.FileInformation

	// Toolset
	lexer  lexer.Lexer
	parser parser.Parser
	gen    generator.Generator
}

func NewEvaluator(gen generator.Generator) Evaluator {
	return Evaluator{
		walkers:    make(map[string]*walker.Walker),
		walkerList: make([]*walker.Walker, 0),
		lexer:      lexer.NewLexer(),
		parser:     parser.NewParser(),
		gen:        gen,
	}
}

func (e *Evaluator) AssignFile(file helpers.FileInformation) {
	e.files = append(e.files, file)
	e.walkerList = append(e.walkerList, walker.NewWalker(file.NewPath("/dynamic", ".lua")))
}

func (e *Evaluator) Action(cwd, outputDir string) error {
	walker.SetupLibraryEnvironments()

	for i := range e.walkerList {
		sourcePath := e.files[i].Path()
		sourceFile, err := os.ReadFile(filepath.Join(cwd, sourcePath))
		if err != nil {
			return fmt.Errorf("failed to read source file: %v", err)
		}

		fmt.Printf("Tokenizing %d characters\n", len(sourceFile))
		start := time.Now()

		e.lexer.AssignSource(sourceFile)
		e.lexer.Tokenize()
		if e.lexer.HasAlerts {
			e.lexer.PrintAlerts(alerts.Lexer, sourceFile, sourcePath)
			return fmt.Errorf("failed to tokenize source file")
		}

		fmt.Printf("Tokenizing time: %f seconds\n\n", time.Since(start).Seconds())
		start = time.Now()

		fmt.Printf("Parsing %d tokens\n", len(e.lexer.Tokens))

		e.parser.AssignTokens(e.lexer.Tokens)
		prog := e.parser.ParseTokens()
		if e.parser.HasAlerts {
			e.parser.PrintAlerts(alerts.Parser, sourceFile, sourcePath)

		}
		// if len(e.parser.Errors) != 0 {
		// 	color.Println("[red]Syntax error")
		// 	for _, err := range e.parser.Errors {
		// 		color.Printf("[red]Error: %+v\n", err)
		// 	}
		// }
		if e.parser.HasAlerts /*|| len(e.parser.Errors) != 0*/ {
			return fmt.Errorf("failed to parse source file")
		}
		fmt.Printf("Parsing time: %f seconds\n\n", time.Since(start).Seconds())

		//ast.DrawNodes(prog)

		start = time.Now()
		fmt.Println("[Pass 1] Walking through the nodes...")
		if env, ok := prog[0].(*ast.EnvironmentStmt); ok {
			e.walkerList[i].Environment.Type = env.EnvType.Type
		}
		pass1.Action(e.walkerList[i], prog, e.walkers)
		fmt.Printf("Pass 1 time: %f seconds\n\n", time.Since(start).Seconds())

		e.lexer = lexer.NewLexer()
		e.parser = parser.NewParser()
	}

	for _, walker := range e.walkerList {

		start := time.Now()
		fmt.Println("[Pass 2] Walking through the nodes...")

		if !walker.Walked {
			pass2.Action(walker, e.walkers)
		}
		fmt.Printf("Pass 2 time: %f seconds\n\n", time.Since(start).Seconds())
		// if walker.HasAlerts == true {
		// 	color.Printf("[red]Failed walking (%s):\n", e.files[i].Path())

		// 	printAlerts(e.files[i].Path(), walker.Errors)
		// 	return fmt.Errorf("failed to walk through the nodes")
		// }
	}

	fmt.Println("Preparing values for generation...")
	for _, walker := range e.walkerList {
		e.gen.SetUniqueEnvName(walker.Environment.Name)
	}

	for i, walker := range e.walkerList {
		start := time.Now()
		fmt.Println("Generating the lua code...")

		//e.gen.Scope.Src.Grow(len(sourceFile))
		e.gen.SetEnv(walker.Environment.Name, walker.Environment.Type)
		if e.files[i].FileName == "level" {
			e.gen.GenerateWithBuiltins(walker.Nodes)
		} else if e.walkerList[i].Environment.Type != ast.Level {
			e.gen.Generate(walker.Nodes, e.walkerList[i].Environment.UsedBuiltinVars)
		} else {
			e.gen.Generate(walker.Nodes, []string{})
		}
		// if len(e.gen.Errors) != 0 {
		// 	color.Println("[red]Failed generating:")
		// 	printAlerts(e.files[i].Path(), e.gen.GetErrors())
		// }
		fmt.Printf("Generating time: %f seconds\n", time.Since(start).Seconds())

		err := os.MkdirAll(filepath.Join(cwd, outputDir, e.files[i].DirectoryPath), os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to write transpiled file to destination: %v", err)
		}
		err = os.WriteFile(e.files[i].NewPath(filepath.Join(cwd, outputDir), ".lua"), []byte(e.gen.GetSrc()), 0644)
		if err != nil {
			return fmt.Errorf("failed to write transpiled file to destination: %v", err)
		}

		e.gen.Clear()
	}

	e.walkerList = make([]*walker.Walker, 0)
	e.files = make([]helpers.FileInformation, 0)
	e.walkers = make(map[string]*walker.Walker)

	return nil
}
