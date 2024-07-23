package evaluator

import (
	"fmt"
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
	"strings"
	"time"

	"github.com/mitchellh/colorstring"
)

type Evaluator struct {
	walkers    map[string]*walker.Walker
	walkerList map[helpers.FileInformation]*walker.Walker

	// Toolset
	lexer  lexer.Lexer
	parser parser.Parser
	gen    generator.Generator
}

func NewEvaluator(gen generator.Generator) Evaluator {
	return Evaluator{
		walkers:    make(map[string]*walker.Walker),
		walkerList: make(map[helpers.FileInformation]*walker.Walker),
		lexer:      lexer.NewLexer(),
		parser:     parser.NewParser(),
		gen:        gen,
	}
}

func (e *Evaluator) AssignFile(file helpers.FileInformation) {
	e.walkerList[file] = walker.NewWalker(file.NewPath("/dynamic", ".lua"))
}

func (e *Evaluator) Action(cwd, outputDir string) error {
	walker.SetupLibraryEnvironments()

	for file := range e.walkerList {
		sourceFile, err := os.ReadFile(filepath.Join(cwd, file.Path()))
		if err != nil {
			return fmt.Errorf("failed to read source file: %v", err)
		}

		fmt.Printf("Tokenizing %v characters\n", len(sourceFile))
		start := time.Now()

		e.lexer.AssignSource(sourceFile)
		e.lexer.Tokenize()
		if len(e.lexer.Errors) != 0 {
			fmt.Println("[red]Failed tokenizing:")
			printAlerts(file.Path(), e.lexer.Errors)
			return fmt.Errorf("failed to tokenize source file")
		}

		fmt.Printf("Tokenizing time: %v seconds\n\n", time.Since(start).Seconds())
		start = time.Now()

		fmt.Printf("Parsing %v tokens\n", len(e.lexer.Tokens))

		e.parser.AssignTokens(e.lexer.Tokens)
		prog := e.parser.ParseTokens()
		if len(e.parser.Errors) != 0 {
			colorstring.Println("[red]Syntax error found:")
			for _, err := range e.parser.Errors {
				e.writeSyntaxAlert(string(sourceFile), err)
				colorstring.Printf("[red]Error: %+v\n", err)
			}
			return fmt.Errorf("failed to parse source file")
		}
		fmt.Printf("Parsing time: %v seconds\n\n", time.Since(start).Seconds())

		start = time.Now()
		fmt.Println("[Pass 1] Walking through the nodes...")
		pass1.Action(e.walkerList[file], prog, e.walkers)
		fmt.Printf("Pass 1 time: %v seconds\n\n", time.Since(start).Seconds())

		e.lexer = lexer.NewLexer()
		e.parser = parser.NewParser()
	}

	for file, walker := range e.walkerList {
		start := time.Now()
		fmt.Println("[Pass 2] Walking through the nodes...")

		pass2.Action(walker, e.walkers)
		fmt.Printf("Pass 2 time: %v seconds\n\n", time.Since(start).Seconds())
		if len(walker.Errors) != 0 {
			colorstring.Printf("[red]Failed walking (%s):\n", file.Path())
			printAlerts(file.Path(), walker.Errors)
			return fmt.Errorf("failed to walk through the nodes")
		}
		if len(walker.Warnings) != 0 {
			printAlerts(file.Path(), walker.Warnings)
		}
	}

	fmt.Println("Preparing values for generation...")
	for _, walker := range e.walkerList {
		e.gen.SetUniqueEnvName(walker.Environment.Name)
	}

	for file, walker := range e.walkerList {
		start := time.Now()
		fmt.Println("Generating the lua code...")

		//e.gen.Scope.Src.Grow(len(sourceFile))
		e.gen.SetEnvName(walker.Environment.Name)
		e.gen.Generate(walker.Nodes)
		if len(e.gen.Errors) != 0 {
			colorstring.Println("[red]Failed generating:")
			printAlerts(file.Path(), e.gen.GetErrors())
		}
		fmt.Printf("Generating time: %v seconds\n", time.Since(start).Seconds())

		err := os.MkdirAll(filepath.Join(cwd, outputDir, file.DirectoryPath), os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to write transpiled file to destination: %v", err)
		}
		err = os.WriteFile(file.NewPath(filepath.Join(cwd, outputDir), ".lua"), []byte(e.gen.GetSrc()), 0644)
		if err != nil {
			return fmt.Errorf("failed to write transpiled file to destination: %v", err)
		}

		e.gen.Clear()
	}

	return nil
}

func printAlerts[T ast.Alert](filePath string, errs []T) {
	for _, err := range errs {
		tokenLocation := err.GetToken().Location
		str := fmt.Sprintf("%s in %s at line %v: %s\n",
			err.GetHeader(),
			filePath,
			tokenLocation.LineStart,
			err.GetMessage())
		fmt.Print(colorstring.Color(str))
	}
	fmt.Println()
}

func (e *Evaluator) writeSyntaxAlert(source string, errMsg ast.Alert) {
	token := errMsg.GetToken()

	sourceLines := strings.Split(source, "\n")
	line := sourceLines[token.Location.LineStart-1]

	fmt.Printf("line: %v in file \n", token.Location.LineStart)
	fmt.Println(line)
	if token.Location.ColStart-6 < 0 {
		fmt.Printf("%s^%s\n", strings.Repeat(" ", token.Location.ColStart-1), strings.Repeat("-", 5))
	} else {
		fmt.Printf("%s%s^\n", strings.Repeat(" ", token.Location.ColStart-6), strings.Repeat("-", 5))
	}
	fmt.Println("message: " + errMsg.GetMessage() + "\n")
}
