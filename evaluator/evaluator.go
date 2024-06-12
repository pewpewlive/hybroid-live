package evaluator

import (
	"fmt"
	"hybroid/ast"
	"hybroid/walker"
	"strings"
	"time"

	"hybroid/generators/lua"
	"hybroid/lexer"
	"hybroid/parser"
	"os"

	"github.com/mitchellh/colorstring"
)

type Evaluator struct {
	walkers *map[string]*walker.Walker
	lexer      *lexer.Lexer
	parser     *parser.Parser
	walker     *walker.Walker
	gen        lua.Generator
	SrcPath    string
	DstPath    string
}

func NewEvaluator(gen lua.Generator, walkers *map[string]*walker.Walker) Evaluator {
	return Evaluator{
		walkers: walkers,
		lexer:      lexer.NewLexer(),
		parser:     parser.NewParser(),
		gen:        gen,
	}
}

func (e *Evaluator) AssignFile(src string, dst string) {
	e.SrcPath, e.DstPath = src, dst
	e.walker = walker.NewWalker(e.SrcPath)
}

func (e *Evaluator) Action(writeEnabled bool) (string, error) {
	sourceFile, err := os.ReadFile(e.SrcPath)
	if err != nil {
		return "", fmt.Errorf("failed to read source file: %v", err)
	}

	fmt.Printf("Tokenizing %v characters\n", len(sourceFile))
	start := time.Now()

	e.lexer.AssignSource(sourceFile)
	e.lexer.Tokenize()
	if len(e.lexer.Errors) != 0 {
		fmt.Println("[red]Failed tokenizing:")
		printAlerts(e.lexer.Errors)
		return "", fmt.Errorf("failed to tokenize source file")
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
			//colorstring.Printf("[red]Error: %+v\n", err)
		}
		return "", fmt.Errorf("failed to parse source file")
	}
	fmt.Printf("Parsing time: %v seconds\n\n", time.Since(start).Seconds())
	start = time.Now()

	fmt.Println("Walking through the nodes...")

	prog = e.walker.Pass1(&prog, e.walkers)
	if len(e.walker.Errors) != 0 {
		colorstring.Println("[red]Failed walking:")
		printAlerts(e.walker.Errors)
		return "", fmt.Errorf("failed to walk through the nodes")
	}
	if len(e.walker.Warnings) != 0 {
		printAlerts(e.walker.Warnings)
	}

	fmt.Printf("Walking time: %v seconds\n\n", time.Since(start).Seconds())
	start = time.Now()

	fmt.Println("Generating the lua code...")

	//e.gen.Scope.Src.Grow(len(sourceFile))
	e.gen.Generate(prog)
	if len(e.gen.Errors) != 0 {
		colorstring.Println("[red]Failed generating:")
		printAlerts(e.gen.GetErrors())
	}
	fmt.Printf("Generating time: %v seconds\n", time.Since(start).Seconds())

	if !writeEnabled {
		return e.gen.GetSrc(), nil
	}
	err = os.WriteFile(e.DstPath, []byte(e.gen.GetSrc()), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write transpiled file to destination: %v", err)
	}

	return e.gen.GetSrc(), nil
}

func printAlerts[T ast.Alert](errs []T) {
	for _, err := range errs {
		tokenLocation := err.GetToken().Location
		str := fmt.Sprintf("%v at %v:%v-%v: %s\n", 
			err.GetHeader(),
			tokenLocation.LineStart, 
			tokenLocation.ColStart, 
			tokenLocation.ColEnd, 
			err.GetMessage())
		fmt.Print(colorstring.Color(str))
	}
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
