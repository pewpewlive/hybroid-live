package commands

import (
	"fmt"
	"hybroid/generators/lua"
	"hybroid/lexer"
	"hybroid/parser"
	"os"
)

func Build() error {
	file, _ := os.ReadFile("example.hyb")

	l := lexer.New(file)
	l.Tokenize()
	fmt.Println("Tokenized")

	p := parser.New(l.Tokens)
	prog := p.ParseTokens()
	fmt.Println("Parsed")

	global := lua.Global{Scope: lua.Scope{Global: nil, Parent: nil, Variables: make(map[string]lua.Value)}}
	global.Scope.Global = &global

	gen := lua.Generator{}
	gen.Generate(prog, &global.Scope)
	fmt.Println("Generated")

	fmt.Print(gen.Src)
	
	for _, err := range gen.Errors {
		fmt.Printf("Error: %s, at line: %v (%v)\n", err.Message, err.Token.Location.LineStart, err.Token.ToString())
	}

	return nil
}
