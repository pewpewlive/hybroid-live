package main

import (
	"fmt"
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
	tokens := l.Tokenize()
	
	p := parser.New()
	p.ParseTokens(tokens)
	if len(p.Errors) != 0 {
		for _, err := range p.Errors {
			fmt.Println(err.Msg())
		}
	}
	// actual compilation
	//var evaluators []evaluator.Evaluator

}
