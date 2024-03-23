package evaluator

import (
	"fmt"
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
	SrcPath string
	DstPath string
	gen     lua.Generator
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

	fmt.Printf("Tokenizing time: %v seconds\n", time.Since(start).Seconds())

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

	fmt.Println("Generating the lua code...")

	global := lua.Global{
		Scope: lua.Scope{
			Global:    nil,
			Parent:    nil,
			Variables: make(map[string]lua.Value),
		},
	}

	e.gen.Src.Grow(len(sourceFile))
	global.Scope.Global = &global
	e.gen.Generate(prog, &global.Scope)
	if len(e.gen.Errors) != 0 {
		colorstring.Println("[red]Failed generating:")
		for _, err := range e.gen.GetErrors() {
			colorstring.Printf("[red]Error: %+v\n", err)
		}
	}
	fmt.Printf("Build time: %v seconds\n", time.Since(start).Seconds())

	err = os.WriteFile(e.DstPath, []byte(e.gen.GetSrc()), 0644)
	if err != nil {
		return fmt.Errorf("failed to write transpiled file to destination: %v", err)
	}

	return nil
}

func (e *Evaluator) writeSyntaxError(source string, err parser.ParserError) {
	token := err.Token

	sourceLines := strings.Split(source, "\n")
	line := sourceLines[token.Location.LineStart-1]

	fmt.Printf("line: %v in file yes\n", token.Location.LineStart)
	fmt.Println(line)
	if token.Location.ColStart-6 < 0 {
		fmt.Printf("%s^%s\n", strings.Repeat(" ", token.Location.ColStart-1), strings.Repeat("-", 5))
	} else {
		fmt.Printf("%s%s^\n", strings.Repeat(" ", token.Location.ColStart-6), strings.Repeat("-", 5))
	}
	fmt.Println("message: " + err.Message + "\n")
}
