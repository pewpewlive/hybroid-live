package evaluator

import (
	"fmt"
	"hybroid/ast"
	"hybroid/walker"
	"strings"
	"time"

	//"hybroid/generators"
	"hybroid/generators/lua"
	"hybroid/lexer"
	"hybroid/parser"
	"os"

	"github.com/mitchellh/colorstring"
)

type Evaluator struct {
	lexer   *lexer.Lexer
	parser  *parser.Parser
	walker  walker.Walker
	gen     lua.Generator
	SrcPath string
	DstPath string
}

func New(gen lua.Generator) Evaluator {
	return Evaluator{
		lexer:  lexer.New(),
		parser: parser.New(),
		gen:    gen,
	}
}

func (e *Evaluator) AssignFile(src string, dst string) {
	e.SrcPath, e.DstPath = src, dst
}

func (e *Evaluator) Action() error {
	sourceFile, err := os.ReadFile(e.SrcPath)
	if err != nil {
		return fmt.Errorf("failed to read source file: %v", err)
	}

	fmt.Printf("Tokenizing %v characters\n", len(sourceFile))
	start := time.Now()

	e.lexer.AssignSource(sourceFile)
	e.lexer.Tokenize()
	if len(e.lexer.Errors) != 0 {
		fmt.Println("[red]Failed tokenizing:")
		for _, err := range e.lexer.Errors {
			colorstring.Printf("[red]Error: %v\n", err)
		}
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
			e.writeSyntaxError(string(sourceFile), err)
			//colorstring.Printf("[red]Error: %+v\n", err)
		}
		return fmt.Errorf("failed to parse source file")
	}

	fmt.Printf("Parsing time: %v seconds\n\n", time.Since(start).Seconds())
	start = time.Now()

	fmt.Println("Walking through the nodes...")

	global := walker.NewGlobal()

	prog = e.walker.Walk(&prog, &global)
	if len(e.walker.Errors) != 0 {
		colorstring.Println("[red]Failed walking:")
		for _, err := range e.walker.Errors {
			colorstring.Printf("[red]Error: %+v\n", err)
		}
		return fmt.Errorf("failed to walk through the nodes")
	}

	fmt.Printf("Walking time: %v seconds\n\n", time.Since(start).Seconds())
	start = time.Now()

	fmt.Println("Generating the lua code...")

	e.gen.Src.Grow(len(sourceFile))
	e.gen.Generate(prog)
	if len(e.gen.Errors) != 0 {
		colorstring.Println("[red]Failed generating:")
		for _, err := range e.gen.GetErrors() {
			colorstring.Printf("[red]Error: %+v\n", err)
		}
	}
	fmt.Printf("Generating time: %v seconds\n", time.Since(start).Seconds())

	err = os.WriteFile(e.DstPath, []byte(e.gen.GetSrc()), 0644)
	if err != nil {
		return fmt.Errorf("failed to write transpiled file to destination: %v", err)
	}

	return nil
}

func (e *Evaluator) writeSyntaxError(source string, errMsg ast.Error) {
	token := errMsg.Token

	sourceLines := strings.Split(source, "\n")
	line := sourceLines[token.Location.LineStart-1]

	fmt.Printf("line: %v in file yes\n", token.Location.LineStart)
	fmt.Println(line)
	if token.Location.ColStart-6 < 0 {
		fmt.Printf("%s^%s\n", strings.Repeat(" ", token.Location.ColStart-1), strings.Repeat("-", 5))
	} else {
		fmt.Printf("%s%s^\n", strings.Repeat(" ", token.Location.ColStart-6), strings.Repeat("-", 5))
	}
	fmt.Println("message: " + errMsg.Message + "\n")
}
