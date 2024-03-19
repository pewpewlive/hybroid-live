package main

import (
	"fmt"
	"hybroid/generators/lua"
	"hybroid/lexer"
	"hybroid/parser"
	"os"
)

func main() {
	file, err := os.ReadFile("./example.hyb")
	if err != nil {
		fmt.Printf("[Error -> Main] reading file: %s\n", err.Error())
		return
	}

	l := lexer.New(file)
	l.Tokenize()

	p := parser.New(l.Tokens)
	prog := p.ParseTokens()

	global := lua.Global{Scope: lua.Scope{Global: nil, Parent: nil, Variables: make(map[string]lua.Value)}}
	global.Scope.Global = &global

	gen := lua.Generator{}
	gen.Program(prog, &global.Scope)

	fmt.Print(gen.Src)

	if len(p.Errors) != 0 {
		for _, err := range p.Errors {
			fmt.Printf("Error: %s, at line: %v (%v)\n", err.Message, err.Token.Location.LineStart, err.Token.ToString())
		}
	}
	// actual compilation
	//var evaluators []evaluator.Evaluator

}
