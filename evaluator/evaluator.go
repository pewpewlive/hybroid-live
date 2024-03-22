package evaluator

import (
	"fmt"
	"time"

	//"hybroid/generators"
	"hybroid/generators/lua"
	"hybroid/lexer"
	"hybroid/parser"
	"os"

	"github.com/mitchellh/colorstring"
)

func (e *Evaluator) HasValidSrc() bool {
	_, ok := os.ReadFile(e.SrcPath)

	return ok == nil
}

type Evaluator struct {
	lexer   lexer.Lexer
	parser  parser.Parser
	SrcPath string
	DstPath string
	gen     lua.Generator
}

func New(src string, dst string, gen lua.Generator) Evaluator {
	file, _ := os.ReadFile(src)
	return Evaluator{
		*lexer.New(file),
		parser.Parser{},
		src,
		dst,
		gen,
	}
}

func (e *Evaluator) Action() {
	lcsrc, _ := os.ReadFile(e.SrcPath)

	fmt.Printf("Tokenizing %v characters\n", len(lcsrc))
	start := time.Now()

	e.lexer.ChangeSrc(lcsrc)
	e.lexer.Tokenize()
	if len(e.lexer.Errors) != 0 {
		fmt.Println("[red]Failed tokenizing:")
		for _, err := range e.lexer.Errors {
			colorstring.Printf("[red]Error: %v\n", err)
		}
		return
	}

	fmt.Printf("Tokenizing time: %v seconds\n", time.Since(start).Seconds())

	fmt.Printf("Parsing %v tokens\n", len(e.lexer.Tokens))

	e.parser.UpdateTokens(e.lexer.Tokens)
	prog := e.parser.ParseTokens()
	if len(e.parser.Errors) != 0 {
		colorstring.Println("[red]Failed parsing:")
		for _, err := range e.parser.Errors {
			colorstring.Printf("[red]Error: %+v\n", err)
		}
		return
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

	e.gen.Src.Grow(len(lcsrc)) // for some reason this doesnt work
	global.Scope.Global = &global
	e.gen.Generate(prog, &global.Scope)
	if len(e.gen.Errors) != 0 {
		colorstring.Println("[red]Failed generating:")
		for _, err := range e.gen.GetErrors() {
			colorstring.Printf("[red]Error: %+v\n", err)
		}
	}
	fmt.Printf("Build time: %v seconds\n", time.Since(start).Seconds())

	os.WriteFile(e.DstPath, []byte(e.gen.GetSrc()), 0677)
}
