package evaluator

import (
	"fmt"
	//"hybroid/generators"
	"hybroid/generators/lua"
	"hybroid/lexer"
	"hybroid/parser"
	"os"
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
	e.lexer.ChangeSrc(lcsrc)
	e.lexer.Tokenize()
	if len(e.lexer.Errors) != 0 {
		fmt.Println("Failed tokenizing:")
		for _, err := range e.lexer.Errors {
			fmt.Printf("Error: %v\n", err)
		}
		return
	}
	e.parser.UpdateTokens(e.lexer.Tokens)
	prog := e.parser.ParseTokens()
	if len(e.parser.Errors) != 0 {
		fmt.Println("Failed parsing:")
		for _, err := range e.lexer.Errors {
			fmt.Printf("Error: %v\n", err)
		}
		return
	}
	global := lua.Global{
		Scope: lua.Scope{
			Global:    nil,
			Parent:    nil,
			Variables: make(map[string]lua.Value),
		},
	}
	global.Scope.Global = &global
	e.gen.Generate(prog, &global.Scope)
	if len(e.parser.Errors) != 0 {
		fmt.Println("Failed generating:")
		for _, err := range e.gen.GetErrors() {
			fmt.Printf("Error: %v\n", err)
		}
	}
	os.WriteFile(e.DstPath, []byte(e.gen.Src), 0677)
}
