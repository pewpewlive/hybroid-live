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
	// for _, token := range tokens {
	// 	fmt.Printf("Token { type: %v, lex: %v, lit: %v, line: %v }\n", token.Type.ToString(), token.Lexeme, token.Literal, token.Line)
	// }

	p := parser.New()
	p.ParseTokens(tokens)
	if len(p.Errors) != 0 {
		for err := range p.Errors {
			fmt.Println(err)
		}
	}
	// actual compilation
	//var evaluators []evaluator.Evaluator

}
