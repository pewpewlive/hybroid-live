package commands

import (
	"fmt"
	"hybroid/generators/lua"
	"hybroid/lexer"
	"hybroid/parser"
	"os"
)

func Transpile() error {
	file, err := os.ReadFile("example.hyb")

	l := lexer.New(file)
	l.Tokenize()

	l := lexer.New(file)
	l.Tokenize()

	p := parser.New(l.Tokens)
	prog := p.ParseTokens()

	global := lua.Global{Scope: lua.Scope{Global: nil, Parent: nil, Variables: make(map[string]lua.Value)}}
	global.Scope.Global = &global

	gen := lua.Generator{}
	gen.Generate(prog, &global.Scope)

	fmt.Print(gen.Src)

	
	for _, err := range gen.Errors {
		fmt.Printf("Error: %s, at line: %v (%v)\n", err.Message, err.Token.Location.LineStart, err.Token.ToString())
	}

	return nil
}
